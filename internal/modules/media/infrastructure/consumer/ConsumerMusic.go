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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerMusic struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	musicRepo IRepositoryMongodb.IMusicLibraryRepository
}

func NewConsumerMusic(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, musicRepo IRepositoryMongodb.IMusicLibraryRepository) *ConsumerMusic {
	return &ConsumerMusic{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		musicRepo: musicRepo,
	}
}

func (c *ConsumerMusic) ConsumerMusic(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicMusic.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = errors.New("unknown event type: " + event.Type)
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit task for event " + event.ID + ": " + err.Error())
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
func (c *ConsumerMusic) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.CreateMusicLibraryPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	albumEntity, err := mapper.ToMusicLibraryEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	_, err = c.musicRepo.CreateMusicLibrary(ctx, albumEntity)
	if err != nil {
		return fmt.Errorf("failed to create music library entry for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerMusic) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.UpdateMusicLibraryPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataalbum, err := c.musicRepo.GetMusicLibraryByID(ctx, data.MusicID)
	if err != nil {
		return fmt.Errorf("failed to get music library entry for event %s: %w", event.ID, err)
	}
	if dataalbum == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("music library entry not found for event %s", event.ID))
	}
	entity := mapper.ApplyUpdateMusicLibrary(dataalbum, data)
	_, err = c.musicRepo.UpdateMusicLibrary(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to update music library entry for event %s: %w", event.ID, err)
	}
	// Optionally, you can add logic here to invalidate cache or perform other post-update actions
	return nil
}
func (c *ConsumerMusic) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.DeleteMusicLibraryPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.ArtisID != "" {
		err = c.musicRepo.DeleteMusicLibraryByArtist(ctx, data.ArtisID)
		if err != nil {
			return fmt.Errorf("failed to delete music library entries for artist %s in event %s: %w", data.ArtisID, event.ID, err)
		}
		return nil
	}
	if data.MusicID != "" {
		err = c.musicRepo.DeleteMusicLibrary(ctx, data.MusicID)
		if err != nil {
			return fmt.Errorf("failed to delete music library entry for event %s: %w", event.ID, err)
		}
	} 
	return nil
}
func (c *ConsumerMusic) ConsumerFailedMusic(ctx context.Context) error {
	return nil
}
