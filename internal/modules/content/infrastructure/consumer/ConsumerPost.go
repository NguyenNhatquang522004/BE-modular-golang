package content

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type ConsumerPost struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	postRepo      IRepositoryMongodb.IPostRepository
	mediaRepo     IRepositoryMongodb.IPostMediaRepository
	extensionRepo IRepositoryMongodb.IPostExtensionRepository
	settingRepo   IRepositoryMongodb.IPostSettingRepository
	insightRepo   IRepositoryCassandra.IPostInsights
}

func NewConsumerPost(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, postRepo IRepositoryMongodb.IPostRepository, mediaRepo IRepositoryMongodb.IPostMediaRepository, extensionRepo IRepositoryMongodb.IPostExtensionRepository, settingRepo IRepositoryMongodb.IPostSettingRepository, insightRepo IRepositoryCassandra.IPostInsights) *ConsumerPost {
	return &ConsumerPost{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		postRepo:      postRepo,
		mediaRepo:     mediaRepo,
		extensionRepo: extensionRepo,
		settingRepo:   settingRepo,
		insightRepo:   insightRepo,
	}
}

func (c *ConsumerPost) ConsumerPost(ctx context.Context) error {
	// Xử lý logic tiêu thụ sự kiện liên quan đến bài viết
	err := c.events.SubscribeBatch(ctx, constants.TopicPost.String(), 100, time.Duration(5)*time.Second, func(ctx context.Context, events []events.IntegrationEvent) error {
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
			ev := event
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreatedPost(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedPost(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedPost(ctx, ev)
				default:
					processErr = errors.New("unknown event type: " + ev.Type)
				}
			})
			if processErr != nil {
				c.redisRepo.Unlock(ctx, ev.ID) // Thất bại -> Mở khoá để lần sau làm lại
				errchan <- errors.New("event " + ev.ID + " failed: " + processErr.Error())
			} else {
				// Thành công -> Cập nhật trạng thái đã xử lý trong Redis
				c.redisRepo.MarkCompleted(ctx, ev.ID)
				errchan <- nil
			}
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + ev.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, ev.ID) // Mở khoá nếu không chạy được worker
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
func (c *ConsumerPost) handleCreatedPost(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi tiêu thụ sự kiện tạo bài viết, ví dụ: insert bài viết vào DB, xử lý media, extension/setting, v.v.
	data, err := utils.ParsePayload[contentEvent.CreatePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}

	return nil
}

func (c *ConsumerPost) handleUpdatedPost(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi tiêu thụ sự kiện cập nhật bài viết, ví dụ: cập nhật nội dung bài viết trong DB, xử lý thay đổi media/extension/setting, v.v.
	data, err := utils.ParsePayload[contentEvent.UpdatePostReq](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	return nil
}
func (c *ConsumerPost) handleDeletedPost(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi tiêu thụ sự kiện xoá bài viết, ví dụ: xoá media liên quan, xoá extension/setting, v.v.
	data, err := utils.ParsePayload[contentEvent.DeletePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	return nil
}
func (c *ConsumerPost) ConsumerFailedPost(ctx context.Context) error {
	// Xử lý logic khi tiêu thụ sự kiện thất bại, ví dụ: ghi log, retry, v.v.
	return nil
}
