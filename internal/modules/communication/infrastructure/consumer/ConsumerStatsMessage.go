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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
)

type ConsumerStatsMessage struct {
	events              events.EventBus
	pool                IRepositoryShare.IWorkerPool
	redisRepo           IRepositoryShare.IRedis
	messageReactionRepo IRepositoryCassandra.IMessageReactionsRepository
}

func NewConsumerStatsMessage(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, messageReactionRepo IRepositoryCassandra.IMessageReactionsRepository) *ConsumerStatsMessage {
	return &ConsumerStatsMessage{
		events:              events,
		pool:                pool,
		redisRepo:           redisRepo,
		messageReactionRepo: messageReactionRepo,
	}
}

func (c *ConsumerStatsMessage) ConsumerStatsMessage(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicStatsMessage.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedStatsMessage(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedStatsMessage(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedStatsMessage(ctx, event)
				default:
					processErr = errors.New("unknown event type: " + event.Type)
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
func (c *ConsumerStatsMessage) handleCreatedStatsMessage(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communicationEvent.MessageStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataMessageReaction, err := c.messageReactionRepo.GetReactionByUser(ctx, data.ConversationID, data.MessageID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to get message reaction: %w", err)
	}
	if dataMessageReaction == nil {
		return fmt.Errorf("message reaction not found for user %s on message %s in conversation %s", data.UserID, data.MessageID, data.ConversationID)
	}
	dataMessageReaction.ReactionCode = data.ReactionCode
	err = c.messageReactionRepo.UpdateOrInsertReaction(ctx, dataMessageReaction)
	if err != nil {
		return fmt.Errorf("failed to update message reaction: %w", err)
	}
	return nil
}
func (c *ConsumerStatsMessage) handleUpdatedStatsMessage(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communicationEvent.MessageStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataMessageReaction, err := c.messageReactionRepo.GetReactionByUser(ctx, data.ConversationID, data.MessageID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to get message reaction: %w", err)
	}
	if dataMessageReaction == nil {
		return fmt.Errorf("message reaction not found for user %s on message %s in conversation %s", data.UserID, data.MessageID, data.ConversationID)
	}
	dataMessageReaction.ReactionCode = data.ReactionCode
	err = c.messageReactionRepo.UpdateOrInsertReaction(ctx, dataMessageReaction)
	if err != nil {
		return fmt.Errorf("failed to update message reaction: %w", err)
	}
	return nil
}
func (c *ConsumerStatsMessage) handleDeletedStatsMessage(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communicationEvent.MessageStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataMessageReaction, err := c.messageReactionRepo.GetReactionByUser(ctx, data.ConversationID, data.MessageID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to get message reaction: %w", err)
	}
	if dataMessageReaction == nil {
		return fmt.Errorf("message reaction not found for user %s on message %s in conversation %s", data.UserID, data.MessageID, data.ConversationID)
	}
	err = c.messageReactionRepo.DeleteReaction(ctx, data.ConversationID, data.MessageID, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete message reaction: %w", err)
	}
	return nil
}
func (c *ConsumerStatsMessage) ConsumerFailedStatsMessage(ctx context.Context) error {
	return nil
}
