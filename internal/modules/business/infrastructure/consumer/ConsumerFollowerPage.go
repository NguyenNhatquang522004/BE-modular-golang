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

type ConsumerFollowerPage struct {
	events           events.EventBus
	redisRepo        IRepositoryShare.IRedis
	pageFollowerRepo IRepositoryMongodb.IPageFollowersRepository
	pool             IRepositoryShare.IWorkerPool
}

func NewConsumerFollowerPage(events events.EventBus, redisRepo IRepositoryShare.IRedis, pageFollowerRepo IRepositoryMongodb.IPageFollowersRepository, pool IRepositoryShare.IWorkerPool) *ConsumerFollowerPage {
	return &ConsumerFollowerPage{
		events:           events,
		redisRepo:        redisRepo,
		pageFollowerRepo: pageFollowerRepo,
		pool:             pool,
	}
}

func (c *ConsumerFollowerPage) ConsumerFollowerPage(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicFollowerPage.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					return
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
func (c *ConsumerFollowerPage) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.CreatePageFollowerPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity, err := mapper.ToEntityFromCreatePayload(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	err = c.pageFollowerRepo.CreateFollower(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create page follower for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerFollowerPage) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.UpdatePageFollowerPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	existing, err := c.pageFollowerRepo.GetFollowerByPageIDAndUserID(ctx, data.PageID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing page follower for event %s: %w", event.ID, err)
	}
	if existing == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("page follower not found for update in event %s", event.ID))
	}
	updatedEntity := mapper.UpdateEntityFromPayload(existing, data)
	err = c.pageFollowerRepo.UpdateFollower(ctx, updatedEntity)
	if err != nil {
		return fmt.Errorf("failed to update page follower for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerFollowerPage) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.DeletedPageFollowerPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.DeleteAll {
		// Xóa tất cả followers của user này trên page
		err := c.pageFollowerRepo.DeleteFollowerByPageID(ctx, data.PageID)
		if err != nil {
			return fmt.Errorf("failed to delete all followers for page %s in event %s: %w", data.PageID, event.ID, err)
		}
	} else {
		err := c.pageFollowerRepo.DeleteFollower(ctx, data.PageID, data.UserID)
		if err != nil {
			return fmt.Errorf("failed to delete follower for page %s and user %s in event %s: %w", data.PageID, data.UserID, event.ID, err)
		}
	}
	return nil
}
func (c *ConsumerFollowerPage) ConsumerFailedFollowerPage(ctx context.Context) error
