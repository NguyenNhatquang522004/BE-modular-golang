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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
)

type ConsumerUserSetting struct {
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	UserSetting IRepositoryMongodb.IUserNotificationSettingsRepository
	redisRepo   IRepositoryShare.IRedis
}

func NewConsumerUserSetting(events events.EventBus, pool IRepositoryShare.IWorkerPool, userSetting IRepositoryMongodb.IUserNotificationSettingsRepository, redisRepo IRepositoryShare.IRedis) *ConsumerUserSetting {
	return &ConsumerUserSetting{
		events:      events,
		pool:        pool,
		UserSetting: userSetting,
		redisRepo:   redisRepo,
	}
}

func (c *ConsumerUserSetting) ConsumerUserNotificationSettings(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicUserNotificationSettings.String(), 100, 5*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			redisKeyPrefix := "consumer_user_setting_lock:" + event.ID
			// 1. Thử khóa event này trong Redis để đảm bảo chỉ 1 worker xử lý 1 eventID nhất định (Distributed Lock)
			status, acquired, err := c.redisRepo.Lock(ctx, redisKeyPrefix)
			if err != nil {
				if status != "" && status != constants.StatusProcessing {
					return err
				}
				return err
			}
			if !acquired {
				if status == constants.StatusProcessing {
					errchan <- errors.New("event is being processed by another worker")
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
					processErr = c.handleCreatedUserNotificationSettings(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedUserNotificationSettings(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedUserNotificationSettings(ctx, ev)
				default:
					processErr = fmt.Errorf("unknown event type: %s", ev.Type)
				}
			})
			if processErr != nil {
				c.redisRepo.Unlock(ctx, redisKeyPrefix) // Thất bại -> Mở khoá để lần sau làm lại
				errchan <- fmt.Errorf("event %s failed: %w", ev.ID, processErr)
			} else {
				// CỰC KỲ QUAN TRỌNG: Thành công -> Đánh dấu Vĩnh viễn (Hoặc 24h)
				c.redisRepo.MarkCompleted(ctx, redisKeyPrefix)
				errchan <- nil
			}
			if err != nil {
				errchan <- err
				c.redisRepo.Unlock(ctx, redisKeyPrefix) // Mở khóa ngay nếu có lỗi khi chạy goroutine
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
		return err
	}
	return nil
}
func (c *ConsumerUserSetting) handleCreatedUserNotificationSettings(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[notificationEvent.NotificationChangePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(err)
	}
	if data == nil {
		return kafka.NewNonRetryableError(err)
	}
	entity := mapper.ToEntityUserNotificationSettingPayload(data)
	if entity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s", event.ID))
	}
	err = c.UserSetting.CreateUserNotificationSettings(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create user notification setting for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerUserSetting) handleUpdatedUserNotificationSettings(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[notificationEvent.NotificationChangePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(err)
	}
	if data == nil {
		return kafka.NewNonRetryableError(err)
	}
	existing, err := c.UserSetting.GetUserNotificationSettingsByUserID(ctx, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to get existing user notification setting for event %s: %w", event.ID, err)
	}
	if existing == nil {
		return fmt.Errorf("no existing user notification setting found for user_id %s in event %s", data.UserID, event.ID)
	}
	mapper.UpdateToEntityUserNotificationSettingPayload(data, existing)
	err = c.UserSetting.UpdateUserNotificationSettings(ctx, existing)
	if err != nil {
		return fmt.Errorf("failed to update user notification setting for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerUserSetting) handleDeletedUserNotificationSettings(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[notificationEvent.NotificationChangePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(err)
	}
	if data == nil {
		return kafka.NewNonRetryableError(err)
	}
	err = c.UserSetting.DeleteUserNotificationSettings(ctx, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete user notification setting for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerUserSetting) ConsumerFailedUserNotificationSettings(ctx context.Context) error {
	return nil
}
