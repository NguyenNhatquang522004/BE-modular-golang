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
	return nil
}

func (c *ConsumerConversation) handleUpdatedConversation(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated conversation messages
	return nil
}

func (c *ConsumerConversation) handleDeletedConversation(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted conversation messages
	return nil
}
func (c *ConsumerConversation) ConsumerFailedConversation(ctx context.Context) error {
	// Implement the logic for handling failed conversation messages
	return nil
}
