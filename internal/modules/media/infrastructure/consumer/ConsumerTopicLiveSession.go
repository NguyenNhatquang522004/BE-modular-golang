package consumer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerTopicLiveSession struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	livesession   IRepositoryMongodb.ILiveSessionRepository
	seaweedfsRepo IRepositoryShare.ISeaweedfs
	cfg           *configs.Config
}

func NewConsumerTopicLiveSession(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, livesession IRepositoryMongodb.ILiveSessionRepository, seaweedfsRepo IRepositoryShare.ISeaweedfs, cfg *configs.Config) *ConsumerTopicLiveSession {
	return &ConsumerTopicLiveSession{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		livesession:   livesession,
		seaweedfsRepo: seaweedfsRepo,
		cfg:           cfg,
	}
}

func (c *ConsumerTopicLiveSession) ConsumerTopicLiveSession(ctx context.Context) error {
	// Implement the logic for consuming the live session topic
	err := c.events.SubscribeBatch(ctx, constants.TopicLiveSession.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				wg.Done()
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
func (c *ConsumerTopicLiveSession) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.CreateLiveSessionPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	streamkey := utils.GenerateSecureStreamKey(data.UserID, data.SessionID, c.cfg.LiveStream.SECRET_KEY, time.Duration(24)*time.Hour)
	var playbackURL string
	if data.PageID != nil {
		playbackURL = c.seaweedfsRepo.GetStreamURL(*data.PageID, data.SessionID, utils.BucketPageLiveStream)
	}
	if data.GroupID != nil {
		playbackURL = c.seaweedfsRepo.GetStreamURL(*data.GroupID, data.SessionID, utils.BucketGroupStream)
	}
	if data.PageID == nil && data.GroupID == nil {
		playbackURL = c.seaweedfsRepo.GetStreamURL(data.UserID, data.SessionID, utils.BucketLive)
	}
	entity := mapper.MapCreatePayloadToEntity(data, streamkey, playbackURL)
	err = c.livesession.CreateLiveSession(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create live session in database: %w", err)
	}
	// Implement the logic for handling the created event using the parsed data
	return nil
}

// Implement the logic for handling the created event using the parsed data
func (c *ConsumerTopicLiveSession) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.UpdateLiveSessionPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	existingLive, err := c.livesession.GetLiveSessionByID(ctx, data.SessionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing live session from database: %w", err)
	}
	if existingLive == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("live session with ID %s not found", data.SessionID))
	}
	updatedLive := mapper.ApplyUpdatePayloadToEntity(existingLive, data)
	if updatedLive == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("no valid fields to update for live session with ID %s", data.SessionID))
	}
	err = c.livesession.UpdateLiveSession(ctx, updatedLive)
	if err != nil {
		return fmt.Errorf("failed to update live session in database: %w", err)
	}
	return nil
}
func (c *ConsumerTopicLiveSession) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.DeleteLiveSessionPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.DeleteAll {
		if data.PageID != nil {
			err = c.livesession.DeleteLiveSessionsByPageID(ctx, *data.PageID)
			if err != nil {
				return fmt.Errorf("failed to delete live sessions by PageID in database: %w", err)
			}

		}
		if data.GroupID != nil {
			err = c.livesession.DeleteLiveSessionsByGroupID(ctx, *data.GroupID)
			if err != nil {
				return fmt.Errorf("failed to delete live sessions by GroupID in database: %w", err)
			}
		}
		if data.PageID == nil && data.GroupID == nil {
			err = c.livesession.DeleteLiveSessionsByHostUserID(ctx, *data.UserID)
			if err != nil {
				return fmt.Errorf("failed to delete live sessions by HostUserID in database: %w", err)
			}
		}
	} else {
		err = c.livesession.DeleteLiveSessionByID(ctx, data.SessionID)
		if err != nil {
			return fmt.Errorf("failed to delete live session by ID in database: %w", err)
		}
	}

	return nil
}
func (c *ConsumerTopicLiveSession) ConsumerFailedTopicLiveSession(ctx context.Context) error {
	// Implement the logic for consuming the failed live session topic
	return nil
}
