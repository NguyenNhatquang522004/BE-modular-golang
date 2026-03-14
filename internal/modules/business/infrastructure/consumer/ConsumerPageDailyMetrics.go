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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/gocql/gocql"
)

type ConsumerPageDailyMetrics struct {
	// Define the fields for the ConsumerPageDailyMetrics struct here
	events         events.EventBus
	pageMetricRepo IRepositoryCassandra.IPageDailyMetricsRepository
	pool           IRepositoryShare.IWorkerPool
	redisRepo      IRepositoryShare.IRedis
}

func NewConsumerPageDailyMetrics(events events.EventBus, pageMetricRepo IRepositoryCassandra.IPageDailyMetricsRepository, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis) *ConsumerPageDailyMetrics {
	return &ConsumerPageDailyMetrics{
		events:         events,
		pageMetricRepo: pageMetricRepo,
		pool:           pool,
		redisRepo:      redisRepo,
	}
}

func (c *ConsumerPageDailyMetrics) ConsumerPageDailyMetric(ctx context.Context) error {
	// Implement the logic for consuming page daily metrics here
	err := c.events.SubscribeBatch(ctx, constants.TopicDailyMetricsPage.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				errchan <- errors.New("failed to submit task for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Manually mark as done since the task was not submitted
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
func (c *ConsumerPageDailyMetrics) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created events here
	data, err := utils.ParsePayload[businessEvent.PageDailyMetricPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datametric, err := c.pageMetricRepo.GetPageDailyMetricLatest(ctx, data.PageID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get latest page daily metric: %w", err))
	}
	if datametric == nil {
		convertPageID, err := gocql.ParseUUID(data.PageID)
		normalizedDate := time.Date(
			data.MetricDate.Year(),
			data.MetricDate.Month(),
			data.MetricDate.Day(),
			0, 0, 0, 0, time.UTC,
		)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("invalid page_id format: %w", err))
		}
		entity := &entity.PageDailyMetric{
			ID:               convertPageID,
			MetricDate:       normalizedDate,
			ReachTotal:       *data.ReachTotal + 0,
			ReachPaid:        *data.ReachPaid + 0,
			ReachOrganic:     *data.ReachOrganic + 0,
			ImpressionsTotal: *data.ImpressionsTotal + 0,
			NewFollowers:     *data.NewFollowers + 0,
			Unfollows:        *data.Unfollows + 0,
			ProfileViews:     *data.ProfileViews + 0,
			WebsiteClicks:    *data.WebsiteClicks + 0,
		}
		err = c.pageMetricRepo.InsertOrUpdatePageDailyMetric(ctx, entity)
		if err != nil {
			return errors.New("failed to insert new page daily metric: " + err.Error())
		}
	} else {
		datametric.ReachTotal = *data.ReachTotal + datametric.ReachTotal
		datametric.ReachPaid = *data.ReachPaid + datametric.ReachPaid
		datametric.ReachOrganic = *data.ReachOrganic + datametric.ReachOrganic
		datametric.ImpressionsTotal = *data.ImpressionsTotal + datametric.ImpressionsTotal
		datametric.NewFollowers = *data.NewFollowers + datametric.NewFollowers
		datametric.Unfollows = *data.Unfollows + datametric.Unfollows
		datametric.ProfileViews = *data.ProfileViews + datametric.ProfileViews
		datametric.WebsiteClicks = *data.WebsiteClicks + datametric.WebsiteClicks
		err = c.pageMetricRepo.InsertOrUpdatePageDailyMetric(ctx, datametric)
		if err != nil {
			return errors.New("failed to update existing page daily metric: " + err.Error())
		}
	}
	return nil
}

func (c *ConsumerPageDailyMetrics) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created events here
	data, err := utils.ParsePayload[businessEvent.PageDailyMetricPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datametric, err := c.pageMetricRepo.GetPageDailyMetricLatest(ctx, data.PageID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get latest page daily metric: %w", err))
	}
	if datametric == nil {
		convertPageID, err := gocql.ParseUUID(data.PageID)
		normalizedDate := time.Date(
			data.MetricDate.Year(),
			data.MetricDate.Month(),
			data.MetricDate.Day(),
			0, 0, 0, 0, time.UTC,
		)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("invalid page_id format: %w", err))
		}
		entity := &entity.PageDailyMetric{
			ID:               convertPageID,
			MetricDate:       normalizedDate,
			ReachTotal:       *data.ReachTotal + 0,
			ReachPaid:        *data.ReachPaid + 0,
			ReachOrganic:     *data.ReachOrganic + 0,
			ImpressionsTotal: *data.ImpressionsTotal + 0,
			NewFollowers:     *data.NewFollowers + 0,
			Unfollows:        *data.Unfollows + 0,
			ProfileViews:     *data.ProfileViews + 0,
			WebsiteClicks:    *data.WebsiteClicks + 0,
		}
		err = c.pageMetricRepo.InsertOrUpdatePageDailyMetric(ctx, entity)
		if err != nil {
			return errors.New("failed to insert new page daily metric: " + err.Error())
		}
	} else {
		datametric.ReachTotal = *data.ReachTotal + datametric.ReachTotal
		datametric.ReachPaid = *data.ReachPaid + datametric.ReachPaid
		datametric.ReachOrganic = *data.ReachOrganic + datametric.ReachOrganic
		datametric.ImpressionsTotal = *data.ImpressionsTotal + datametric.ImpressionsTotal
		datametric.NewFollowers = *data.NewFollowers + datametric.NewFollowers
		datametric.Unfollows = *data.Unfollows + datametric.Unfollows
		datametric.ProfileViews = *data.ProfileViews + datametric.ProfileViews
		datametric.WebsiteClicks = *data.WebsiteClicks + datametric.WebsiteClicks
		err = c.pageMetricRepo.InsertOrUpdatePageDailyMetric(ctx, datametric)
		if err != nil {
			return errors.New("failed to update existing page daily metric: " + err.Error())
		}
	}
	return nil
}

func (c *ConsumerPageDailyMetrics) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted events here
	data, err := utils.ParsePayload[businessEvent.DeletePageDailyMetricPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.DeleteALL {
		err = c.pageMetricRepo.DeletePageDailyMetricsByPageID(ctx, data.PageID)
		if err != nil {
			return errors.New("failed to delete all page daily metrics for page: " + err.Error())
		}
	} else {
		normalizedDate := time.Date(
			data.MetricDate.Year(),
			data.MetricDate.Month(),
			data.MetricDate.Day(),
			0, 0, 0, 0, time.UTC,
		)
		err = c.pageMetricRepo.DeletePageDailyMetric(ctx, data.PageID, normalizedDate)
		if err != nil {
			return errors.New("failed to delete page daily metric: " + err.Error())
		}
	}
	return nil
}
func (c *ConsumerPageDailyMetrics) ConsumerFailedPageDailyMetric(ctx context.Context) error {
	// Implement the logic for consuming failed page daily metrics here
	return nil
}
