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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerGroupEvent struct {
	events         events.EventBus
	pool           IRepositoryShare.IWorkerPool
	redisRepo      IRepositoryShare.IRedis
	groupEventRepo IRepositoryMongodb.IGroupEventsRepository
}

func NewConsumerGroupEvent(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupEventRepo IRepositoryMongodb.IGroupEventsRepository) *ConsumerGroupEvent {
	return &ConsumerGroupEvent{
		events:         events,
		pool:           pool,
		redisRepo:      redisRepo,
		groupEventRepo: groupEventRepo,
	}
}

func (c *ConsumerGroupEvent) ConsumerGroupEvent(ctx context.Context) error {
	// Implement the logic for consuming group events
	err := c.events.SubscribeBatch(ctx, constants.TopicGroupEvent.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
			var processErr error
			wg.Add(1)
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreatedGroupEvent(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedGroupEvent(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedGroupEvent(ctx, event)
				default:
					processErr = errors.New("unknown event type: " + event.Type)
					return
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Decrement the WaitGroup counter since the task won't be processed
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
func (c *ConsumerGroupEvent) handleCreatedGroupEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created group events
	data, err := utils.ParsePayload[communityEvent.CreateGroupEventPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	entity, err := mapper.ToGroupEventEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	if entity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("mapped entity is nil for event %s", event.ID))
	}
	err = c.groupEventRepo.CreateGroupEvent(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create group event in repository for event %s: %w", event.ID, err)
	}
	// Use the parsed data to perform necessary actions
	return nil
}

func (c *ConsumerGroupEvent) handleUpdatedGroupEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated group
	data, err := utils.ParsePayload[communityEvent.UpdateGroupEventPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataEvent, err := c.groupEventRepo.GetGroupEventByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to get old group event in repository for event %s: %w", event.ID, err)
	}
	if dataEvent == nil {
		return fmt.Errorf("no existing group event found with ID %s for event %s", data.ID, event.ID)
	}
	mapper.ApplyUpdateToGroupEvent(dataEvent, data)
	err = c.groupEventRepo.CreateGroupEvent(ctx, dataEvent)
	if err != nil {
		return fmt.Errorf("failed to create updated group event in repository for event %s: %w", event.ID, err)
	}
	return nil
}

func (c *ConsumerGroupEvent) handleDeletedGroupEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted group events
	data, err := utils.ParsePayload[communityEvent.DeleteGroupEventPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	if data.DeleteAll == true {
		err = c.groupEventRepo.DeleteGroupEventsByGroupID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to delete all group events by group ID in repository for event %s: %w", event.ID, err)
		}
	} else {
		err = c.groupEventRepo.DeleteGroupEvent(ctx, data.EventID)
		if err != nil {
			return fmt.Errorf("failed to delete group event in repository for event %s: %w", event.ID, err)
		}
	}
	return nil
}
func (c *ConsumerGroupEvent) ConsumerFailedGroupEvent(ctx context.Context) error {
	// Implement the logic for handling failed group events
	return nil
}
