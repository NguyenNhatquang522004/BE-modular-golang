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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

type ConsumerEntityReaction struct {
	events             events.EventBus
	pool               IRepositoryShare.IWorkerPool
	redisRepo          IRepositoryShare.IRedis
	entityReactionRepo IRepositoryCassandra.IReactionsRepository
}

func NewConsumerEntityReaction(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, entityReactionRepo IRepositoryCassandra.IReactionsRepository) *ConsumerEntityReaction {
	return &ConsumerEntityReaction{
		events:             events,
		pool:               pool,
		redisRepo:          redisRepo,
		entityReactionRepo: entityReactionRepo,
	}
}
func (c *ConsumerEntityReaction) ConsumerEntityReaction(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicEntityReaction.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- err
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
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreatedEntityReaction(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedEntityReaction(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedEntityReaction(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
				}
			})
			if processErr != nil {
				errchan <- processErr
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchan <- err
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
				continue
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
func (c *ConsumerEntityReaction) handleCreatedEntityReaction(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.EntityReactionPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	entity := &entity.EntityReaction{
		TargetID:     data.TargetID,
		UserID:       data.UserID,
		TargetType:   data.TargetType,
		ReactionCode: data.ReactionCode,
		CreatedAt:    data.CreatedAt,
	}
	err = c.entityReactionRepo.CreateReaction(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create entity reaction for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerEntityReaction) handleUpdatedEntityReaction(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.EntityReactionPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	entity := &entity.EntityReaction{
		TargetID:     data.TargetID,
		UserID:       data.UserID,
		TargetType:   data.TargetType,
		ReactionCode: data.ReactionCode,
		CreatedAt:    data.CreatedAt,
	}
	err = c.entityReactionRepo.CreateReaction(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create entity reaction for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerEntityReaction) handleDeletedEntityReaction(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.EntityReactionPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	err = c.entityReactionRepo.DeleteReaction(ctx, data.TargetID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete entity reaction for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerEntityReaction) ConsumerFailedEntityReaction(ctx context.Context) error {
	return nil
}
