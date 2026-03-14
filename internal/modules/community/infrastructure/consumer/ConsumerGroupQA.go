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

type ConsumerGroupQA struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	groupQARepo     IRepositoryMongodb.IGroupjoinQuestionsRepository
}

func NewConsumerGroupQA(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupQARepo IRepositoryMongodb.IGroupjoinQuestionsRepository) *ConsumerGroupQA {
	return &ConsumerGroupQA{
		events:      events,
		pool:        pool,
		redisRepo:   redisRepo,
		groupQARepo: groupQARepo,
	}
}

func (c *ConsumerGroupQA) ConsumerGroupQA(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGroupQA.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Đảm bảo unlock nếu không thể chạy worker
				wg.Done()                         // Giảm wg nếu không thể chạy worker
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID) // Đảm bảo unlock nếu có lỗi trong xử lý
			} else {
				errchan <- nil                           // Xử lý thành công, gửi nil vào channel
				c.redisRepo.MarkCompleted(ctx, event.ID) // Đánh dấu hoàn thành để các worker khác biết mà không xử lý lại
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
func (c *ConsumerGroupQA) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.CreateGroupJoinQuestionPayload](event.Payload)
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
	entity, err := mapper.ToGroupJoinQuestionEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	if entity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("mapped entity is nil for event %s", event.ID))
	}
	err = c.groupQARepo.CreateGroupJoinQuestion(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create GroupJoinQuestion in repository for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerGroupQA) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.UpdateGroupJoinQuestionPayload](event.Payload)
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
	// Lấy entity hiện tại từ DB để áp dụng update
	existingEntity, err := c.groupQARepo.GetGroupJoinQuestionByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing entity for event %s: %w", event.ID, err)
	}
	if existingEntity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("no existing entity found with ID %s for event %s", data.ID, event.ID))
	}
	mapper.ApplyUpdateToEntity(existingEntity, data)
	err = c.groupQARepo.UpdateGroupJoinQuestion(ctx, existingEntity)
	if err != nil {
		return fmt.Errorf("failed to update GroupJoinQuestion in repository for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerGroupQA) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.DeleteGroupJoinQuestionPayload](event.Payload)
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
	if data.DeleteAll == true {
		_, err := c.groupQARepo.DeleteGroupJoinQuestionsByGroupID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to delete GroupJoinQuestions by group ID %s for event %s: %w", data.GroupID, event.ID, err)
		}
	} else {
		err = c.groupQARepo.DeleteGroupJoinQuestionByIDAndGroupID(ctx, data.ID, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to delete GroupJoinQuestion with ID %s and group ID %s for event %s: %w", data.ID, data.GroupID, event.ID, err)
		}
	}
	return nil
}
func (c *ConsumerGroupQA) ConsumerFailedGroupQA(ctx context.Context) error {
	return nil
}
