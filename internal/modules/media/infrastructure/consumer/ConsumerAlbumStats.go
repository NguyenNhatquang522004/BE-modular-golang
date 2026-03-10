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

type ConsumerAlbumStats struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	albumRepo IRepositoryMongodb.IAlbumsRepository
}

func NewConsumerAlbumStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, albumRepo IRepositoryMongodb.IAlbumsRepository) *ConsumerAlbumStats {
	return &ConsumerAlbumStats{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		albumRepo: albumRepo,
	}
}

func (c *ConsumerAlbumStats) ConsumerAlbumStats(ctx context.Context) error {
	// Implement the logic for consuming album stats here
	err := c.events.SubscribeBatch(ctx, constants.TopicAlbumStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
			ev := event
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreatedAlbumStats(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedAlbumStats(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedAlbumStats(ctx, ev)
				default:
					processErr = errors.New("unsupported event type: " + ev.Type)
					return
				}
			})
			if processErr != nil {
				errchan <- errors.New("failed to process event " + ev.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, ev.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, ev.ID)
			}
			if err != nil {
				errchan <- errors.New("failed to submit event " + ev.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, ev.ID)
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
func (c *ConsumerAlbumStats) handleCreatedAlbumStats(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created album stats here
	data, err := utils.ParsePayload[mediaEvent.AlbumStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	// Use the parsed data to perform necessary operations
	datablum, err := c.albumRepo.GetAlbumByID(ctx, data.AlbumID)
	if err != nil {
		return fmt.Errorf("failed to get album by ID %s: %w", data.AlbumID, err)
	}
	if datablum == nil {
		return fmt.Errorf("album with ID %s not found", data.AlbumID)
	}
	datablum.AssetCount = data.AssetCount + datablum.AssetCount
	datablum.Reactions.Angry = data.Angry + datablum.Reactions.Angry
	datablum.Reactions.Haha = data.Haha + datablum.Reactions.Haha
	datablum.Reactions.Like = data.Like + datablum.Reactions.Like
	datablum.Reactions.Sad = data.Sad + datablum.Reactions.Sad
	datablum.Reactions.Wow = data.Wow + datablum.Reactions.Wow
	datablum.Reactions.Total = datablum.Reactions.Total + data.Like + data.Love + data.Haha + data.Wow + data.Sad + data.Angry
	err = c.albumRepo.UpdateAlbum(ctx, datablum)
	if err != nil {
		return fmt.Errorf("failed to update album with ID %s: %w", data.AlbumID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for album ID %s: %w", data.AlbumID, err)
	}
	return nil
}

func (c *ConsumerAlbumStats) handleUpdatedAlbumStats(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created album stats here
	data, err := utils.ParsePayload[mediaEvent.AlbumStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	// Use the parsed data to perform necessary operations
	datablum, err := c.albumRepo.GetAlbumByID(ctx, data.AlbumID)
	if err != nil {
		return fmt.Errorf("failed to get album by ID %s: %w", data.AlbumID, err)
	}
	if datablum == nil {
		return fmt.Errorf("album with ID %s not found", data.AlbumID)
	}
	datablum.AssetCount = data.AssetCount + datablum.AssetCount
	datablum.Reactions.Angry = data.Angry + datablum.Reactions.Angry
	datablum.Reactions.Haha = data.Haha + datablum.Reactions.Haha
	datablum.Reactions.Like = data.Like + datablum.Reactions.Like
	datablum.Reactions.Sad = data.Sad + datablum.Reactions.Sad
	datablum.Reactions.Wow = data.Wow + datablum.Reactions.Wow
	datablum.Reactions.Total = datablum.Reactions.Total + data.Like + data.Love + data.Haha + data.Wow + data.Sad + data.Angry
	err = c.albumRepo.UpdateAlbum(ctx, datablum)
	if err != nil {
		return fmt.Errorf("failed to update album with ID %s: %w", data.AlbumID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for album ID %s: %w", data.AlbumID, err)
	}
	return nil
}

func (c *ConsumerAlbumStats) handleDeletedAlbumStats(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted album stats here
	data, err := utils.ParsePayload[mediaEvent.AlbumStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	// Use the parsed data to perform necessary operations
	datablum, err := c.albumRepo.GetAlbumByID(ctx, data.AlbumID)
	if err != nil {
		return fmt.Errorf("failed to get album by ID %s: %w", data.AlbumID, err)
	}
	if datablum == nil {
		return fmt.Errorf("album with ID %s not found", data.AlbumID)
	}
	datablum.AssetCount = datablum.AssetCount - data.AssetCount
	datablum.Reactions.Angry = datablum.Reactions.Angry - data.Angry
	datablum.Reactions.Haha = datablum.Reactions.Haha - data.Haha
	datablum.Reactions.Like = datablum.Reactions.Like - data.Like
	datablum.Reactions.Sad = datablum.Reactions.Sad - data.Sad
	datablum.Reactions.Wow = datablum.Reactions.Wow - data.Wow
	datablum.Reactions.Total = datablum.Reactions.Total - (data.Like + data.Love + data.Haha + data.Wow + data.Sad + data.Angry)
	err = c.albumRepo.UpdateAlbum(ctx, datablum)
	if err != nil {
		return fmt.Errorf("failed to update album with ID %s: %w", data.AlbumID, err)
	}
	err = c.handlemappingReactionCode(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to handle mapping reaction code for album ID %s: %w", data.AlbumID, err)
	}
	return nil
}
func (c *ConsumerAlbumStats) handlemappingReactionCode(ctx context.Context, data *mediaEvent.AlbumStatsPayload) error {
	typ := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := value.Field(i).Interface()
		if fieldName == "asset_count" || fieldName == "user_id" || fieldName == "album_id" || fieldName == "event_type" {
			continue
		}
		convertcql, err := gocql.ParseUUID(data.UserID)
		if err != nil {
			return fmt.Errorf("failed to convert user ID %s to CQL UUID: %w", data.UserID, err)
		}
		if fieldValue.(int) != 0 && fieldValue.(int) > 0 {
			code := c.mappingreactioncode(fieldName)
			if code == sharedEnums.ReactionUnknown {
				continue
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.AlbumID,
				UserID:       convertcql,
				TargetType:   sharedEnums.ReactionTargetAlbum,
				ReactionCode: code,
				CreatedAt:    time.Now(),
				Type:         data.EventType,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, data.EventType.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish entity reaction event for album ID %s: %w", data.AlbumID, err)
			}
		}
	}
	return nil
}
func (c *ConsumerAlbumStats) mappingreactioncode(data string) sharedEnums.ReactionCode {
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
func (c *ConsumerAlbumStats) ConsumerFailedAlbumStats(ctx context.Context) error {
	// Implement the logic for consuming failed album stats here
	return nil
}
