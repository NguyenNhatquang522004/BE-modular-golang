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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
	"github.com/gocql/gocql"
)

type ConsumerCommentStats struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	comment   IRepositoryMongoDB.ICommentRepository
}

func NewConsumerCommentStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, comment IRepositoryMongoDB.ICommentRepository) *ConsumerCommentStats {
	return &ConsumerCommentStats{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		comment:   comment,
	}
}
func (c *ConsumerCommentStats) ConsumerCommentStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicCommentStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedCommentStats(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedCommentStats(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedCommentStats(ctx, event)
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
func (c *ConsumerCommentStats) handleCreatedCommentStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.CommentStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datacomment, err := c.comment.GetCommentByID(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get comment by ID: %w", err)
	}
	if datacomment == nil {
		return fmt.Errorf("comment with ID %s not found", data.TargetID)
	}
	datacomment.ReplyCount = datacomment.ReplyCount + data.ReplyCount
	datacomment.MentionCount = datacomment.MentionCount + data.MentionCount
	datacomment.Reactions.Like = datacomment.Reactions.Like + data.Like
	datacomment.Reactions.Love = datacomment.Reactions.Love + data.Love
	datacomment.Reactions.Wow = datacomment.Reactions.Wow + data.Wow
	datacomment.Reactions.Sad = datacomment.Reactions.Sad + data.Sad
	datacomment.Reactions.Angry = datacomment.Reactions.Angry + data.Angry
	datacomment.Reactions.Total = datacomment.Reactions.Total + data.Like + data.Love + data.Wow + data.Sad + data.Angry
	err = c.comment.UpdateComment(ctx, datacomment)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	err = c.SendEntityReactionEvent(ctx, *data)
	if err != nil {
		return fmt.Errorf("failed to send entity reaction event: %w", err)
	}
	return nil
}
func (c *ConsumerCommentStats) SendEntityReactionEvent(ctx context.Context, payloadinit interactionEvent.CommentStatsPayload) error {
	val := reflect.ValueOf(payloadinit)
	typ := reflect.TypeOf(payloadinit)
	converuserid, err := gocql.ParseUUID(payloadinit.UserID)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
		if field.Name == "Type" || field.Name == "TargetID" || field.Name == "UserID" || field.Name == "MentionCount" || field.Name == "Total" {
			continue
		}
		if field.Name == "reply_count" {
			payload := interactionEvent.EntityReactionPayload{
				TargetID:     payloadinit.TargetID,
				UserID:       converuserid,
				TargetType:   sharedEnums.ReactionTargetReplyComment,
				ReactionCode: sharedEnums.ReactionUnknown, // Giả sử reply_count chỉ tăng khi có like, bạn có thể điều chỉnh logic này tùy theo yêu cầu thực tế
				CreatedAt:    time.Now(),
				Type:         payloadinit.Type,
			}
			err := c.events.Publish(ctx, constants.TopicEntityReaction.String(), payloadinit.TargetID, string(payloadinit.Type), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event: %w", err)
			}
		}
		if value.(int) > 0 {
			reactionCode := mappingreactioncode(field.Name)
			payload := interactionEvent.EntityReactionPayload{
				TargetID:     payloadinit.TargetID,
				UserID:       converuserid,
				TargetType:   sharedEnums.ReactionTargetComment,
				ReactionCode: reactionCode,
				CreatedAt:    time.Now(),
				Type:         payloadinit.Type,
			}
			err := c.events.Publish(ctx, constants.TopicEntityReaction.String(), payloadinit.TargetID, string(payloadinit.Type), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event: %w", err)
			}
		}
	}
	return nil
}
func mappingreactioncode(fieldName string) sharedEnums.ReactionCode {
	switch fieldName {
	case "Like":
		return sharedEnums.ReactionLike
	case "Love":
		return sharedEnums.ReactionLove
	case "Wow":
		return sharedEnums.ReactionWow
	case "Sad":
		return sharedEnums.ReactionSad
	case "Angry":
		return sharedEnums.ReactionAngry
	default:
		return sharedEnums.ReactionUnknown
	}
}
func (c *ConsumerCommentStats) handleUpdatedCommentStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.CommentStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datacomment, err := c.comment.GetCommentByID(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get comment by ID: %w", err)
	}
	if datacomment == nil {
		return fmt.Errorf("comment with ID %s not found", data.TargetID)
	}
	datacomment.ReplyCount = datacomment.ReplyCount + data.ReplyCount
	datacomment.MentionCount = datacomment.MentionCount + data.MentionCount
	datacomment.Reactions.Like = datacomment.Reactions.Like + data.Like
	datacomment.Reactions.Love = datacomment.Reactions.Love + data.Love
	datacomment.Reactions.Wow = datacomment.Reactions.Wow + data.Wow
	datacomment.Reactions.Sad = datacomment.Reactions.Sad + data.Sad
	datacomment.Reactions.Angry = datacomment.Reactions.Angry + data.Angry
	datacomment.Reactions.Total = datacomment.Reactions.Total + data.Like + data.Love + data.Wow + data.Sad + data.Angry
	err = c.comment.UpdateComment(ctx, datacomment)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	err = c.SendEntityReactionEvent(ctx, *data)
	if err != nil {
		return fmt.Errorf("failed to send entity reaction event: %w", err)
	}
	return nil
}
func (c *ConsumerCommentStats) handleDeletedCommentStats(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.CommentStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datacomment, err := c.comment.GetCommentByID(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get comment by ID: %w", err)
	}
	if datacomment == nil {
		return fmt.Errorf("comment with ID %s not found", data.TargetID)
	}
	datacomment.ReplyCount = datacomment.ReplyCount - data.ReplyCount
	datacomment.MentionCount = datacomment.MentionCount - data.MentionCount
	datacomment.Reactions.Like = datacomment.Reactions.Like - data.Like
	datacomment.Reactions.Love = datacomment.Reactions.Love - data.Love
	datacomment.Reactions.Wow = datacomment.Reactions.Wow - data.Wow
	datacomment.Reactions.Sad = datacomment.Reactions.Sad - data.Sad
	datacomment.Reactions.Angry = datacomment.Reactions.Angry - data.Angry
	datacomment.Reactions.Total = datacomment.Reactions.Total - data.Like - data.Love - data.Wow - data.Sad - data.Angry
	err = c.comment.UpdateComment(ctx, datacomment)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}
func (c *ConsumerCommentStats) ConsumerFailedCommentStats(ctx context.Context) error {
	return nil
}
