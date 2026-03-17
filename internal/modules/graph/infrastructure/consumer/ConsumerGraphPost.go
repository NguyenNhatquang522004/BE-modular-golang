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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/entity"
)

type ConsumerGraphPost struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	graphRepo neo4j.IGraphRepository
}

func NewConsumerGraphPost(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, graphRepo neo4j.IGraphRepository) *ConsumerGraphPost {
	return &ConsumerGraphPost{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		graphRepo: graphRepo,
	}
}

func (c *ConsumerGraphPost) ConsumerGraphPost(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGraphPost.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
func (c *ConsumerGraphPost) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý sự kiện Created cho GraphPost
	data, err := utils.ParsePayload[graphEvent.PostNodePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.graphRepo.UpsertPostNode(ctx, &entity.PostNode{
		PostID:    data.PostID,
		CreatedAt: data.CreatedAt,
		TTL:       data.TTL,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert post node: %w", err)
	}
	for _, topic := range data.Topic {
		topicID, _, err := c.graphRepo.UpsertTopicAndConnectNeighbors(ctx, topic.Name)
		if err != nil {
			return fmt.Errorf("failed to upsert topic and connect neighbors: %w", err)
		}
		err = c.graphRepo.LinkPostToTopicByID(ctx, data.PostID, topicID, topic.ConfidenceScore)
		if err != nil {
			return fmt.Errorf("failed to link post to topic: %w", err)
		}
	}
	err = c.graphRepo.LinkAuthorToPost(ctx, data.UserID, data.PostID, data.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to link author to post: %w", err)
	}

	switch data.TargetType {
	case sharedEnums.ContextTypeGroup:
		err = c.graphRepo.LinkPostToGroup(ctx, data.PostID, data.TargetType.String(), data.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to link post to group: %w", err)
		}
	case sharedEnums.ContextTypePage:
		err = c.graphRepo.LinkPageToPost(ctx, data.TargetType.String(), data.PostID, data.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to link page to post: %w", err)
		}
	case sharedEnums.ContextTypeUserWall:
		err = c.graphRepo.IncrementTopicInterest(ctx, data.UserID, 50) // Cập nhật sở thích của người dùng dựa trên tương tác gần đây
		if err != nil {
			return fmt.Errorf("failed to increment topic interest for event %s: %w", event.ID, err)
		}
	default:
		return fmt.Errorf("unknown target type: %s", data.TargetType.String())
	}
	return nil
}

func (c *ConsumerGraphPost) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý sự kiện Created cho GraphPost
	data, err := utils.ParsePayload[graphEvent.PostNodePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.graphRepo.UpsertPostNode(ctx, &entity.PostNode{
		PostID:    data.PostID,
		CreatedAt: data.CreatedAt,
		TTL:       data.TTL,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert post node: %w", err)
	}
	for _, topic := range data.Topic {
		topicID, _, err := c.graphRepo.UpsertTopicAndConnectNeighbors(ctx, topic.Name)
		if err != nil {
			return fmt.Errorf("failed to upsert topic and connect neighbors: %w", err)
		}
		err = c.graphRepo.LinkPostToTopicByID(ctx, data.PostID, topicID, topic.ConfidenceScore)
		if err != nil {
			return fmt.Errorf("failed to link post to topic: %w", err)
		}
	}
	err = c.graphRepo.LinkAuthorToPost(ctx, data.UserID, data.PostID, data.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to link author to post: %w", err)
	}
	switch data.TargetType {
	case sharedEnums.ContextTypeGroup:
		err = c.graphRepo.LinkPostToGroup(ctx, data.PostID, data.TargetType.String(), data.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to link post to group: %w", err)
		}
	case sharedEnums.ContextTypePage:
		err = c.graphRepo.LinkPageToPost(ctx, data.TargetType.String(), data.PostID, data.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to link page to post: %w", err)
		}
	case sharedEnums.ContextTypeUserWall:
		err = c.graphRepo.IncrementTopicInterest(ctx, data.UserID, 50) // Cập nhật sở thích của người dùng dựa trên tương tác gần đây
		if err != nil {
			return fmt.Errorf("failed to increment topic interest for event %s: %w", event.ID, err)
		}
	default:
		return fmt.Errorf("unknown target type: %s", data.TargetType.String())
	}
	return nil
}

func (c *ConsumerGraphPost) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý sự kiện Deleted cho GraphPost
	data, err := utils.ParsePayload[graphEvent.DeletePostNodePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.graphRepo.DeletePostNode(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to delete post node: %w", err)
	}
	return nil
}
func (c *ConsumerGraphPost) ConsumerFailedGraphPost(ctx context.Context) error {

	return nil
}
