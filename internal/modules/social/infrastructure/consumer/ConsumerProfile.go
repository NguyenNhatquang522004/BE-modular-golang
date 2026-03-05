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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type ConsumerProfile struct {
	profileRepo IRepositoryMongodb.IProfileRepositoryMongodb
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	pb          pb.IdentityServiceClient
}

func NewConsumerProfile(profileRepo IRepositoryMongodb.IProfileRepositoryMongodb, events events.EventBus, pbClient pb.IdentityServiceClient) *ConsumerProfile {
	return &ConsumerProfile{
		profileRepo: profileRepo,
		events:      events,
		pb:          pbClient,
	}
}

func (c *ConsumerProfile) ConsumerProfile(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicProfile.String(), 100, time.Duration(5)*time.Second, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			redisKeyPrefix := "consumer_profile_lock:" + event.ID
			// 1. Thử khóa event này trong Redis để đảm bảo chỉ 1 worker xử lý 1 eventID nhất định (Distributed Lock)
			status, acquired, err := c.redisRepo.Lock(ctx, redisKeyPrefix)
			if err != nil {
				if status != constants.StatusProcessing {
					log.Printf("Event %s is already processed with status %s. Skipping.\n", event.ID, status)
					return err
				}
				return err
			}
			if !acquired {
				if status == constants.StatusProcessing {
					log.Printf("Event %s is currently being processed by another worker. Skipping.\n", event.ID)
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
					processErr = c.handleCreatedProfile(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedProfile(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedProfile(ctx, ev)
				default:
					log.Printf("Unsupported event type %s for event ID %s. Marking as failed.\n", ev.Type, ev.ID)

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
	return err
}

func (c *ConsumerProfile) handleCreatedProfile(ctx context.Context, event events.IntegrationEvent) error {
	// Gọi hàm parse chuẩn
	data, err := utils.ParsePayload[socialEvent.ProfilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	entity, err := mapper.ToEntityProfilePayload(data)
	if err != nil {
		return kafka.NewNonRetryableError(err) // Lỗi mapper thường là lỗi data sai, cũng ném vào DLQ
	}
	datausersetting, err := c.pb.GetUserSettingByID(ctx, &pb.UserSettingIDRequest{UserId: data.UserID})
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	entity.Settings.AllowSearchEngine = datausersetting.AllowSearchEngineIndexing
	entity.Settings.IsPrivate = datausersetting.Allow_Profile_View_From
	// Lỗi DB thì cứ trả về bình thường để Retry
	return c.profileRepo.CreateProfile(ctx, entity)
}
func (c *ConsumerProfile) handleUpdatedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[socialEvent.ProfilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	dataprofile, err := c.profileRepo.GetProfileByID(ctx, data.UserID)
	if err != nil {
		return err // Lỗi kết nối DB -> Retry
	}
	datausersetting, err := c.pb.GetUserSettingByID(ctx, &pb.UserSettingIDRequest{UserId: data.UserID})
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	dataprofile.Settings.AllowSearchEngine = datausersetting.AllowSearchEngineIndexing
	dataprofile.Settings.IsPrivate = datausersetting.Allow_Profile_View_From
	if dataprofile == nil {
		// CẬP NHẬT: Tuỳ vào logic nghiệp vụ của bạn.
		// Lệnh Update mà user không tồn tại thì không bao giờ thành công được -> NonRetryableError
		return kafka.NewNonRetryableError(errors.New("profile not found for update"))
	}
	entity, err := mapper.ToEntityUpdateProfilePayload(dataprofile, data)
	err = c.profileRepo.UpdateProfile(ctx, entity)
	if err != nil {
		return err
	}
	return c.profileRepo.UpdateProfile(ctx, entity)
}
func (c *ConsumerProfile) handleDeletedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[socialEvent.ProfilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	err = c.profileRepo.DeleteProfile(ctx, data.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerProfile) ConsumerFailedProfile(ctx context.Context) error {
	return nil
}
