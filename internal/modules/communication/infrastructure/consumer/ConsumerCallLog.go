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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type ConsumerCallLog struct {
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	callLogRepo IRepositoryMongodb.ICallLogsRepository
}

func NewConsumerCallLog(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, callLogRepo IRepositoryMongodb.ICallLogsRepository) *ConsumerCallLog {
	return &ConsumerCallLog{
		events:      events,
		pool:        pool,
		redisRepo:   redisRepo,
		callLogRepo: callLogRepo,
	}
}

func (c *ConsumerCallLog) ConsumerCallLogEvents(ctx context.Context) error {
	// Implement the logic for consuming call log events
	err := c.events.SubscribeBatch(ctx, constants.TopicCallLog.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = errors.New("unsupported event type " + event.Type + " for event " + event.ID)
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
func (c *ConsumerCallLog) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created call log events
	data, err := utils.ParsePayload[communicationEvent.CreateCallLogPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event " + event.ID + ": " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event " + event.ID))
	}
	entity, err := mapper.ToCallLogEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event " + event.ID + ": " + err.Error()))
	}
	err = c.callLogRepo.CreateCallLog(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create call log in repository for event " + event.ID + ": " + err.Error())
	}
	return nil
}

func (c *ConsumerCallLog) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated call log events
	data, err := utils.ParsePayload[communicationEvent.UpdateCallLogPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event " + event.ID + ": " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event " + event.ID))
	}
	existingEntity, err := c.callLogRepo.GetCallLogByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing call log from repository for event " + event.ID + ": " + err.Error())
	}
	if existingEntity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("call log not found with ID " + data.ID + " for event " + event.ID))
	}
	updatedEntity := mapper.MapUpdateCallLogEntity(existingEntity, data)
	err = c.callLogRepo.UpdateCallLog(ctx, updatedEntity)
	if err != nil {
		return fmt.Errorf("failed to update call log in repository for event " + event.ID + ": " + err.Error())
	}
	return nil
}

func (c *ConsumerCallLog) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted call log events
	data, err := utils.ParsePayload[communicationEvent.DeleteCallLogPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event " + event.ID + ": " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event " + event.ID))
	}
	if data.DeleteAll {
		err = c.callLogRepo.DeleteCallLogByConversationID(ctx, data.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to delete call logs by conversation ID in repository for event " + event.ID + ": " + err.Error())
		}
	} else {
		err = c.callLogRepo.DeleteCallLog(ctx, data.ID)
		if err != nil {
			return fmt.Errorf("failed to delete call log in repository for event " + event.ID + ": " + err.Error())
		}
	}
	return nil
}
func (c *ConsumerCallLog) ConsumerFailedCallLogEvents(ctx context.Context) error {
	// Implement the logic for consuming failed call log events
	return nil
}
