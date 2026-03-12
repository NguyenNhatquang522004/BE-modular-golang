package consumer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerReel struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	reelRepo  IRepositoryMongodb.IReelRepository
}

func NewConsumerReel(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, reelRepo IRepositoryMongodb.IReelRepository) *ConsumerReel {
	return &ConsumerReel{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		reelRepo:  reelRepo,
	}
}
func (c *ConsumerReel) ConsumerReel(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, string(constants.TopicReel), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
				}
			})
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
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
func (c *ConsumerReel) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.CreateReelPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity, err := mapper.ToReelEntity(data, data.UserID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	err = c.reelRepo.CreateReel(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create reel in repository for event %s: %w", event.ID, err)
	}

	payloadpost := &contentEvent.CreatePostPayload{
		ID:      entity.ID.Hex(), // ID của Post sẽ trùng với ID của Reel để dễ quản lý, client không gửi lên
		UserID:  data.UserID,
		Type:    sharedEnums.PostTypeMedia,
		Content: data.Caption,
		Privacy: contentEvent.PostPrivacyPayload{
			Scope:        data.Privacy,
			AllowComment: data.AllowComment,
			AllowShare:   data.AllowShare,
		},
		Context: &contentEvent.PostContextPayload{
			TargetID: entity.ID.Hex(),
			Type:     sharedEnums.ContextTypeReel,
		},
		Hashtags: data.Hashtags,
		Mentions: data.Mentions,
		Media: []contentEvent.MediaItemPayload{
			{
				MediaType:    sharedEnums.MediaTypeVideo,
				URL:          data.Video.URL,
				ThumbnailURL: data.Video.ThumbnailURL,
				Width:        data.Video.Width,
				Height:       data.Video.Height,
				Duration:     data.Video.Duration,
				SizeBytes:    data.Video.SizeBytes,
				MimeType:     data.Video.MimeType,
				TaggedUsers: func() []contentEvent.TaggedUserPayload {
					var taggedUsers []contentEvent.TaggedUserPayload
					for _, mention := range data.Mentions {
						taggedUsers = append(taggedUsers, contentEvent.TaggedUserPayload{
							UserID: mention,
						})
					}
					return taggedUsers
				}(), // Reel không hỗ trợ tag người dùng trong video, chỉ tag ở caption
			},
		},
	}
	err = c.events.Publish(ctx, constants.TopicReel.String(), data.UserID, constants.Created.String(), payloadpost)
	if err != nil {
		return fmt.Errorf("failed to publish post creation event for reel %s: %w", entity.ID.Hex(), err)
	}
	return nil
}
func (c *ConsumerReel) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.UpdateReelPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataReel, err := c.reelRepo.GetReelByID(ctx, data.ReelID)
	if err != nil {
		return fmt.Errorf("failed to get existing reel from repository for event %s: %w", event.ID, err)
	}
	if dataReel == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("reel with ID %s not found for event %s", data.ReelID, event.ID))
	}
	updatedEntity := mapper.UpdateReelEntity(dataReel, data)
	err = c.reelRepo.UpdateReel(ctx, updatedEntity)
	if err != nil {
		return fmt.Errorf("failed to update reel in repository for event %s: %w", event.ID, err)
	}
	payloadpost := &contentEvent.UpdatePostPayload{
		PostID:  data.ReelID,
		Content: data.Caption,
		Privacy: &contentEvent.UpdatePrivacyPayload{
			Scope:        data.Privacy,
			AllowComment: data.AllowComment,
			AllowShare:   data.AllowShare,
		},
		Hashtags: &data.Hashtags,
		Mentions: &data.Mentions,
	}
	err = c.events.Publish(ctx, constants.TopicReel.String(), data.UserID, constants.Updated.String(), payloadpost)
	if err != nil {
		return fmt.Errorf("failed to publish post update event for reel %s: %w", data.ReelID, err)
	}
	return nil
}
func (c *ConsumerReel) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.DeleteReelPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.reelRepo.DeleteReel(ctx, data.ReelID)
	if err != nil {
		return fmt.Errorf("failed to delete reel in repository for event %s: %w", event.ID, err)
	}
	payloadpost := &contentEvent.DeletePostPayload{
		PostID: data.ReelID,
	}
	err = c.events.Publish(ctx, constants.TopicReel.String(), data.UserID, constants.Deleted.String(), payloadpost)
	if err != nil {
		return fmt.Errorf("failed to publish post deletion event for reel %s: %w", data.ReelID, err)
	}
	return nil
}
func (c *ConsumerReel) ConsumerFailedReel(ctx context.Context) error {
	return nil
}
