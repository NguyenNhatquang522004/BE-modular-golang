package consumer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"github.com/gocql/gocql"
)

type ConsumerStoryStats struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	storyRepo     IRepositoryMongodb.IStoryRepository
	storyViewRepo IRepositoryCassandra.IStoryViewRepository
	pb            pb.SocialServiceClient
}

func NewConsumerStoryStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, storyRepo IRepositoryMongodb.IStoryRepository, storyViewRepo IRepositoryCassandra.IStoryViewRepository, pbClient pb.SocialServiceClient) *ConsumerStoryStats {
	return &ConsumerStoryStats{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		storyRepo:     storyRepo,
		storyViewRepo: storyViewRepo,
		pb:            pbClient,
	}
}

func (c *ConsumerStoryStats) ConsumerStoryStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicStoryStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- errors.New("failed to acquire lock for event " + event.ID + ": " + err.Error())
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
					processErr = c.handleCreatedStoryStats(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedStoryStats(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedStoryStats(ctx, event)
				default:
					processErr = errors.New("unsupported event type " + event.Type + " for event " + event.ID)
					return
				}
			})
			if processErr != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)

			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
			}
		}
		wg.Wait()
		close(errchan)
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
func (c *ConsumerStoryStats) handleCreatedStoryStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.StoryStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event " + event.ID + ": " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event " + event.ID))
	}
	datastory, err := c.storyRepo.GetStoryByID(ctx, data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to get story by ID " + data.StoryID + ": " + err.Error())
	}
	if datastory == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("story not found for ID " + data.StoryID))
	}
	datainfo, err := c.pb.GetInfoUserByID(ctx, &pb.UserSocialIDRequest{UserId: data.UserID})
	if err != nil {
		return fmt.Errorf("failed to get user info by ID %s: %w", data.UserID, err)
	}
	datastory.Stats.Likes = datastory.Stats.Likes + data.Like
	datastory.Stats.Love = datastory.Stats.Love + data.Love
	datastory.Stats.Haha = datastory.Stats.Haha + data.Haha
	datastory.Stats.Wow = datastory.Stats.Wow + data.Wow
	datastory.Stats.Sad = datastory.Stats.Sad + data.Sad
	datastory.Stats.Angry = datastory.Stats.Angry + data.Angry
	datastory.Stats.ReplyCount = datastory.Stats.ReplyCount + data.ReplyCount
	datastory.Stats.ViewsCount = datastory.Stats.ViewsCount + data.ViewsCount
	if len(datastory.PreviewViewers) == 3 {
		datastory.PreviewViewers = datastory.PreviewViewers[1:]
		datastory.PreviewViewers = append(datastory.PreviewViewers, entity.ViewerPreview{
			UserID: data.UserID,
			Avatar: datainfo.GetAuthorAvatar(),
			Name:   datainfo.GetAuthorName(),
		})
	} else {
		datastory.PreviewViewers = append(datastory.PreviewViewers, entity.ViewerPreview{
			UserID: data.UserID,
			Avatar: datainfo.GetAuthorAvatar(),
			Name:   datainfo.GetAuthorName(),
		})
	}
	err = c.storyRepo.UpdateStory(ctx, datastory)
	if err != nil {
		return fmt.Errorf("failed to update story with ID %s: %w", data.StoryID, err)
	}
	convertcqlstoryID, err := gocql.ParseUUID(data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to convert story ID %s to CQL UUID: %w", data.StoryID, err)
	}
	convertcqlUserID, err := gocql.ParseUUID(data.UserID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID %s to CQL UUID: %w", data.UserID, err)
	}
	entity := &entity.StoryView{
		StoryID:         convertcqlstoryID,
		ViewerID:        convertcqlUserID,
		ViewedAt:        time.Now(),
		ViewerName:      datainfo.GetAuthorName(),
		ViewerAvatarURL: datainfo.GetAuthorAvatar(),
		InteractionType: data.InteractionType,
		ReactionCode:    c.GetReactionCode(ctx, data),
		Content:         data.Content,
		PollOptionIndex: data.PollOptionIndex,
	}
	err = c.storyViewRepo.CreateStoryView(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create story view for story ID %s and user ID %s: %w", data.StoryID, data.UserID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for story ID %s: %w", data.StoryID, err)
	}

	return nil
}

