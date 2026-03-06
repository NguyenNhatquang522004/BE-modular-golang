package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type ConsumerFollowerUser struct {
	redisRepo  IRepositoryShare.IRedis
	events     events.EventBus
	pool       IRepositoryShare.IWorkerPool
	followRepo IRepositoryPostgres.IFollowersRepository
}

func NewConsumerFollowerUser(redisRepo IRepositoryShare.IRedis, events events.EventBus, pool IRepositoryShare.IWorkerPool) *ConsumerFollowerUser {
	return &ConsumerFollowerUser{
		redisRepo: redisRepo,
		events:    events,
		pool:      pool,
	}
}
func (c *ConsumerFollowerUser) ConsumerFollowUser(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicFollowUser.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			redisprefix := "consumer_follow_user_lock:" + event.ID
			// 1. Thử khóa event này trong Redis để đảm bảo chỉ 1 worker xử lý 1 eventID nhất định (Distributed Lock)
			status, acquired, err := c.redisRepo.Lock(ctx, redisprefix)
			if err != nil {
				errchan <- fmt.Errorf("failed to acquire lock for event %s: %w", event.ID, err)
				continue
			}
			if !acquired {
				if status == constants.StatusProcessing {
					log.Printf("Event %s is currently being processed by another worker. Skipping.\n", event.ID)
					errchan <- fmt.Errorf("event %s is being processed by another worker", event.ID)
				} else {
					log.Printf("Event %s has already been processed with status %s. Skipping.\n", event.ID, status)
				}
				continue
			}
			ev := event
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreateFollowerUserEvent(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdateFollowerUserEvent(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeleteFollowerUserEvent(ctx, ev)
				case constants.Failed.String():
					processErr = c.ConsumerFailedFollowUser(ctx)
				default:
					log.Printf("Unknown event type %s for event %s. Skipping.\n", ev.Type, ev.ID)
				}
			})
			if processErr != nil {
				errchan <- fmt.Errorf("failed to process event %s: %w", ev.ID, processErr)
				c.redisRepo.Unlock(ctx, ev.ID)
				wg.Done()
			} else {
				c.redisRepo.MarkCompleted(ctx, ev.ID)
				errchan <- nil
			}
			if err != nil {
				wg.Done()
				errchan <- fmt.Errorf("failed to run event %s in worker pool: %w", ev.ID, err)
				c.redisRepo.Unlock(ctx, redisprefix) // Mở khóa ngay nếu có lỗi khi chạy goroutine
			}
		}
		wg.Wait()
		close(errchan)
		var batchErr error
		for err := range errchan {
			if err != nil {
				batchErr = errors.Join(batchErr, err)
			}
		}

		return batchErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerFollowerUser) handleCreateFollowerUserEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[socialEvent.FollowerUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	id := uuid.New()
	entity := &entity.Followers{
		ID:              id,
		Follower_UserID: uuid.MustParse(data.FollowerUserID),
		Followed_UserID: uuid.MustParse(data.FollowedUserID),
		IsMuted:         false,
	}
	err = c.followRepo.CreateFollower(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create follower relationship for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerFollowerUser) handleUpdateFollowerUserEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[socialEvent.FollowerUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	
	return nil
}
func (c *ConsumerFollowerUser) handleDeleteFollowerUserEvent(ctx context.Context, event events.IntegrationEvent) error {

	return nil
}
func (c *ConsumerFollowerUser) ConsumerFailedFollowUser(ctx context.Context) error {
	return nil
}
