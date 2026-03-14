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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerGroupMember struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewConsumerGroupMember(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository) *ConsumerGroupMember {
	return &ConsumerGroupMember{
		events:          events,
		pool:            pool,
		redisRepo:       redisRepo,
		groupMemberRepo: groupMemberRepo,
	}
}

func (c *ConsumerGroupMember) ConsumerGroupMember(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGroupMember.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events)) // Channel để thu thập lỗi từ goroutine
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
					processErr = c.handleCreatedGroupMember(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedGroupMember(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedGroupMember(ctx, event)

				default:
					processErr = errors.New("unsupported event type: " + event.Type)
					return
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Mở khóa ngay nếu không thể xử lý
				wg.Done()                         // Giảm counter ngay vì không vào được goroutine
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Mở khóa ngay nếu xử lý thất bại
			} else {
				errchan <- nil // Không có lỗi
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
		}
		wg.Wait()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err) // Gom tất cả lỗi lại
			}
		}
		return finalErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerGroupMember) handleCreatedGroupMember(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.CreateGroupMemberPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	memberEntity, err := mapper.ToGroupMemberEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	err = c.groupMemberRepo.CreateGroupMember(ctx, memberEntity)
	if err != nil {
		return fmt.Errorf("failed to create group member in repository for event %s: %w", event.ID, err)
	}
	payload := &communityEvent.GroupStatsPayload{
		GroupID:            data.GroupID,
		MemberCount:        0,
		PostCount:          0,
		PendingMemberCount: 1,
		PendingPostCount:   0,
		EventType:          constants.Created,
	}
	err = c.events.Publish(ctx, constants.TopicGroupStats.String(), data.GroupID, constants.Created.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish group stats event after creating group member for event %s: %w", event.ID, err)
	}
	return nil
}

func (c *ConsumerGroupMember) handleUpdatedGroupMember(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.UpdateGroupMemberPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataaction, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserActionID, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to fetch action user group member from repository for event %s: %w", event.ID, err)
	}
	if dataaction == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("action user group member with user_id %s and group_id %s not found for event %s", data.UserActionID, data.GroupID, event.ID))
	}
	if dataaction.Role != sharedEnums.RoleTypeAdmin {
		return kafka.NewNonRetryableError(fmt.Errorf("user with ID %s does not have permission to update group member in group %s for event %s", data.UserActionID, data.GroupID, event.ID))
	}
	datamember, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserID, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get existing group member from repository for event %s: %w", event.ID, err)
	}
	if datamember == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("group member not found for user_id %s and group_id %s for event %s", data.UserID, data.GroupID, event.ID))
	}
	if datamember.Role == dataaction.Role {
		return kafka.NewNonRetryableError(fmt.Errorf("no role change for group member with user_id %s in group %s for event %s", data.UserID, data.GroupID, event.ID))
	}
	memberEntity := mapper.ApplyUpdateToGroupMember(datamember, data)
	if memberEntity.Status == sharedEnums.ProcessingFailed {
		payloadDelete := &communityEvent.DeleteGroupMemberPayload{
			UserID:    data.UserID,
			GroupID:   data.GroupID,
			DeleteALL: false,
		}
		err = c.events.Publish(ctx, constants.TopicGroupMember.String(), event.ID, constants.Deleted.String(), payloadDelete)
		if err != nil {
			return fmt.Errorf("failed to publish group member deletion event for processing failure for event %s: %w", event.ID, err)
		}
		payloadStats := &communityEvent.GroupStatsPayload{
			GroupID:            data.GroupID,
			MemberCount:        0,
			PostCount:          0,
			PendingMemberCount: -1,
			PendingPostCount:   0,
			EventType:          constants.Updated,
		}
		err = c.events.Publish(ctx, constants.TopicGroupStats.String(), data.GroupID, constants.Updated.String(), payloadStats)
		if err != nil {
			return fmt.Errorf("failed to publish group stats event after processing failure for event %s: %w", event.ID, err)
		}
		return nil // Không cần lỗi nữa vì đã xử lý bằng cách tạo sự kiện xóa
	}
	if memberEntity.Status == sharedEnums.ProcessingActive {
		payloadStats := &communityEvent.GroupStatsPayload{
			GroupID:            data.GroupID,
			MemberCount:        1,
			PostCount:          0,
			PendingMemberCount: 0,
			PendingPostCount:   0,
			EventType:          constants.Updated,
		}
		err = c.events.Publish(ctx, constants.TopicGroupStats.String(), data.GroupID, constants.Updated.String(), payloadStats)
		if err != nil {
			return fmt.Errorf("failed to publish group stats event after member activation for event %s: %w", event.ID, err)
		}
	}
	err = c.groupMemberRepo.UpdateGroupMember(ctx, memberEntity)
	if err != nil {
		return fmt.Errorf("failed to update group member in repository for event %s: %w", event.ID, err)
	}
	return nil
}

func (c *ConsumerGroupMember) handleDeletedGroupMember(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.DeleteGroupMemberPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	if data.DeleteALL == true {

		err = c.groupMemberRepo.DeleteGroupMember(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to delete all group members by group_id in repository for event %s: %w", event.ID, err)
		}
	} else {
		dataaction, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserActionID, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to fetch action user group member from repository for event %s: %w", event.ID, err)
		}
		if dataaction == nil {
			return kafka.NewNonRetryableError(fmt.Errorf("action user group member with user_id %s and group_id %s not found for event %s", data.UserActionID, data.GroupID, event.ID))
		}
		if dataaction.Role != sharedEnums.RoleTypeAdmin {
			return kafka.NewNonRetryableError(fmt.Errorf("user with ID %s does not have permission to update group member in group %s for event %s", data.UserActionID, data.GroupID, event.ID))
		}
		err = c.groupMemberRepo.DeleteGroupMemberByUserIDAndGroupID(ctx, data.UserID, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to delete group member in repository for event %s: %w", event.ID, err)
		}
	}
	payload := &communityEvent.GroupStatsPayload{
		GroupID:            data.GroupID,
		MemberCount:        1,
		PostCount:          0,
		PendingMemberCount: 0,
		PendingPostCount:   0,
		EventType:          constants.Deleted,
	}
	err = c.events.Publish(ctx, constants.TopicGroupStats.String(), data.GroupID, constants.Deleted.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish group stats event after deleting group member for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerGroupMember) ConsumerFailedGroupMember(ctx context.Context) error {
	return nil
}
