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

type ConsumerConversationParticipant struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	participantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewConsumerConversationParticipant(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, participantRepo IRepositoryMongodb.IConversationParticipantsRepository) *ConsumerConversationParticipant {
	return &ConsumerConversationParticipant{
		events:          events,
		pool:            pool,
		redisRepo:       redisRepo,
		participantRepo: participantRepo,
	}
}

func (c *ConsumerConversationParticipant) ConsumerConversationParticipantEvents(ctx context.Context) error {
	// Implement the logic for consuming conversation participant events
	err := c.events.SubscribeBatch(ctx, constants.TopicConversationParticipant.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
func (c *ConsumerConversationParticipant) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created conversation participant events
	data, err := utils.ParsePayload[communicationEvent.CreateParticipantPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	entity, err := mapper.MapCreatePayloadToEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	err = c.participantRepo.CreateConversationParticipant(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create conversation participant for event %s: %w", event.ID, err)
	}
	return nil
}

func (c *ConsumerConversationParticipant) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated conversation participant
	data, err := utils.ParsePayload[communicationEvent.UpdateParticipantPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	existingEntity, err := c.participantRepo.GetConversationParticipant(ctx, data.ConversationID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to get existing conversation participant for event %s: %w", event.ID, err)
	}
	if existingEntity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("conversation participant not found for event %s with conversation ID %s and user ID %s", event.ID, data.ConversationID, data.UserID))
	}
	mapper.ApplyUpdatePayloadToEntity(existingEntity, data)
	err = c.participantRepo.UpdateConversationParticipant(ctx, existingEntity)
	if err != nil {
		return fmt.Errorf("failed to update conversation participant for event %s: %w", event.ID, err)
	}
	return nil
}

func (c *ConsumerConversationParticipant) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted conversation participant events
	data, err := utils.ParsePayload[communicationEvent.DeleteParticipantPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	if data.DeleteAll == true {
		err = c.participantRepo.DeleteConversationParticipant(ctx, data.ConversationID, data.UserID)
		if err != nil {
			return fmt.Errorf("failed to delete conversation participant for event %s: %w", event.ID, err)
		}
		return nil
	}
	if data.UserID != "" {
		err = c.participantRepo.DeleteConversationParticipant(ctx, data.ConversationID, data.UserID)
		if err != nil {
			return fmt.Errorf("failed to delete conversation participant for event %s: %w", event.ID, err)
		}
	}

	return nil
}
func (c *ConsumerConversationParticipant) ConsumerFailedConversationParticipantEvents(ctx context.Context) error {
	// Implement the logic for consuming failed conversation participant events
	return nil
}
