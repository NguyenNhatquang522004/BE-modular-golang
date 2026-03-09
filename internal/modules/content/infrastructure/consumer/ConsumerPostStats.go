package consumer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/gocql/gocql"
)

type ConsumerPostStats struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	postRepo  IRepositoryMongodb.IPostRepository
}

func NewConsumerPostStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, postRepo IRepositoryMongodb.IPostRepository) *ConsumerPostStats {
	return &ConsumerPostStats{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		postRepo:  postRepo,
	}
}
func (c *ConsumerPostStats) ConsumerPostStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicPostStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- err
				continue
			}
			if !can {
				if status == constants.StatusProcessing {
					errchan <- errors.New("event " + event.ID + " is currently being processed by another worker. Skipping.")
				} else {
					errchan <- errors.New("event " + event.ID + " has already been processed with status " + status.String() + ". Skipping.")
				}
				continue
			}
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreateStats(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdateStats(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeleteStats(ctx, event)
					// Xử lý logic cập nhật số liệu thống kê cho bài viết
					// Ví dụ: tăng/giảm số lượt thích, số lượt xem, v.v.
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
				}
			})
			if processErr != nil {
				errchan <- processErr
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchan <- err
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
				continue
			}
		}
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err)
			}
		}
		return finalErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerPostStats) handleCreateStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[contentEvent.PostStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to get post by ID %s: %w", data.PostID, err)
	}
	if datapost == nil {
		return fmt.Errorf("post not found for ID %s", data.PostID)
	}
	datapost.Stats.TotalReactions = datapost.Stats.TotalReactions + c.handleTotalReactions(ctx, *data)
	datapost.Stats.Comments = datapost.Stats.Comments + data.Comments
	datapost.Stats.Shares = datapost.Stats.Shares + data.Shares
	datapost.Stats.Views = datapost.Stats.Views + data.Views
	datapost.Stats.Like = datapost.Stats.Like + data.Like
	datapost.Stats.Love = datapost.Stats.Love + data.Love
	datapost.Stats.Haha = datapost.Stats.Haha + data.Haha
	datapost.Stats.Wow = datapost.Stats.Wow + data.Wow
	datapost.Stats.Sad = datapost.Stats.Sad + data.Sad
	datapost.Stats.Angry = datapost.Stats.Angry + data.Angry
	datapost.Stats.TopReactionTypes = c.handleTopReactionTypes(ctx, datapost.Stats)
	c.handleStatsWithReflect(ctx, *data)
	_, err = c.postRepo.UpdatePost(ctx, datapost)
	if err != nil {
		return fmt.Errorf("failed to update post stats for ID %s: %w", data.PostID, err)
	}
	return nil

}
func (c *ConsumerPostStats) handleTotalReactions(ctx context.Context, data contentEvent.PostStatsPayload) int {
	return data.Like + data.Love + data.Haha + data.Wow + data.Sad + data.Angry
}
func (c *ConsumerPostStats) handleTopReactionTypes(ctx context.Context, data entity.PostStats) []string {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	type reaction struct {
		ReactionCode string
		Count        int
	}
	var reactions []reaction
	for i := 0; i < val.NumField(); i++ {
		fieldInfo := typ.Field(i)
		jsonKey := fieldInfo.Tag.Get("json")
		if jsonKey == "like" || jsonKey == "love" || jsonKey == "haha" || jsonKey == "wow" || jsonKey == "sad" || jsonKey == "angry" {
			count := val.Field(i).Int()
			reactions = append(reactions, reaction{
				ReactionCode: jsonKey,
				Count:        int(count),
			})
		}
	}
	sort.Slice(reactions, func(i, j int) bool {
		return reactions[i].Count > reactions[j].Count
	})
	topReactions := make([]string, 0, 2)
	for i := 0; i < len(reactions) && i < 2; i++ {
		topReactions = append(topReactions, reactions[i].ReactionCode)
	}
	return topReactions
}
func (c *ConsumerPostStats) handleStatsWithReflect(ctx context.Context, payload contentEvent.PostStatsPayload) {
	val := reflect.ValueOf(payload)
	typ := reflect.TypeOf(payload)

	// Duyệt qua tất cả các field của struct
	for i := 0; i < val.NumField(); i++ {
		fieldInfo := typ.Field(i)
		fieldName := fieldInfo.Name
		// Bỏ qua PostID và TotalReactions, chỉ lấy từ Comments trở đi
		if fieldName == "PostID" || fieldName == "TotalReactions" {
			continue
		}
		if fieldName == "comments" {
			convertcql, err := gocql.ParseUUID(payload.UserID)
			if err != nil {
				continue
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     payload.PostID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetComment,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Topic:        constants.Created,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), payload.UserID.String(), constants.Created.String(), payload)
			if err != nil {
				continue
			}
			continue
		}
		if fieldName == "shares" {
			convertcql, err := gocql.ParseUUID(payload.UserID)
			if err != nil {
				continue
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     payload.PostID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetSharePost,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Topic:        constants.Created,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), payload.UserID.String(), constants.Created.String(), payload)
			if err != nil {
				continue
			}
			continue
		}
		if fieldName == "views" {
			convertcql, err := gocql.ParseUUID(payload.UserID)
			if err != nil {
				continue
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     payload.PostID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetViewPost,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Topic:        constants.Created,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), payload.UserID.String(), constants.Created.String(), payload)
			if err != nil {
				continue
			}
			continue
		}
		// Lấy giá trị của field (ép kiểu về int)
		fieldValue := val.Field(i).Int()

		if fieldValue == 1 {
			jsonKey := fieldInfo.Tag.Get("json")
			reactionCode := mapperReactioncodeToTargetType(jsonKey)
			convertcql, err := gocql.ParseUUID(payload.UserID)
			if err != nil {
				continue
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     payload.PostID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetPost,
				ReactionCode: reactionCode,
				CreatedAt:    time.Now(),
				Topic:        constants.Created,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), payload.UserID.String(), constants.Created.String(), payload)
			if err != nil {
				continue
			}
		}
	}
}
func (c *ConsumerPostStats) handleUpdateStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[contentEvent.PostStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to get post by ID %s: %w", data.PostID, err)
	}
	if datapost == nil {
		return fmt.Errorf("post not found for ID %s", data.PostID)
	}
	datapost.Stats.TotalReactions = datapost.Stats.TotalReactions + data.TotalReactions
	datapost.Stats.Comments = datapost.Stats.Comments + data.Comments
	datapost.Stats.Shares = datapost.Stats.Shares + data.Shares
	datapost.Stats.Views = datapost.Stats.Views + data.Views
	datapost.Stats.Like = datapost.Stats.Like + data.Like
	datapost.Stats.Love = datapost.Stats.Love + data.Love
	datapost.Stats.Haha = datapost.Stats.Haha + data.Haha
	datapost.Stats.Wow = datapost.Stats.Wow + data.Wow
	datapost.Stats.Sad = datapost.Stats.Sad + data.Sad
	datapost.Stats.Angry = datapost.Stats.Angry + data.Angry
	datapost.Stats.TopReactionTypes = c.handleTopReactionTypes(ctx, datapost.Stats)
	c.handleStatsWithReflect(ctx, *data)
	_, err = c.postRepo.UpdatePost(ctx, datapost)
	if err != nil {
		return fmt.Errorf("failed to update post stats for ID %s: %w", data.PostID, err)
	}
	return nil
}
func (c *ConsumerPostStats) handleDeleteStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[contentEvent.PostStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to get post by ID %s: %w", data.PostID, err)
	}
	if datapost == nil {
		return fmt.Errorf("post not found for ID %s", data.PostID)
	}
	datapost.Stats.TotalReactions = datapost.Stats.TotalReactions - data.TotalReactions
	datapost.Stats.Comments = datapost.Stats.Comments - data.Comments
	datapost.Stats.Shares = datapost.Stats.Shares - data.Shares
	datapost.Stats.Views = datapost.Stats.Views - data.Views
	datapost.Stats.Like = datapost.Stats.Like - data.Like
	datapost.Stats.Love = datapost.Stats.Love - data.Love
	datapost.Stats.Haha = datapost.Stats.Haha - data.Haha
	datapost.Stats.Wow = datapost.Stats.Wow - data.Wow
	datapost.Stats.Sad = datapost.Stats.Sad - data.Sad
	datapost.Stats.Angry = datapost.Stats.Angry - data.Angry
	datapost.Stats.TopReactionTypes = c.handleTopReactionTypes(ctx, datapost.Stats)
	c.handleStatsWithReflect(ctx, *data)
	_, err = c.postRepo.UpdatePost(ctx, datapost)
	if err != nil {
		return fmt.Errorf("failed to update post stats for ID %s: %w", data.PostID, err)
	}
	return nil
}

func mapperReactioncodeToTargetType(reactionCode string) sharedEnums.ReactionCode {
	switch reactionCode {
	case "like":
		return sharedEnums.ReactionLike
	case "love":
		return sharedEnums.ReactionLove
	case "haha":
		return sharedEnums.ReactionHaha
	case "wow":
		return sharedEnums.ReactionWow
	case "sad":
		return sharedEnums.ReactionSad
	case "angry":
		return sharedEnums.ReactionAngry
	default:
		return sharedEnums.ReactionUnknown
	}
}
func (c *ConsumerPostStats) ConsumerFailedPostStats(ctx context.Context) error {
	return nil
}
