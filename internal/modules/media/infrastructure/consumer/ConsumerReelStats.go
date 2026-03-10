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

type ReelStatsConsumer struct {
	// Add any necessary fields for the consumer, such as a reference to the event bus or database
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	reelRepo  IRepositoryMongodb.IReelRepository
}

func NewReelStatsConsumer(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, reelRepo IRepositoryMongodb.IReelRepository) *ReelStatsConsumer {
	return &ReelStatsConsumer{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		reelRepo:  reelRepo,
	}
}

func (c *ReelStatsConsumer) ConsumerReelStats(ctx context.Context) error {
	// Implement the logic to consume reel stats events here
	err := c.events.SubscribeBatch(ctx, constants.TopicReelStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedEvent(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedEvent(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedEvent(ctx, event)
					// Add more cases for other event types as needed
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
					errchan <- processErr
					return
				}
			})
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				// Unlock the event after processing
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchan <- errors.New("failed to run event " + event.ID + " in worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Ensure we call Done() if we fail to run the task
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
func (c *ReelStatsConsumer) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic to handle created reel stats events here
	data, err := utils.ParsePayload[mediaEvent.ReelStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datareel, err := c.reelRepo.GetReelByID(ctx, data.ReelID)
	if err != nil {
		return fmt.Errorf("failed to get reel by ID %s: %w", data.ReelID, err)
	}
	if datareel == nil {
		return fmt.Errorf("reel not found for ID %s", data.ReelID)
	}
	datareel.Stats.Views = data.Views + datareel.Stats.Views
	datareel.Stats.Shares = data.Shares + datareel.Stats.Shares
	datareel.Stats.Comments = data.Comments + datareel.Stats.Comments
	datareel.Stats.Saves = data.Saves + datareel.Stats.Saves
	datareel.Stats.Like = data.Like + datareel.Stats.Like
	datareel.Stats.Love = data.Love + datareel.Stats.Love
	datareel.Stats.Haha = data.Haha + datareel.Stats.Haha
	datareel.Stats.Wow = data.Wow + datareel.Stats.Wow
	datareel.Stats.Sad = data.Sad + datareel.Stats.Sad
	datareel.Stats.Angry = data.Angry + datareel.Stats.Angry
	datareel.Stats.Total = datareel.Stats.Views + datareel.Stats.Like + datareel.Stats.Love + datareel.Stats.Haha + datareel.Stats.Wow + datareel.Stats.Sad + datareel.Stats.Angry
	err = c.reelRepo.UpdateReel(ctx, datareel)
	// Process the data as needed
	if err != nil {
		return fmt.Errorf("failed to update reel stats for reel ID %s: %w", data.ReelID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for reel ID %s: %w", data.ReelID, err)
	}
	return nil
}

func (c *ReelStatsConsumer) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic to handle created reel stats events here
	data, err := utils.ParsePayload[mediaEvent.ReelStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datareel, err := c.reelRepo.GetReelByID(ctx, data.ReelID)
	if err != nil {
		return fmt.Errorf("failed to get reel by ID %s: %w", data.ReelID, err)
	}
	if datareel == nil {
		return fmt.Errorf("reel not found for ID %s", data.ReelID)
	}
	datareel.Stats.Views = data.Views + datareel.Stats.Views
	datareel.Stats.Shares = data.Shares + datareel.Stats.Shares
	datareel.Stats.Comments = data.Comments + datareel.Stats.Comments
	datareel.Stats.Saves = data.Saves + datareel.Stats.Saves
	datareel.Stats.Like = data.Like + datareel.Stats.Like
	datareel.Stats.Love = data.Love + datareel.Stats.Love
	datareel.Stats.Haha = data.Haha + datareel.Stats.Haha
	datareel.Stats.Wow = data.Wow + datareel.Stats.Wow
	datareel.Stats.Sad = data.Sad + datareel.Stats.Sad
	datareel.Stats.Angry = data.Angry + datareel.Stats.Angry
	datareel.Stats.Total = datareel.Stats.Views + data.Like + data.Love + data.Haha + data.Wow + data.Sad + data.Angry
	err = c.reelRepo.UpdateReel(ctx, datareel)
	// Process the data as needed
	if err != nil {
		return fmt.Errorf("failed to update reel stats for reel ID %s: %w", data.ReelID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for reel ID %s: %w", data.ReelID, err)
	}
	return nil
}

func (c *ReelStatsConsumer) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic to handle created reel stats events here
	data, err := utils.ParsePayload[mediaEvent.ReelStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datareel, err := c.reelRepo.GetReelByID(ctx, data.ReelID)
	if err != nil {
		return fmt.Errorf("failed to get reel by ID %s: %w", data.ReelID, err)
	}
	if datareel == nil {
		return fmt.Errorf("reel not found for ID %s", data.ReelID)
	}
	datareel.Stats.Views = data.Views - datareel.Stats.Views
	datareel.Stats.Shares = data.Shares - datareel.Stats.Shares
	datareel.Stats.Comments = data.Comments - datareel.Stats.Comments
	datareel.Stats.Saves = data.Saves - datareel.Stats.Saves
	datareel.Stats.Like = data.Like - datareel.Stats.Like
	datareel.Stats.Love = data.Love - datareel.Stats.Love
	datareel.Stats.Haha = data.Haha - datareel.Stats.Haha
	datareel.Stats.Wow = data.Wow - datareel.Stats.Wow
	datareel.Stats.Sad = data.Sad - datareel.Stats.Sad
	datareel.Stats.Angry = data.Angry - datareel.Stats.Angry
	datareel.Stats.Total = datareel.Stats.Views - (data.Like + data.Love + data.Haha + data.Wow + data.Sad + data.Angry)
	err = c.reelRepo.UpdateReel(ctx, datareel)
	// Process the data as needed
	if err != nil {
		return fmt.Errorf("failed to update reel stats for reel ID %s: %w", data.ReelID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for reel ID %s: %w", data.ReelID, err)
	}
	return nil
}

func (c *ReelStatsConsumer) ConsumerFailedReelStats(ctx context.Context) error {
	// Implement the logic to consume failed reel stats events here
	return nil
}

func (c *ReelStatsConsumer) handlemappingReactionCode(ctx context.Context, data *mediaEvent.ReelStatsPayload) error {
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
		if fieldName == "user_id" || fieldName == "reel_id" || fieldName == "event_type" {
			continue
		}
		if fieldName == "views" && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.ReelID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetViewReel,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for reel ID %s: %w", data.ReelID, err)
			}
		}
		if fieldName == "shares" && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.ReelID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetShareReel,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for reel ID %s: %w", data.ReelID, err)
			}
		}
		if fieldName == "comments" && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.ReelID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetCommentReel,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for reel ID %s: %w", data.ReelID, err)
			}
		}
		if fieldName == "saves" && fieldValue.(int) > 0 {
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.ReelID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetSaveReel,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for reel ID %s: %w", data.ReelID, err)
			}
		}
		if fieldValue.(int) != 0 && fieldValue.(int) > 0 {
			reactionCode := c.mappingreactioncode(fieldName)
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.ReelID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetReel,
				ReactionCode: reactionCode,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for reel ID %s: %w", data.ReelID, err)
			}
		}
	}
	return nil
}
func (c *ReelStatsConsumer) GetReactionCode(ctx context.Context, data *mediaEvent.ReelStatsPayload) sharedEnums.ReactionCode {
	typ := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	flag := false
	var reactionCode sharedEnums.ReactionCode
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := value.Field(i).Interface()
		if fieldName == "asset_count" || fieldName == "user_id" || fieldName == "reel_id" || fieldName == "event_type" || fieldName == "views" || fieldName == "shares" || fieldName == "comments" || fieldName == "saves" {
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
func (c *ReelStatsConsumer) mappingreactioncode(data string) sharedEnums.ReactionCode {
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
