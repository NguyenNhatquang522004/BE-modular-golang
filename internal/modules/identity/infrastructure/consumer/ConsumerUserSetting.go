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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryPostgres"
)

type ConsumerUserSetting struct {
	events          events.EventBus
	redisRepo       IRepositoryShare.IRedis
	pool            IRepositoryShare.IWorkerPool
	userRepo        IRepositoryPostgres.IUserRepository
	usersettingRepo IRepositoryMongodb.IUserSettingRepository
}

func NewConsumerUserSetting(events events.EventBus, redisRepo IRepositoryShare.IRedis, pool IRepositoryShare.IWorkerPool, userRepo IRepositoryPostgres.IUserRepository, usersettingRepo IRepositoryMongodb.IUserSettingRepository) *ConsumerUserSetting {
	return &ConsumerUserSetting{
		events:          events,
		redisRepo:       redisRepo,
		pool:            pool,
		userRepo:        userRepo,
		usersettingRepo: usersettingRepo,
	}
}

func (c *ConsumerUserSetting) ConsumerUserSettingEvent(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicUserSettings.String(), 100, 5*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			redisKeyPrefix := "consumer_user_setting_lock" + event.ID
			status, can, err := c.redisRepo.Lock(ctx, redisKeyPrefix)
			if err != nil {
				if status == constants.StatusProcessing {
					log.Printf("Event %s is currently being processed by another worker. Skipping.\n", event.ID)
					return err
				}
				return err
			}
			if !can {
				if status == constants.StatusProcessing {
					log.Printf("Event %s is currently being processed by another worker. Skipping.\n", event.ID)
					errchan <- errors.New("event is being processed by another worker")
				} else {
					log.Printf("Event %s has already been processed with status %s. Skipping.\n", event.ID, status)
				}
				continue
			}
			ev := event
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreatedUserSetting(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedUserSetting(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedUserSetting(ctx, ev)
				default:
					log.Printf("Unknown event type: %s\n", ev.Type)
				}
			})
			wg.Add(1)
			if processErr != nil {
				c.redisRepo.Unlock(ctx, redisKeyPrefix)
				errchan <- fmt.Errorf("event %s failed: %w", ev.ID, processErr)
			} else {
				// Thành công -> Giữ khoá để tránh các worker khác xử lý lại
				log.Printf("Event %s processed successfully. Keeping lock to prevent reprocessing.\n", event.ID)
				c.redisRepo.Set(ctx, redisKeyPrefix, constants.StatusProcessing, 30*time.Minute) // Cập nhật lại khoá với status processing và TTL mới
				errchan <- nil
			}
			if err != nil {
				c.redisRepo.Unlock(ctx, redisKeyPrefix)
				errchan <- err
				wg.Done()
			}
		}
		c.pool.Wait()
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
func (c *ConsumerUserSetting) handleCreatedUserSetting(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[req.UserSettingReq](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("invalid parse payload for created user setting event: " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(errors.New("invalid parse payload for created user setting event: " + err.Error()))
	}
	entity := mapper.ToCreateIniEntityUserSetting(data.User_ID)
	err = c.usersettingRepo.CreateUserSettingDefault(ctx, entity)
	if err != nil {
		return errors.New("failed to create default user settings: " + err.Error())
	}
	return nil
}
func (c *ConsumerUserSetting) handleUpdatedUserSetting(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[req.UserSettingReq](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("invalid parse payload for updated user setting event: " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(errors.New("invalid parse payload for updated user setting event: " + err.Error()))
	}
	entity := mapper.ToEntityUserSetting(data)
	err = c.usersettingRepo.UpdateUserSettings(ctx, data.User_ID, entity)
	if err != nil {
		return errors.New("failed to update user settings: " + err.Error())
	}
	return nil
}
func (c *ConsumerUserSetting) handleDeletedUserSetting(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[req.UserSettingReq](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("invalid parse payload for deleted user setting event: " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(errors.New("invalid parse payload for deleted user setting event: " + err.Error()))
	}
	err = c.usersettingRepo.DeleteUserSettings(ctx, data.User_ID)
	if err != nil {
		return errors.New("failed to delete user settings: " + err.Error())
	}
	return nil
}
func (c *ConsumerUserSetting) ConsumerFailedUserSettingEvent(ctx context.Context) error {
	return nil
}
