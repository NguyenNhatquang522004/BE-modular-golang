package IConsumer

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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerArtistStats struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	artisRepo IRepositoryMongodb.IArtistRepository
}

func NewConsumerArtistStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, artisRepo IRepositoryMongodb.IArtistRepository) *ConsumerArtistStats {
	return &ConsumerArtistStats{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		artisRepo: artisRepo,
	}
}
func (c *ConsumerArtistStats) ConsumerArtistStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicArtistStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					return
				}
			})
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				c.redisRepo.MarkCompleted(ctx, event.ID)
				errchan <- nil
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
func (c *ConsumerArtistStats) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.ArtistStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse ArtistStatsPayload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataartis, err := c.artisRepo.GetArtistByID(ctx, data.ArtistID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get artist by ID %s for event %s: %w", data.ArtistID, event.ID, err))
	}
	if dataartis == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("artist not found with ID %s for event %s", data.ArtistID, event.ID))
	}
	dataartis.FollowerCount = dataartis.FollowerCount + data.FollowerCount
	dataartis.TotalStreams = dataartis.TotalStreams + data.TotalStreams
	err = c.artisRepo.UpdatedArtist(ctx, dataartis)
	if err != nil {
		return fmt.Errorf("failed to update artist stats for artist ID %s for event %s: %w", data.ArtistID, event.ID, err)
	}
	return nil
}
func (c *ConsumerArtistStats) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.ArtistStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse ArtistStatsPayload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataartis, err := c.artisRepo.GetArtistByID(ctx, data.ArtistID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get artist by ID %s for event %s: %w", data.ArtistID, event.ID, err))
	}
	if dataartis == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("artist not found with ID %s for event %s", data.ArtistID, event.ID))
	}
	dataartis.FollowerCount = dataartis.FollowerCount + data.FollowerCount
	dataartis.TotalStreams = dataartis.TotalStreams + data.TotalStreams
	err = c.artisRepo.UpdatedArtist(ctx, dataartis)
	if err != nil {
		return fmt.Errorf("failed to update artist stats for artist ID %s for event %s: %w", data.ArtistID, event.ID, err)
	}
	return nil

}
func (c *ConsumerArtistStats) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.ArtistStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse ArtistStatsPayload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataartis, err := c.artisRepo.GetArtistByID(ctx, data.ArtistID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get artist by ID %s for event %s: %w", data.ArtistID, event.ID, err))
	}
	if dataartis == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("artist not found with ID %s for event %s", data.ArtistID, event.ID))
	}
	dataartis.FollowerCount = dataartis.FollowerCount - data.FollowerCount
	dataartis.TotalStreams = dataartis.TotalStreams - data.TotalStreams
	err = c.artisRepo.UpdatedArtist(ctx, dataartis)
	if err != nil {
		return fmt.Errorf("failed to update artist stats for artist ID %s for event %s: %w", data.ArtistID, event.ID, err)
	}
	return nil
}
func (c *ConsumerArtistStats) ConsumerFailedArtistStats(ctx context.Context) error {

	return nil
}