func (c *ConsumerStoryStats) handleUpdatedStoryStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.StoryStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event " + event.ID + ": " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event " + event.ID))
	}
	datastory, err := c.storyRepo.GetStoryByID(ctx, data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to get story by ID " + data.StoryID + ": " + err.Error())
	}
	if datastory == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("story not found for ID " + data.StoryID))
	}
	datainfo, err := c.pb.GetInfoUserByID(ctx, &pb.UserSocialIDRequest{UserId: data.UserID})
	if err != nil {
		return fmt.Errorf("failed to get user info by ID %s: %w", data.UserID, err)
	}
	datastory.Stats.Likes = datastory.Stats.Likes + data.Like
	datastory.Stats.Love = datastory.Stats.Love + data.Love
	datastory.Stats.Haha = datastory.Stats.Haha + data.Haha
	datastory.Stats.Wow = datastory.Stats.Wow + data.Wow
	datastory.Stats.Sad = datastory.Stats.Sad + data.Sad
	datastory.Stats.Angry = datastory.Stats.Angry + data.Angry
	datastory.Stats.ReplyCount = datastory.Stats.ReplyCount + data.ReplyCount
	datastory.Stats.ViewsCount = datastory.Stats.ViewsCount + data.ViewsCount
	if len(datastory.PreviewViewers) == 3 {
		datastory.PreviewViewers = datastory.PreviewViewers[1:]
		datastory.PreviewViewers = append(datastory.PreviewViewers, entity.ViewerPreview{
			UserID: data.UserID,
			Avatar: datainfo.GetAuthorAvatar(),
			Name:   datainfo.GetAuthorName(),
		})
	} else {
		datastory.PreviewViewers = append(datastory.PreviewViewers, entity.ViewerPreview{
			UserID: data.UserID,
			Avatar: datainfo.GetAuthorAvatar(),
			Name:   datainfo.GetAuthorName(),
		})
	}
	err = c.storyRepo.UpdateStory(ctx, datastory)
	if err != nil {
		return fmt.Errorf("failed to update story with ID %s: %w", data.StoryID, err)
	}
	convertcqlstoryID, err := gocql.ParseUUID(data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to convert story ID %s to CQL UUID: %w", data.StoryID, err)
	}
	convertcqlUserID, err := gocql.ParseUUID(data.UserID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID %s to CQL UUID: %w", data.UserID, err)
	}
	entity := &entity.StoryView{
		StoryID:         convertcqlstoryID,
		ViewerID:        convertcqlUserID,
		ViewedAt:        time.Now(),
		ViewerName:      datainfo.GetAuthorName(),
		ViewerAvatarURL: datainfo.GetAuthorAvatar(),
		InteractionType: data.InteractionType,
		ReactionCode:    c.GetReactionCode(ctx, data),
		Content:         data.Content,
		PollOptionIndex: data.PollOptionIndex,
	}
	err = c.storyViewRepo.CreateStoryView(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create story view for story ID %s and user ID %s: %w", data.StoryID, data.UserID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for story ID %s: %w", data.StoryID, err)
	}

	return nil
}
func (c *ConsumerStoryStats) handleDeletedStoryStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.StoryStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event " + event.ID + ": " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event " + event.ID))
	}
	datastory, err := c.storyRepo.GetStoryByID(ctx, data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to get story by ID " + data.StoryID + ": " + err.Error())
	}
	if datastory == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("story not found for ID " + data.StoryID))
	}
	datainfo, err := c.pb.GetInfoUserByID(ctx, &pb.UserSocialIDRequest{UserId: data.UserID})
	if err != nil {
		return fmt.Errorf("failed to get user info by ID %s: %w", data.UserID, err)
	}
	datastory.Stats.Likes = datastory.Stats.Likes - data.Like
	datastory.Stats.Love = datastory.Stats.Love - data.Love
	datastory.Stats.Haha = datastory.Stats.Haha - data.Haha
	datastory.Stats.Wow = datastory.Stats.Wow - data.Wow
	datastory.Stats.Sad = datastory.Stats.Sad - data.Sad
	datastory.Stats.Angry = datastory.Stats.Angry - data.Angry
	datastory.Stats.ReplyCount = datastory.Stats.ReplyCount - data.ReplyCount
	datastory.Stats.ViewsCount = datastory.Stats.ViewsCount - data.ViewsCount
	if len(datastory.PreviewViewers) == 3 {
		datastory.PreviewViewers = datastory.PreviewViewers[1:]
		datastory.PreviewViewers = append(datastory.PreviewViewers, entity.ViewerPreview{
			UserID: data.UserID,
			Avatar: datainfo.GetAuthorAvatar(),
			Name:   datainfo.GetAuthorName(),
		})
	} else {
		datastory.PreviewViewers = append(datastory.PreviewViewers, entity.ViewerPreview{
			UserID: data.UserID,
			Avatar: datainfo.GetAuthorAvatar(),
			Name:   datainfo.GetAuthorName(),
		})
	}
	err = c.storyRepo.UpdateStory(ctx, datastory)
	if err != nil {
		return fmt.Errorf("failed to update story with ID %s: %w", data.StoryID, err)
	}
	convertcqlstoryID, err := gocql.ParseUUID(data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to convert story ID %s to CQL UUID: %w", data.StoryID, err)
	}
	convertcqlUserID, err := gocql.ParseUUID(data.UserID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID %s to CQL UUID: %w", data.UserID, err)
	}
	entity := &entity.StoryView{
		StoryID:         convertcqlstoryID,
		ViewerID:        convertcqlUserID,
		ViewedAt:        time.Now(),
		ViewerName:      datainfo.GetAuthorName(),
		ViewerAvatarURL: datainfo.GetAuthorAvatar(),
		InteractionType: data.InteractionType,
		ReactionCode:    c.GetReactionCode(ctx, data),
		Content:         data.Content,
		PollOptionIndex: data.PollOptionIndex,
	}
	err = c.storyViewRepo.CreateStoryView(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create story view for story ID %s and user ID %s: %w", data.StoryID, data.UserID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for story ID %s: %w", data.StoryID, err)
	}

	return nil
}

