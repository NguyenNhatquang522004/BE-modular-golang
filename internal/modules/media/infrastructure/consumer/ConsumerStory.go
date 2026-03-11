package consumer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerStory struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	storyRepo     IRepositoryMongodb.IStoryRepository
	storyViewRepo IRepositoryCassandra.IStoryViewRepository
}

func NewConsumerStory(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, storyRepo IRepositoryMongodb.IStoryRepository, storyViewRepo IRepositoryCassandra.IStoryViewRepository) *ConsumerStory {
	return &ConsumerStory{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		storyRepo:     storyRepo,
		storyViewRepo: storyViewRepo,
	}
}
func (c *ConsumerStory) ConsumerStory(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicStory.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
		wg.Done()
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
func (c *ConsumerStory) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {

	return nil
}
func (c *ConsumerStory) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {

	return nil
}
func (c *ConsumerStory) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {

	return nil
}
func (c *ConsumerStory) ConsumerFailedStory(ctx context.Context) error {

	return nil
}
