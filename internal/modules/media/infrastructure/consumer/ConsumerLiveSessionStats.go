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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
	"github.com/gocql/gocql"
)

type ConsumerLiveSessionStats struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	liveSessionRepo IRepositoryMongodb.ILiveSessionRepository
}

func NewConsumerLiveSessionStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, liveSessionRepo IRepositoryMongodb.ILiveSessionRepository) *ConsumerLiveSessionStats {
	return &ConsumerLiveSessionStats{
		events:          events,
		pool:            pool,
		redisRepo:       redisRepo,
		liveSessionRepo: liveSessionRepo,
	}
}

func (c *ConsumerLiveSessionStats) ConsumerLiveSessionStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicLiveSessionStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		var finalErr error
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
					processErr = c.handleCreatedEvent(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedEvent(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedEvent(ctx, event)
				default:
					processErr = errors.New("unknown event type: " + event.Type)

				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Decrement the WaitGroup counter since the task won't be processed
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
		}
		wg.Wait()
		close(errchan)
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
func (c *ConsumerLiveSessionStats) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.LiveSessionStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datalive, err := c.liveSessionRepo.GetLiveSessionByID(ctx, data.LiveSessionID)
	if err != nil {
		return fmt.Errorf("failed to get live session by ID %s: %w", data.LiveSessionID, err)
	}
	if datalive == nil {
		return fmt.Errorf("live session with ID %s not found", data.LiveSessionID)
	}
	datalive.Stats.PeakViewers = data.PeakViewers + datalive.Stats.PeakViewers
	datalive.Stats.TotalViews = data.TotalViews + datalive.Stats.TotalViews
	datalive.Stats.TotalComments = data.TotalComments + datalive.Stats.TotalComments
	// Cập nhật reaction counts nếu có
	datalive.Stats.Like += data.Like
	datalive.Stats.Love += data.Love
	datalive.Stats.Haha += data.Haha
	datalive.Stats.Wow += data.Wow
	datalive.Stats.Sad += data.Sad
	datalive.Stats.Angry += data.Angry
	err = c.liveSessionRepo.UpdateLiveSession(ctx, datalive)
	if err != nil {
		return fmt.Errorf("failed to update live session with ID %s: %w", data.LiveSessionID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for live session ID %s: %w", data.LiveSessionID, err)
	}
	return nil

}

func (c *ConsumerLiveSessionStats) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.LiveSessionStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datalive, err := c.liveSessionRepo.GetLiveSessionByID(ctx, data.LiveSessionID)
	if err != nil {
		return fmt.Errorf("failed to get live session by ID %s: %w", data.LiveSessionID, err)
	}
	if datalive == nil {
		return fmt.Errorf("live session with ID %s not found", data.LiveSessionID)
	}
	datalive.Stats.PeakViewers = data.PeakViewers + datalive.Stats.PeakViewers
	datalive.Stats.TotalViews = data.TotalViews + datalive.Stats.TotalViews
	datalive.Stats.TotalComments = data.TotalComments + datalive.Stats.TotalComments
	// Cập nhật reaction counts nếu có
	datalive.Stats.Like += data.Like
	datalive.Stats.Love += data.Love
	datalive.Stats.Haha += data.Haha
	datalive.Stats.Wow += data.Wow
	datalive.Stats.Sad += data.Sad
	datalive.Stats.Angry += data.Angry
	err = c.liveSessionRepo.UpdateLiveSession(ctx, datalive)
	if err != nil {
		return fmt.Errorf("failed to update live session with ID %s: %w", data.LiveSessionID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for live session ID %s: %w", data.LiveSessionID, err)
	}
	return nil
}
func (c *ConsumerLiveSessionStats) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.LiveSessionStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datalive, err := c.liveSessionRepo.GetLiveSessionByID(ctx, data.LiveSessionID)
	if err != nil {
		return fmt.Errorf("failed to get live session by ID %s: %w", data.LiveSessionID, err)
	}
	if datalive == nil {
		return fmt.Errorf("live session with ID %s not found", data.LiveSessionID)
	}
	datalive.Stats.TotalViews = data.TotalViews - datalive.Stats.TotalViews
	datalive.Stats.TotalComments = data.TotalComments - datalive.Stats.TotalComments
	// Cập nhật reaction counts nếu có
	datalive.Stats.Like -= data.Like
	datalive.Stats.Love -= data.Love
	datalive.Stats.Haha -= data.Haha
	datalive.Stats.Wow -= data.Wow
	datalive.Stats.Sad -= data.Sad
	datalive.Stats.Angry -= data.Angry
	err = c.liveSessionRepo.UpdateLiveSession(ctx, datalive)
	if err != nil {
		return fmt.Errorf("failed to update live session with ID %s: %w", data.LiveSessionID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for live session ID %s: %w", data.LiveSessionID, err)
	}
	return nil
}
func (c *ConsumerLiveSessionStats) ConsumerFailedLiveSessionStats(ctx context.Context) error {
	return nil
}
func (c *ConsumerLiveSessionStats) handlemappingReactionCode(ctx context.Context, data *mediaEvent.LiveSessionStatsPayload) error {
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
		if fieldName == "asset_count" || fieldName == "user_id" || fieldName == "live_session_id" || fieldName == "event_type" {
			continue
		}
		if fieldName == "total_views" && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.LiveSessionID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetViewLive,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for live session ID %s: %w", data.LiveSessionID, err)
			}
		}
		if fieldName == "total_comments" && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.LiveSessionID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetCommentLive,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for live session ID %s: %w", data.LiveSessionID, err)
			}
		}
		if fieldValue.(int) != 0 && fieldValue.(int) > 0 {
			reactionCode := c.mappingreactioncode(fieldName)
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.LiveSessionID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetLive,
				ReactionCode: reactionCode,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for live session ID %s: %w", data.LiveSessionID, err)
			}
		}
	}
	return nil
}
func (c *ConsumerLiveSessionStats) GetReactionCode(ctx context.Context, data *mediaEvent.LiveSessionStatsPayload) sharedEnums.ReactionCode {
	typ := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	flag := false
	var reactionCode sharedEnums.ReactionCode
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := value.Field(i).Interface()
		if fieldName == "asset_count" || fieldName == "user_id" || fieldName == "live_session_id" || fieldName == "event_type" || fieldName == "total_comments" || fieldName == "total_views" || fieldName == "peak_viewers" {
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
func (c *ConsumerLiveSessionStats) mappingreactioncode(data string) sharedEnums.ReactionCode {
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
