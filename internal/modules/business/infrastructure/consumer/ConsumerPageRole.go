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

type ConsumerPageRole struct {
	events       events.EventBus
	pool         IRepositoryShare.IWorkerPool
	redisRepo    IRepositoryShare.IRedis
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
}

func NewConsumerPageRole(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, pageRoleRepo IRepositoryMongodb.IPageRolesRepository) *ConsumerPageRole {
	return &ConsumerPageRole{
		events:       events,
		pool:         pool,
		redisRepo:    redisRepo,
		pageRoleRepo: pageRoleRepo,
	}
}
func (c *ConsumerPageRole) ConsumerPage(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicPageRole.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Đảm bảo unlock nếu không thể xử lý
				wg.Done()                         // Giảm counter ngay vì không có goroutine nào được tạo ra
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Đảm bảo unlock nếu có lỗi trong quá trình xử lý

			} else {
				errchan <- nil // Xử lý thành công, gửi nil vào channel
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
func (c *ConsumerPageRole) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.CreatePageRolePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	pageRoleEntity, err := mapper.ToPageRoleEntity(data)
	if err != nil {
		return fmt.Errorf("failed to convert payload to entity: %w", err)
	}
	err = c.pageRoleRepo.CreatePageRole(ctx, pageRoleEntity)
	if err != nil {
		return fmt.Errorf("failed to create page role in repository: %w", err)
	}
	return nil
}
func (c *ConsumerPageRole) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.UpdatePageRolePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	existingEntity, err := c.pageRoleRepo.GetPageRoleByID(ctx, data.PageRoleID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing page role: %w", err)
	}
	if existingEntity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("page role with ID %s not found", data.PageRoleID))
	}
	mapper.MapUpdateToPageRoleEntity(existingEntity, data)
	err = c.pageRoleRepo.UpdatePageRole(ctx, existingEntity)
	if err != nil {
		return fmt.Errorf("failed to update page role in repository: %w", err)
	}
	return nil
}
func (c *ConsumerPageRole) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[businessEvent.DeletedPageRolePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.DeLeteAll == true {
		err := c.pageRoleRepo.DeletePageRoleByPageID(ctx, data.PageRoleID)
		if err != nil {
			return fmt.Errorf("failed to delete all page roles by page ID in repository: %w", err)
		}
	} else {
		err = c.pageRoleRepo.DeletePageRoleByUserIDandPageID(ctx, data.UserID, data.PageRoleID)
		if err != nil {
			return fmt.Errorf("failed to delete page role in repository: %w", err)
		}
	}
	return nil
}
func (c *ConsumerPageRole) ConsumerFailedPage(ctx context.Context) error {

	return nil
}
