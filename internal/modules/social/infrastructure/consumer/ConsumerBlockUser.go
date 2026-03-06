package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
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

type ConsumerBlockUser struct {
	redisRepo IRepositoryShare.IRedis
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	blockRepo IRepositoryPostgres.IBlockRepository
}

func NewConsumerBlockUser(redisRepo IRepositoryShare.IRedis, events events.EventBus, pool IRepositoryShare.IWorkerPool, blockRepo IRepositoryPostgres.IBlockRepository) *ConsumerBlockUser {
	return &ConsumerBlockUser{
		redisRepo: redisRepo,
		events:    events,
		pool:      pool,
		blockRepo: blockRepo,
	}
}

func (c *ConsumerBlockUser) ConsumerBlockUser(ctx context.Context) error {
	// Implement the logic for consuming block user events
	err := c.events.SubscribeBatch(ctx, constants.TopicBlockUser.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			wg.Add(1)
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- fmt.Errorf("failed to acquire lock for event %s: %w", event.ID, err)
				continue
			}
			if !can {
				if status == constants.StatusProcessing {
					errchan <- fmt.Errorf("event %s is already being processed by another consumer", event.ID)
				} else {
					log.Printf("Event %s has already been processed with status %s. Skipping.\n", event.ID, status)
				}
				continue
			}
			var processErr error
			ev := event
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreatedBlockUser(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedBlockUser(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedBlockUser(ctx, ev)
				default:
					errchan <- fmt.Errorf("unsupported event type %s for event ID %s. Marking as failed", ev.Type, ev.ID)
				}
			})
			if processErr != nil {
				errchan <- fmt.Errorf("failed to process event %s: %w", ev.ID, processErr)
				c.redisRepo.Unlock(ctx, ev.ID)
			} else {
				c.redisRepo.MarkCompleted(ctx, ev.ID)
				errchan <- nil
			}
			if err != nil {
				errchan <- fmt.Errorf("failed to run worker for event %s: %w", ev.ID, err)
				c.redisRepo.Unlock(ctx, ev.ID)
				wg.Done()
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
		return fmt.Errorf("failed to subscribe to block user events: %w", err)
	}
	return nil
}
func (c *ConsumerBlockUser) handleCreatedBlockUser(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created block user events
	data, err := utils.ParsePayload[socialEvent.BlockUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	entity := &entity.UserBlock{
		ID:             uuid.New(),
		Blocker_UserID: uuid.MustParse(data.BlockerUserID),
		Blocked_UserID: uuid.MustParse(data.BlockedUserID),
		Type_Block:     data.TypeBlock,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = c.blockRepo.CreateBlock(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create block relationship for event %s: %w", event.ID, err)
	}
	if slices.Contains(data.TypeBlock, sharedEnums.Type_Block_Full) || slices.Contains(data.TypeBlock, sharedEnums.Type_Block_Profile) {
		payloadFollower := socialEvent.FollowerUserPayload{
			FollowerUserID: data.BlockerUserID,
			FollowedUserID: data.BlockedUserID,
			EventType:      constants.Deleted,
		}
		err = c.events.Publish(ctx, constants.TopicFollowUser.String(), data.BlockedUserID, constants.Deleted.String(), payloadFollower)
		if err != nil {
			log.Printf("Failed to publish follow user deletion event for block user event %s: %v\n", event.ID, err)
		}
		payloadFriend := socialEvent.FriendUserPayload{
			RequesterID: data.BlockerUserID,
			RecipientID: data.BlockedUserID,
			EventType:   constants.Deleted,
		}
		err = c.events.Publish(ctx, constants.TopicFriendship.String(), data.BlockedUserID, constants.Deleted.String(), payloadFriend)
		if err != nil {
			log.Printf("Failed to publish friend user deletion event for block user event %s: %v\n", event.ID, err)
		}
	}

	return nil
}
func (c *ConsumerBlockUser) handleUpdatedBlockUser(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated block user events
	data, err := utils.ParsePayload[socialEvent.BlockUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datauser, err := c.blockRepo.GetBlockByUserIDs(ctx, data.BlockerUserID, data.BlockedUserID)
	if err != nil {
		return fmt.Errorf("failed to get existing block relationship for event %s: %w", event.ID, err)
	}
	if datauser == nil {
		return fmt.Errorf("no existing block relationship found for event %s", event.ID)
	}
	datauser.Type_Block = data.TypeBlock
	datauser.UpdatedAt = time.Now()
	err = c.blockRepo.UpdateBlock(ctx, datauser)
	if err != nil {
		return fmt.Errorf("failed to update block relationship for event %s: %w", event.ID, err)
	}
	if slices.Contains(data.TypeBlock, sharedEnums.Type_Block_Full) || slices.Contains(data.TypeBlock, sharedEnums.Type_Block_Profile) {
		payloadFollower := socialEvent.FollowerUserPayload{
			FollowerUserID: data.BlockerUserID,
			FollowedUserID: data.BlockedUserID,
			EventType:      constants.Deleted,
		}
		err = c.events.Publish(ctx, constants.TopicFollowUser.String(), data.BlockedUserID, constants.Deleted.String(), payloadFollower)
		if err != nil {
			log.Printf("Failed to publish follow user deletion event for block user event %s: %v\n", event.ID, err)
		}
		payloadFriend := socialEvent.FriendUserPayload{
			RequesterID: data.BlockerUserID,
			RecipientID: data.BlockedUserID,
			EventType:   constants.Deleted,
		}
		err = c.events.Publish(ctx, constants.TopicFriendship.String(), data.BlockedUserID, constants.Deleted.String(), payloadFriend)
		if err != nil {
			log.Printf("Failed to publish friend user deletion event for block user event %s: %v\n", event.ID, err)
		}
	}
	return nil
}

func (c *ConsumerBlockUser) handleDeletedBlockUser(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted block user events
	data, err := utils.ParsePayload[socialEvent.BlockUserPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	err = c.blockRepo.DeleteBlockByUserIDs(ctx, data.BlockerUserID, data.BlockedUserID)
	if err != nil {
		return fmt.Errorf("failed to delete block relationship for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerBlockUser) ConsumerFailedBlockUser(ctx context.Context) error {
	// Implement the logic for consuming failed block user events
	return nil
}
