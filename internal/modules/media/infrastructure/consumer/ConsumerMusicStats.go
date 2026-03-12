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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerMusicStats struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	musicRepo IRepositoryMongodb.IMusicLibraryRepository
}

func NewConsumerMusicStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, musicRepo IRepositoryMongodb.IMusicLibraryRepository) *ConsumerMusicStats {
	return &ConsumerMusicStats{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		musicRepo: musicRepo,
	}
}

func (c *ConsumerMusicStats) ConsumerMusicStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicMusicStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		// Xử lý batch events ở đây
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
				errchan <- errors.New("failed to run event " + event.ID + " in worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Đảm bảo unlock nếu không thể xử lý
				wg.Done()                         // Giảm wg nếu không thể xử lý
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Đảm bảo unlock nếu xử lý thất bại
			} else {
				errchan <- nil // Gửi nil nếu xử lý thành công để biết rằng event này đã được xử lý
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
func (c *ConsumerMusicStats) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic cho event CREATED ở đây
	data, err := utils.ParsePayload[mediaEvent.MusicStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataMusic, err := c.musicRepo.GetMusicLibraryByID(ctx, data.ID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get music library by ID for event %s: %w", event.ID, err))
	}
	if dataMusic == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("music library not found for event %s", event.ID))
	}
	dataMusic.UsageCount = dataMusic.UsageCount + data.UsageCount
	_, err = c.musicRepo.UpdateMusicLibrary(ctx, dataMusic)
	if err != nil {
		return fmt.Errorf("failed to update music library for event %s: %w", event.ID, err)
	}
	payload := &mediaEvent.ArtistStatsPayload{
		ArtistID:      data.ArtistID,
		FollowerCount: 0,               // Tạm thời chưa có logic cập nhật follower count, có thể mở rộng sau
		TotalStreams:  data.UsageCount, // Cập nhật tổng số lần nhạc được dùng
	}
	err = c.events.Publish(ctx, constants.TopicArtistStats.String(), data.ArtistID, constants.Updated.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish artist stats event for artist %s: %w", data.ArtistID, err)
	}
	return nil
}

func (c *ConsumerMusicStats) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic cho event CREATED ở đây
	data, err := utils.ParsePayload[mediaEvent.MusicStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataMusic, err := c.musicRepo.GetMusicLibraryByID(ctx, data.ID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get music library by ID for event %s: %w", event.ID, err))
	}
	if dataMusic == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("music library not found for event %s", event.ID))
	}
	dataMusic.UsageCount = dataMusic.UsageCount + data.UsageCount
	_, err = c.musicRepo.UpdateMusicLibrary(ctx, dataMusic)
	if err != nil {
		return fmt.Errorf("failed to update music library for event %s: %w", event.ID, err)
	}
	payload := &mediaEvent.ArtistStatsPayload{
		ArtistID:      data.ArtistID,
		FollowerCount: 0,               // Tạm thời chưa có logic cập nhật follower count, có thể mở rộng sau
		TotalStreams:  data.UsageCount, // Cập nhật tổng số lần nhạc được dùng
	}
	err = c.events.Publish(ctx, constants.TopicArtistStats.String(), data.ArtistID, constants.Updated.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish artist stats event for artist %s: %w", data.ArtistID, err)
	}
	return nil
}

func (c *ConsumerMusicStats) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic cho event DELETED ở đây
	data, err := utils.ParsePayload[mediaEvent.MusicStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataMusic, err := c.musicRepo.GetMusicLibraryByID(ctx, data.ID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get music library by ID for event %s: %w", event.ID, err))
	}
	if dataMusic == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("music library not found for event %s", event.ID))
	}
	dataMusic.UsageCount = dataMusic.UsageCount - data.UsageCount
	_, err = c.musicRepo.UpdateMusicLibrary(ctx, dataMusic)
	if err != nil {
		return fmt.Errorf("failed to update music library for event %s: %w", event.ID, err)
	}
	payload := &mediaEvent.ArtistStatsPayload{
		ArtistID:      data.ArtistID,
		FollowerCount: 0,               // Tạm thời chưa có logic cập nhật follower count, có thể mở rộng sau
		TotalStreams:  data.UsageCount, // Cập nhật tổng số lần nhạc được dùng
	}
	err = c.events.Publish(ctx, constants.TopicArtistStats.String(), data.ArtistID, constants.Deleted.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish artist stats event for artist %s: %w", data.ArtistID, err)
	}
	return nil
}
func (c *ConsumerMusicStats) ConsumerFailedMusicStats(ctx context.Context) error {
	return nil
}
