package consumer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type ConsumerConversation struct {
	events           events.EventBus
	pool             IRepositoryShare.IWorkerPool
	redisRepo        IRepositoryShare.IRedis
	conversationRepo IRepositoryMongodb.IConversationsRepository
}

func NewConsumerConversation(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, conversationRepo IRepositoryMongodb.IConversationsRepository) *ConsumerConversation {
	return &ConsumerConversation{
		events:           events,
		pool:             pool,
		redisRepo:        redisRepo,
		conversationRepo: conversationRepo,
	}
}

func (c *ConsumerConversation) ConsumerConversation(ctx context.Context) error {
	// Implement the logic for consuming conversation messages
	err := c.events.SubscribeBatch(ctx, constants.TopicConversation.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedConversation(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedConversation(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedConversation(ctx, event)
				default:
					processErr = errors.New("unknown event type: " + event.Type)
				}
			})
			if err != nil {
				errchan <- errors.New("failed to run task for event " + event.ID + ": " + err.Error())
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
func (c *ConsumerConversation) handleCreatedConversation(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created conversation messages
	data, err := utils.ParsePayload[communicationEvent.CreateConversationPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.Type == sharedEnums.TypePrivate {
		var users []string
		users = append(users, data.UserCreatorAndOwnerID)
		users = append(users, data.ParticipantIDs[0].UserID)
		sort.Strings(users)
		privateKey := fmt.Sprintf("private_%s_%s", users[0], users[1])
		existingConv, err := c.conversationRepo.CheckConversationExistsByPrivateChatKey(ctx, privateKey)
		if err != nil {
			return fmt.Errorf("failed to check existing conversation by private chat key: %w", err)
		}
		if existingConv != nil {
			return fmt.Errorf("conversation with private chat key %s already exists", privateKey)
		}
	}
	conversationEntity, err := mapper.ToConversationEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to conversation entity for event %s: %w", event.ID, err))
	}
	if conversationEntity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("mapped conversation entity is nil for event %s", event.ID))
	}
	err = c.conversationRepo.CreateConversation(ctx, conversationEntity)
	if err != nil {
		return fmt.Errorf("failed to create conversation in repository for event %s: %w", event.ID, err)
	}
	payloadCreated := &communicationEvent.CreateParticipantPayload{
		ConversationID: conversationEntity.ID.Hex(),
		UserID:         data.UserCreatorAndOwnerID,
		AddedByUserID:  data.UserCreatorAndOwnerID,
		Nickname:       data.UserCreatorName,
		Role:           sharedEnums.RoleTypeAdmin,
	}
	err = c.events.Publish(ctx, constants.TopicConversationParticipant.String(), conversationEntity.ID.Hex(), constants.Created.String(), payloadCreated)
	for _, participant := range data.ParticipantIDs {
		payload := &communicationEvent.CreateParticipantPayload{
			ConversationID: conversationEntity.ID.Hex(),
			UserID:         participant.UserID,
			AddedByUserID:  data.UserCreatorAndOwnerID,
			Nickname:       participant.Nickname,
			Role:           sharedEnums.RoleTypeMember,
		}
		err = c.events.Publish(ctx, constants.TopicConversationParticipant.String(), conversationEntity.ID.Hex(), constants.Created.String(), payload)
		if err != nil {
			return fmt.Errorf("failed to publish create participant event for user %s in conversation %s: %w", participant.UserID, conversationEntity.ID.Hex(), err)
		}
	}
	return nil
}

func (c *ConsumerConversation) handleUpdatedConversation(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated conversation messages
	data, err := utils.ParsePayload[communicationEvent.UpdateConversationReq](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	existingConv, err := c.conversationRepo.GetConversationByID(ctx, data.ConversationID)
	if err != nil {
		return fmt.Errorf("failed to get existing conversation by ID for event %s: %w", event.ID, err)
	}
	if existingConv == nil {
		return fmt.Errorf("conversation with ID %s not found for event %s", data.ConversationID, event.ID)
	}
	hasChanges, err := mapper.ApplyConversationUpdate(existingConv, data)
	if err != nil {
		return fmt.Errorf("failed to apply conversation update for event %s: %w", event.ID, err)
	}
	if !hasChanges {
		return nil
	}
	err = c.conversationRepo.UpdateConversation(ctx, existingConv)
	if err != nil {
		return fmt.Errorf("failed to update conversation in repository for event %s: %w", event.ID, err)
	}
	return nil
}

func (c *ConsumerConversation) handleDeletedConversation(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted conversation messages
	data, err := utils.ParsePayload[communicationEvent.DeletePrivateConversationGroupPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.conversationRepo.DeleteConversation(ctx, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation in repository for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerConversation) ConsumerFailedConversation(ctx context.Context) error {
	// Implement the logic for handling failed conversation messages
	return nil
}
