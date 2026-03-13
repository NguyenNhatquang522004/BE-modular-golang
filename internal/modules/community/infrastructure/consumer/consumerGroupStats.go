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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerGroupStats struct {
	events    events.EventBus
	redisRepo IRepositoryShare.IRedis
	groupRepo IRepositoryMongodb.IGroupRepository
	pool      IRepositoryShare.IWorkerPool
}

func NewConsumerGroupStats(events events.EventBus, redisRepo IRepositoryShare.IRedis, groupRepo IRepositoryMongodb.IGroupRepository, pool IRepositoryShare.IWorkerPool) *ConsumerGroupStats {
	return &ConsumerGroupStats{
		events:    events,
		redisRepo: redisRepo,
		groupRepo: groupRepo,
		pool:      pool,
	}
}

func (c *ConsumerGroupStats) ConsumerGroupStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGroupStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedEvent(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedEvent(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedEvent(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
					return
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
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
func (c *ConsumerGroupStats) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.GroupStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datagroup, err := c.groupRepo.GetGroupByID(ctx, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group by ID %s for event %s: %w", data.GroupID, event.ID, err)
	}
	if datagroup == nil {
		return fmt.Errorf("group not found with ID %s for event %s", data.GroupID, event.ID)
	}
	// Cập nhật lại số liệu thống kê của nhóm dựa trên dữ liệu mới
	datagroup.Stats.MemberCount = data.MemberCount + datagroup.Stats.MemberCount
	datagroup.Stats.PostCount = data.PostCount + datagroup.Stats.PostCount
	datagroup.Stats.PendingMemberCount = data.PendingMemberCount + datagroup.Stats.PendingMemberCount
	datagroup.Stats.PendingPostCount = data.PendingPostCount + datagroup.Stats.PendingPostCount
	datagroup.Stats.ReportedPostCount = data.ReportedPostCount + datagroup.Stats.ReportedPostCount
	err = c.groupRepo.UpdateGroup(ctx, datagroup)
	if err != nil {
		return fmt.Errorf("failed to update group stats for group ID %s for event %s: %w", data.GroupID, event.ID, err)
	}
	return nil
}
func (c *ConsumerGroupStats) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.GroupStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datagroup, err := c.groupRepo.GetGroupByID(ctx, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group by ID %s for event %s: %w", data.GroupID, event.ID, err)
	}
	if datagroup == nil {
		return fmt.Errorf("group not found with ID %s for event %s", data.GroupID, event.ID)
	}
	// Cập nhật lại số liệu thống kê của nhóm dựa trên dữ liệu mới
	datagroup.Stats.MemberCount = data.MemberCount + datagroup.Stats.MemberCount
	datagroup.Stats.PostCount = data.PostCount + datagroup.Stats.PostCount
	datagroup.Stats.PendingMemberCount = data.PendingMemberCount + datagroup.Stats.PendingMemberCount
	datagroup.Stats.PendingPostCount = data.PendingPostCount + datagroup.Stats.PendingPostCount
	datagroup.Stats.ReportedPostCount = data.ReportedPostCount + datagroup.Stats.ReportedPostCount
	err = c.groupRepo.UpdateGroup(ctx, datagroup)
	if err != nil {
		return fmt.Errorf("failed to update group stats for group ID %s for event %s: %w", data.GroupID, event.ID, err)
	}
	return nil
}
func (c *ConsumerGroupStats) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.GroupStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	datagroup, err := c.groupRepo.GetGroupByID(ctx, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group by ID %s for event %s: %w", data.GroupID, event.ID, err)
	}
	if datagroup == nil {
		return fmt.Errorf("group not found with ID %s for event %s", data.GroupID, event.ID)
	}
	// Cập nhật lại số liệu thống kê của nhóm dựa trên dữ liệu mới
	datagroup.Stats.MemberCount = datagroup.Stats.MemberCount - data.MemberCount
	datagroup.Stats.PostCount = datagroup.Stats.PostCount - data.PostCount
	datagroup.Stats.PendingMemberCount = datagroup.Stats.PendingMemberCount - data.PendingMemberCount
	datagroup.Stats.PendingPostCount = datagroup.Stats.PendingPostCount - data.PendingPostCount
	datagroup.Stats.ReportedPostCount = datagroup.Stats.ReportedPostCount - data.ReportedPostCount
	err = c.groupRepo.UpdateGroup(ctx, datagroup)
	if err != nil {
		return fmt.Errorf("failed to update group stats for group ID %s for event %s: %w", data.GroupID, event.ID, err)
	}
	return nil
}
func (c *ConsumerGroupStats) ConsumerFailedGroupStats(ctx context.Context) error {
	return nil
}
