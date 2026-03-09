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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type ConsumerEditTarget struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	editRepo  IRepositoryMongoDB.IEditLogsRepository
}

func NewConsumerEditTarget(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, editRepo IRepositoryMongoDB.IEditLogsRepository) *ConsumerEditTarget {
	return &ConsumerEditTarget{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		editRepo:  editRepo,
	}
}
func (c *ConsumerEditTarget) ConsumerEditTarget(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicEdit.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedEditTarget(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedEditTarget(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedEditTarget(ctx, event)
				default:
					processErr = errors.New("unknown event type: " + event.Type)
				}
			})
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
				continue
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
func (c *ConsumerEditTarget) handleCreatedEditTarget(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.EditTargetPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataversion, err := c.editRepo.GetLatestEditLogByTargetID(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get latest edit log for target %s: %w", data.TargetID, err)
	}
	entity, err := mapper.ToEntityEditLogsPayload(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for target %s: %w", data.TargetID, err))
	}
	entity.Version = dataversion.Version + 1
	err = c.editRepo.CreateEditLog(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create edit log for target %s: %w", data.TargetID, err)
	}
	return nil
}
func (c *ConsumerEditTarget) handleUpdatedEditTarget(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.EditTargetPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataversion, err := c.editRepo.GetLatestEditLogByTargetID(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get latest edit log for target %s: %w", data.TargetID, err)
	}
	entity, err := mapper.ToEntityEditLogsPayload(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for target %s: %w", data.TargetID, err))
	}
	entity.Version = dataversion.Version + 1
	err = c.editRepo.CreateEditLog(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create edit log for target %s: %w", data.TargetID, err)
	}
	return nil
}
func (c *ConsumerEditTarget) handleDeletedEditTarget(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.EditTargetPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.editRepo.DeleteEditLogsByTargetID(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to delete edit logs for target %s: %w", data.TargetID, err)
	}
	return nil
}
func (c *ConsumerEditTarget) ConsumerFailEditTarget(ctx context.Context) error {
	return nil
}
