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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type ConsumerStatsConversation struct {
	events           events.EventBus
	pool             IRepositoryShare.IWorkerPool
	redisRepo        IRepositoryShare.IRedis
	participantRepo  IRepositoryMongodb.IConversationParticipantsRepository
	conversationRepo IRepositoryMongodb.IConversationsRepository
}

func NewConsumerStatsConversation(events events.EventBus,
	pool IRepositoryShare.IWorkerPool,
	redisRepo IRepositoryShare.IRedis,
	participantRepo IRepositoryMongodb.IConversationParticipantsRepository,
	conversationRepo IRepositoryMongodb.IConversationsRepository) *ConsumerStatsConversation {
	return &ConsumerStatsConversation{
		events:           events,
		pool:             pool,
		redisRepo:        redisRepo,
		participantRepo:  participantRepo,
		conversationRepo: conversationRepo,
	}
}
func (c *ConsumerStatsConversation) ConsumerStatsConversation(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicStatsConversation.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedStatsConversation(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedStatsConversation(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedStatsConversation(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
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
func (c *ConsumerStatsConversation) handleCreatedStatsConversation(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communicationEvent.ConversationStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.ParticipantCount != nil || data.LastMessage != nil {
		dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to get conversation from database: %w", err)
		}
		if dataconversation == nil {
			return fmt.Errorf("conversation not found for ID: %s", data.ConversationID)
		}
		if data.ParticipantCount != nil {
			dataconversation.ParticipantCount = dataconversation.ParticipantCount + *data.ParticipantCount
		}
		if data.LastMessage != nil {
			dataconversation.LastMessage.MessageID = data.LastMessage.MessageID
			dataconversation.LastMessage.Content = data.LastMessage.Content
			dataconversation.LastMessage.SenderID = data.LastMessage.SenderID
			dataconversation.LastMessage.Type = data.LastMessage.Type
			dataconversation.LastMessage.CreatedAt = data.LastMessage.CreatedAt
		}
	}
	if data.LastSeenAt != nil || data.LastSeenMessageID != nil && data.UserID != nil {
		participant, err := c.participantRepo.GetConversationParticipant(ctx, data.ConversationID, *data.UserID)
		if err != nil {
			return fmt.Errorf("failed to get conversation participant from database: %w", err)
		}
		if participant == nil {
			return fmt.Errorf("conversation participant not found for ConversationID: %s and UserID: %s", data.ConversationID, data.UserID)
		}
		if data.LastSeenAt != nil {
			participant.LastSeenAt = *data.LastSeenAt
		}
		if data.LastSeenMessageID != nil {
			participant.LastSeenMessageID = *data.LastSeenMessageID
		}
		err = c.participantRepo.UpdateConversationParticipant(ctx, participant)
		if err != nil {
			return fmt.Errorf("failed to update conversation participant in database: %w", err)
		}
	}
	return nil
}
func (c *ConsumerStatsConversation) handleUpdatedStatsConversation(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communicationEvent.ConversationStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.ParticipantCount != nil || data.LastMessage != nil {
		dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to get conversation from database: %w", err)
		}
		if dataconversation == nil {
			return fmt.Errorf("conversation not found for ID: %s", data.ConversationID)
		}
		if data.ParticipantCount != nil {
			dataconversation.ParticipantCount = dataconversation.ParticipantCount + *data.ParticipantCount
		}
		if data.LastMessage != nil {
			dataconversation.LastMessage.MessageID = data.LastMessage.MessageID
			dataconversation.LastMessage.Content = data.LastMessage.Content
			dataconversation.LastMessage.SenderID = data.LastMessage.SenderID
			dataconversation.LastMessage.Type = data.LastMessage.Type
			dataconversation.LastMessage.CreatedAt = data.LastMessage.CreatedAt
		}
	}
	if data.LastSeenAt != nil || data.LastSeenMessageID != nil && data.UserID != nil {
		participant, err := c.participantRepo.GetConversationParticipant(ctx, data.ConversationID, *data.UserID)
		if err != nil {
			return fmt.Errorf("failed to get conversation participant from database: %w", err)
		}
		if participant == nil {
			return fmt.Errorf("conversation participant not found for ConversationID: %s and UserID: %s", data.ConversationID, data.UserID)
		}
		if data.LastSeenAt != nil {
			participant.LastSeenAt = *data.LastSeenAt
		}
		if data.LastSeenMessageID != nil {
			participant.LastSeenMessageID = *data.LastSeenMessageID
		}
		err = c.participantRepo.UpdateConversationParticipant(ctx, participant)
		if err != nil {
			return fmt.Errorf("failed to update conversation participant in database: %w", err)
		}
	}
	return nil
}
func (c *ConsumerStatsConversation) handleDeletedStatsConversation(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communicationEvent.ConversationStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.ParticipantCount != nil || data.LastMessage != nil {
		dataconversation, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to get conversation from database: %w", err)
		}
		if dataconversation == nil {
			return fmt.Errorf("conversation not found for ID: %s", data.ConversationID)
		}
		if data.ParticipantCount != nil {
			dataconversation.ParticipantCount = dataconversation.ParticipantCount - *data.ParticipantCount
		}
		if data.LastMessage != nil {
			dataconversation.LastMessage.MessageID = data.LastMessage.MessageID
			dataconversation.LastMessage.Content = data.LastMessage.Content
			dataconversation.LastMessage.SenderID = data.LastMessage.SenderID
			dataconversation.LastMessage.Type = data.LastMessage.Type
			dataconversation.LastMessage.CreatedAt = data.LastMessage.CreatedAt
		}
	}
	if data.LastSeenAt != nil || data.LastSeenMessageID != nil && data.UserID != nil {
		participant, err := c.participantRepo.GetConversationParticipant(ctx, data.ConversationID, *data.UserID)
		if err != nil {
			return fmt.Errorf("failed to get conversation participant from database: %w", err)
		}
		if participant == nil {
			return fmt.Errorf("conversation participant not found for ConversationID: %s and UserID: %s", data.ConversationID, *data.UserID)
		}
		if data.LastSeenAt != nil {
			participant.LastSeenAt = *data.LastSeenAt
		}
		if data.LastSeenMessageID != nil {
			participant.LastSeenMessageID = *data.LastSeenMessageID
		}
		err = c.participantRepo.UpdateConversationParticipant(ctx, participant)
		if err != nil {
			return fmt.Errorf("failed to update conversation participant in database: %w", err)
		}
	}
	return nil
}
func (c *ConsumerStatsConversation) ConsumerFailedStatsConversation(ctx context.Context) error {

	return nil
}
