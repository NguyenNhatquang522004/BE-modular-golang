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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type ConsumerStats struct {
	events    events.EventBus
	redisRepo IRepositoryShare.IRedis
	pool      IRepositoryShare.IWorkerPool
	pageRepo  IRepositoryMongodb.IPagesRepository
}

func NewConsumerStats(events events.EventBus, redisRepo IRepositoryShare.IRedis, pool IRepositoryShare.IWorkerPool, pageRepo IRepositoryMongodb.IPagesRepository) *ConsumerStats {
	return &ConsumerStats{
		events:    events,
		redisRepo: redisRepo,
		pool:      pool,
		pageRepo:  pageRepo,
	}
}

func (c *ConsumerStats) ConsumerStatsPage(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicStatsPage.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
func (c *ConsumerStats) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.StatsPagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataStats, err := c.pageRepo.GetPageByID(ctx, data.PageID)
	if err != nil {
		return fmt.Errorf("failed to get page data for event %s: %w", event.ID, err)
	}
	dataStats.Stats.FollowersCount = dataStats.Stats.FollowersCount + data.FollowersCount
	dataStats.Stats.LikesCount = dataStats.Stats.LikesCount + data.LikesCount
	dataStats.Stats.RatingScore = dataStats.Stats.RatingScore + data.RatingScore
	dataStats.Stats.ReviewCount = dataStats.Stats.ReviewCount + data.ReviewCount
	err = c.pageRepo.UpdatePage(ctx, dataStats)
	if err != nil {
		return fmt.Errorf("failed to update page stats for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerStats) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.StatsPagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataStats, err := c.pageRepo.GetPageByID(ctx, data.PageID)
	if err != nil {
		return fmt.Errorf("failed to get page data for event %s: %w", event.ID, err)
	}
	dataStats.Stats.FollowersCount = dataStats.Stats.FollowersCount + data.FollowersCount
	dataStats.Stats.LikesCount = dataStats.Stats.LikesCount + data.LikesCount
	dataStats.Stats.RatingScore = dataStats.Stats.RatingScore + data.RatingScore
	dataStats.Stats.ReviewCount = dataStats.Stats.ReviewCount + data.ReviewCount
	err = c.pageRepo.UpdatePage(ctx, dataStats)
	if err != nil {
		return fmt.Errorf("failed to update page stats for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerStats) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.StatsPagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataStats, err := c.pageRepo.GetPageByID(ctx, data.PageID)
	if err != nil {
		return fmt.Errorf("failed to get page data for event %s: %w", event.ID, err)
	}
	dataStats.Stats.FollowersCount = dataStats.Stats.FollowersCount - data.FollowersCount
	dataStats.Stats.LikesCount = dataStats.Stats.LikesCount - data.LikesCount
	dataStats.Stats.RatingScore = dataStats.Stats.RatingScore - data.RatingScore
	dataStats.Stats.ReviewCount = dataStats.Stats.ReviewCount - data.ReviewCount
	err = c.pageRepo.UpdatePage(ctx, dataStats)
	if err != nil {
		return fmt.Errorf("failed to update page stats for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerStats) ConsumerFailedStatsPage(ctx context.Context) error {
	return nil
}
