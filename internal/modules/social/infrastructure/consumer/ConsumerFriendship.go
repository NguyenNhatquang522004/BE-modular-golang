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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type ConsumerFriendship struct {
	redisRepo  IRepositoryShare.IRedis
	events     events.EventBus
	pool       IRepositoryShare.IWorkerPool
	friendRepo IRepositoryPostgres.IFriendshipsRepository
}

func NewConsumerFriendship() *ConsumerFriendship {
	return &ConsumerFriendship{}
}

func (c *ConsumerFriendship) ConsumerFriendUser(ctx context.Context) error {
	// Implement the logic for consuming friend user events
	err := c.events.SubscribeBatch(ctx, constants.TopicFriendship.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			// 1. Thử khóa event này trong Redis để đảm bảo chỉ 1 worker xử lý 1 eventID nhất định (Distributed Lock)
			status, acquired, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- err
				continue
			}
			if !acquired {
				if status == constants.StatusProcessing {
					// Log và bỏ qua event này vì đang có worker khác xử lý
					log.Printf("Event %s is currently being processed by another worker. Skipping.\n", event.ID)
					errchan <- fmt.Errorf("Event %s is currently being processed by another worker. Skipping.", event.ID) // Không coi đây là lỗi, chỉ là tình trạng bình thường khi có nhiều worker
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
					processErr = c.handleCreateFriendshipEvent(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdateFriendshipEvent(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeleteFriendshipEvent(ctx, ev)
				default:
					log.Printf("Unknown event type %s for event %s. Skipping.\n", ev.Type, ev.ID)
					return
				}
			})
			if processErr != nil {
				errchan <- fmt.Errorf("failed to process event %s: %w", event.ID, processErr)
				c.redisRepo.Unlock(ctx, ev.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, ev.ID)
			}
			if err != nil {
				wg.Done()
				errchan <- fmt.Errorf("failed to run event %s in worker pool: %w", event.ID, err)
				c.redisRepo.Unlock(ctx, ev.ID) // Mở khóa ngay nếu có lỗi khi chạy goroutine
			}
		}
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
func (c *ConsumerFriendship) handleCreateFriendshipEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling create friendship events
	data, err := utils.ParsePayload[socialEvent.FriendUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity := &entity.Friendships{
		ID:           uuid.New(),
		Requester_ID: uuid.MustParse(data.RequesterID),
		Recipient_ID: uuid.MustParse(data.RecipientID),
		Status:       data.Status,
		Created_At:   time.Now(),
		Updated_At:   time.Now(),
	}
	err = c.friendRepo.CreateFriendship(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create friendship: %w", err)
	}
	payloadfollower := socialEvent.FollowerUserPayload{
		ID:             uuid.New().String(),
		FollowerUserID: data.RequesterID,
		FollowedUserID: data.RecipientID,
		IsMuted:        false,
		EventType:      constants.Created,
	}
	err = c.events.Publish(ctx, constants.TopicFollowUser.String(), payloadfollower.ID, constants.Created.String(), payloadfollower)
	if err != nil {
		return fmt.Errorf("failed to publish follower event: %w", err)
	}
	return nil
}
func (c *ConsumerFriendship) handleUpdateFriendshipEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling update friendship events
	data, err := utils.ParsePayload[socialEvent.FriendUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity, err := c.friendRepo.GetFriendshipByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to get friendship: %w", err)
	}
	if entity == nil {
		return fmt.Errorf("friendship not found")
	}
	switch data.Status {
	case sharedEnums.StatusFriendship_Accepted:
		entity.Status = sharedEnums.StatusFriendship_Accepted
		entity.Updated_At = time.Now()
		err = c.friendRepo.UpdateFriendship(ctx, entity)
		if err != nil {
			return fmt.Errorf("failed to update friendship: %w", err)
		}
		payloadfollower := socialEvent.FollowerUserPayload{
			ID:             uuid.New().String(),
			FollowerUserID: entity.Recipient_ID.String(),
			FollowedUserID: entity.Requester_ID.String(),
			IsMuted:        false,
			CreatedAt:      data.Created_At,
			UpdatedAt:      data.Updated_At,
			EventType:      constants.Created,
		}
		err = c.events.Publish(ctx, constants.TopicFollowUser.String(), payloadfollower.ID, constants.Created.String(), payloadfollower)
		if err != nil {
			return fmt.Errorf("failed to publish follower event: %w", err)
		}
	case sharedEnums.StatusFriendship_Blocked:
		entity.Status = sharedEnums.StatusFriendship_Blocked
		entity.Updated_At = time.Now()
		err := c.friendRepo.DeleteFriendship(ctx, data.ID)
		if err != nil {
			return fmt.Errorf("failed to delete friendship: %w", err)
		}
		payloadfollower := socialEvent.FollowerUserPayload{
			FollowerUserID: entity.Recipient_ID.String(),
			FollowedUserID: entity.Requester_ID.String(),
			UpdatedAt:      data.Updated_At,
			EventType:      constants.Deleted,
		}
		err = c.events.Publish(ctx, constants.TopicFollowUser.String(), payloadfollower.ID, constants.Deleted.String(), payloadfollower)
		if err != nil {
			return fmt.Errorf("failed to publish follower event: %w", err)
		}
	default:
		return kafka.NewNonRetryableError(fmt.Errorf("invalid friendship status: %s", data.Status))
	}
	return nil
}
func (c *ConsumerFriendship) handleDeleteFriendshipEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling delete friendship events
	data, err := utils.ParsePayload[socialEvent.FriendUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.Status != sharedEnums.StatusFriendship_Declined {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid event for delete friendship: status must be declined"))
	}
	err = c.friendRepo.DeleteFriendshipByUserIDs(ctx, data.RequesterID, data.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to delete friendship: %w", err)
	}
	return nil
}
func (c *ConsumerFriendship) ConsumerFailedFriendUser(ctx context.Context) error {
	// Implement the logic for consuming failed friend user events
	return nil
}