func (c *ConsumerStoryStats) handlemappingReactionCode(ctx context.Context, data *mediaEvent.StoryStatsPayload) error {
	convertcql, err := gocql.ParseUUID(data.UserID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID %s to CQL UUID: %w", data.UserID, err)
	}
	typ := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := value.Field(i).Interface()
		if fieldName == "asset_count" || fieldName == "user_id" || fieldName == "story_id" || fieldName == "event_type" {
			continue
		}
		if fieldName == "reply_count"  && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.StoryID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetReplyStory,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for story ID %s: %w", data.StoryID, err)
			}
		}
		if fieldName == "views_count"  && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.StoryID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetViewStory,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for story ID %s: %w", data.StoryID, err)
			}
		}
		if fieldValue.(int) != 0 && fieldValue.(int) > 0 {
			reactionCode := c.mappingreactioncode(fieldName)
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.StoryID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetStory,
				ReactionCode: reactionCode,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for story ID %s: %w", data.StoryID, err)
			}
		}
	}
	return nil
}
func (c *ConsumerStoryStats) GetReactionCode(ctx context.Context, data *mediaEvent.StoryStatsPayload) sharedEnums.ReactionCode {
	typ := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	flag := false
	var reactionCode sharedEnums.ReactionCode
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := value.Field(i).Interface()
		if fieldName == "asset_count" || fieldName == "user_id" || fieldName == "story_id" || fieldName == "event_type" || fieldName == "reply_count" || fieldName == "views_count" {
			continue
		}
		if fieldValue.(int) != 0 && fieldValue.(int) > 0 {
			flag = true
			reactionCode = c.mappingreactioncode(fieldName)
		}
	}
	if !flag {
		return sharedEnums.ReactionUnknown
	}
	return reactionCode
}
func (c *ConsumerStoryStats) mappingreactioncode(data string) sharedEnums.ReactionCode {
	switch data {
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
func (c *ConsumerStoryStats) ConsumerFailedStoryStats(ctx context.Context) error {
	return nil
}
