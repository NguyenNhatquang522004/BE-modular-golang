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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/graphEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/IRepository/neo4j"
)

type ConsumerGraphInteractionsRecentConsumer struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	graphRepo neo4j.IGraphRepository
}

func NewConsumerGraphInteractionsRecentConsumer(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, graphRepo neo4j.IGraphRepository) *ConsumerGraphInteractionsRecentConsumer {
	return &ConsumerGraphInteractionsRecentConsumer{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		graphRepo: graphRepo,
	}
}
func (c *ConsumerGraphInteractionsRecentConsumer) ConsumerGraphInteractionsRecent(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGraphInteractionsRecent.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
func (c *ConsumerGraphInteractionsRecentConsumer) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[graphEvent.InteractionRecentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	like, comment, share, message, view := 0, 0, 0, 0, 0
	switch data.Type {
	case sharedEnums.ReactionTargetView:
		data.Weight = 1.0
		view = 1
	case sharedEnums.ReactionTargetLike:
		data.Weight = 3.0
		like = 1
	case sharedEnums.ReactionTargetShare:
		data.Weight = 10.0
		share = 1
	case sharedEnums.ReactionTargetComment:
		data.Weight = 5.0
		comment = 1
	default:
		return kafka.NewNonRetryableError(fmt.Errorf("unknown interaction type %s for event %s", data.Type.String(), event.ID))
	}
	err = c.graphRepo.RecordRecentInteraction(ctx, data.UserID, data.PostID, data.Type.String(), data.Weight, data.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to record recent interaction for event %s: %w", event.ID, err)
	}
	err = c.graphRepo.IncrementTopicInterest(ctx, data.UserID, 50) // Cập nhật sở thích của người dùng dựa trên tương tác gần đây
	if err != nil {
		return fmt.Errorf("failed to increment topic interest for event %s: %w", event.ID, err)
	}
	_, err = c.graphRepo.IncrementInteractionByPost(ctx, data.UserID, data.PostID, like, comment, share, message, view, true)
	if err != nil {
		return fmt.Errorf("failed to increment interaction by post for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerGraphInteractionsRecentConsumer) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[graphEvent.InteractionRecentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	like, comment, share, message, view := 0, 0, 0, 0, 0
	switch data.Type {
	case sharedEnums.ReactionTargetView:
		data.Weight = 1.0

	case sharedEnums.ReactionTargetLike:
		data.Weight = 3.0

	case sharedEnums.ReactionTargetShare:
		data.Weight = 10.0

	case sharedEnums.ReactionTargetComment:
		data.Weight = 5.0

	default:
		return kafka.NewNonRetryableError(fmt.Errorf("unknown interaction type %s for event %s", data.Type.String(), event.ID))
	}

	err = c.graphRepo.RecordRecentInteraction(ctx, data.UserID, data.PostID, data.Type.String(), data.Weight, data.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to record recent interaction for event %s: %w", event.ID, err)
	}
	err = c.graphRepo.IncrementTopicInterest(ctx, data.UserID, 50) // Cập nhật sở thích của người dùng dựa trên tương tác gần đây
	if err != nil {
		return fmt.Errorf("failed to increment topic interest for event %s: %w", event.ID, err)
	}
	_, err = c.graphRepo.IncrementInteractionByPost(ctx, data.UserID, data.PostID, like, comment, share, message, view, true)
	if err != nil {
		return fmt.Errorf("failed to increment interaction by post for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerGraphInteractionsRecentConsumer) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[graphEvent.InteractionRecentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	like, comment, share, message, view := 0, 0, 0, 0, 0
	switch data.Type {
	case sharedEnums.ReactionTargetView:
		data.Weight = -1.0
		view = -1
	case sharedEnums.ReactionTargetLike:
		data.Weight = -3.0
		like = -1
	case sharedEnums.ReactionTargetShare:
		data.Weight = -10.0
		share = -1
	case sharedEnums.ReactionTargetComment:
		data.Weight = -5.0
		comment = -1
	default:
		return kafka.NewNonRetryableError(fmt.Errorf("unknown interaction type %s for event %s", data.Type.String(), event.ID))
	}
	err = c.graphRepo.RecordRecentInteraction(ctx, data.UserID, data.PostID, data.Type.String(), data.Weight, data.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to record recent interaction for event %s: %w", event.ID, err)
	}
	err = c.graphRepo.IncrementTopicInterest(ctx, data.UserID, 50) // Cập nhật sở thích của người dùng dựa trên tương tác gần đây
	if err != nil {
		return fmt.Errorf("failed to increment topic interest for event %s: %w", event.ID, err)
	}
	_, err = c.graphRepo.IncrementInteractionByPost(ctx, data.UserID, data.PostID, like, comment, share, message, view, false)
	if err != nil {
		return fmt.Errorf("failed to increment interaction by post for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerGraphInteractionsRecentConsumer) ConsumerFailedGraphInteractionsRecent(ctx context.Context) error {
	return nil
}
