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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type ConsumerPage struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	pageRepo  IRepositoryMongodb.IPagesRepository
}

func NewConsumerPage(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, pageRepo IRepositoryMongodb.IPagesRepository) *ConsumerPage {
	return &ConsumerPage{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		pageRepo:  pageRepo,
	}
}

func (c *ConsumerPage) ConsumerPage(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicPage.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.hnadleCreatedPage(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedPage(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedPage(ctx, event)
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

func (c *ConsumerPage) hnadleCreatedPage(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.CreatePagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	page := mapper.MapCreatePayloadToEntity(data)
	if page == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity"))
	}
	err = c.pageRepo.CreatePage(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to create page in repository: %w", err)
	}
	return nil
}

func (c *ConsumerPage) handleUpdatedPage(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.UpdatePagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	page, err := c.pageRepo.GetPageByID(ctx, data.PageID)
	if err != nil {
		return fmt.Errorf("failed to get page by ID: %w", err)
	}
	if page == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("page not found with ID: %s", data.PageID))
	}
	mapper.ApplyUpdatePayloadToEntity(page, data)
	err = c.pageRepo.UpdatePage(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to update page in repository: %w", err)
	}
	return nil
}

func (c *ConsumerPage) handleDeletedPage(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.DeletedPagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.pageRepo.DeletePage(ctx, data.PageID)
	if err != nil {
		return fmt.Errorf("failed to delete page in repository: %w", err)
	}
	return nil
}
func (c *ConsumerPage) ConsumerFailedPage(ctx context.Context) error {
	return nil
}
