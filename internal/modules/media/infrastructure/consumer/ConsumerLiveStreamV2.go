package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerLiveStreamV2 struct {
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	livesession IRepositoryMongodb.ILiveSessionRepository
	seaweedfs   IRepositoryShare.ISeaweedfs
	
}

func NewConsumerLiveStreamV2(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, livesession IRepositoryMongodb.ILiveSessionRepository, seaweedfs IRepositoryShare.ISeaweedfs) *ConsumerLiveStreamV2 {
	return &ConsumerLiveStreamV2{
		events:      events,
		pool:        pool,
		redisRepo:   redisRepo,
		livesession: livesession,
		seaweedfs:   seaweedfs,
	}
}

func (c *ConsumerLiveStreamV2) ConsumeStartLiveStreamV2(ctx context.Context) error {
	// Implement the logic for starting live stream consumption
	err := c.events.SubscribeBatch(ctx, constants.TopicStartStopLiveV2.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				errchan <- errors.New("failed to run event " + event.ID + " in worker pool: " + err.Error())
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
func (c *ConsumerLiveStreamV2) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created event
	data, err := utils.ParsePayload[dto.OMEWebhookPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datalive, err := c.livesession.GetLiveSessionByID(ctx, data.Stream.Name)
	if err != nil {
		return errors.New("failed to get live session from database for ID " + data.Stream.Name + ": " + err.Error())
	}
	if datalive == nil {
		return errors.New("live session not found in database for ID " + data.Stream.Name)
	}
	datalive.Status = sharedEnums.ProcessingActive
	err = c.livesession.UpdateLiveSession(ctx, datalive)
	if err != nil {
		return errors.New("failed to update live session status in database for ID " + data.Stream.Name + ": " + err.Error())
	}
	// Use the parsed data as needed

	return nil
}

func (c *ConsumerLiveStreamV2) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated event
	data, err := utils.ParsePayload[dto.OMEWebhookPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datalive, err := c.livesession.GetLiveSessionByID(ctx, data.Stream.Name)
	if err != nil {
		return errors.New("failed to get live session from database for ID " + data.Stream.Name + ": " + err.Error())
	}
	if datalive == nil {
		return errors.New("live session not found in database for ID " + data.Stream.Name)
	}
	datalive.Status = sharedEnums.ProcessingEnd
	err = c.livesession.UpdateLiveSession(ctx, datalive)
	if err != nil {
		return errors.New("failed to update live session status in database for ID " + data.Stream.Name + ": " + err.Error())
	}
	go c.processVODUpload(context.Background(), data.Stream.Name)
	return nil
}

func (c *ConsumerLiveStreamV2) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted event
	return nil
}
func (c *ConsumerLiveStreamV2) ConsumerFailedTopicLiveStreamV2(ctx context.Context) error {
	// Implement the logic for handling failed topic in live stream
	return nil
}
func (c *ConsumerLiveStreamV2) processVODUpload(ctx context.Context, sessionID string) {
	// Đường dẫn file MP4 mà OME đã record lại (volume /tmp/lives map từ Docker)
	localMP4Path := fmt.Sprintf("/tmp/lives/%s.mp4", sessionID)

	// Đợi 1 chút để đảm bảo OME đã nhả hoàn toàn file lock
	time.Sleep(2 * time.Second)

	fileInfo, err := os.Stat(localMP4Path)
	if err != nil {
		log.Printf("[VOD Error] Không tìm thấy file MP4 của OME: %v", err)
		return
	}

	file, err := os.Open(localMP4Path)
	if err != nil {
		log.Printf("[VOD Error] Không thể mở file MP4: %v", err)
		return
	}
	defer file.Close()
	datalive, err := c.livesession.GetLiveSessionByID(ctx, sessionID)
	if err != nil {
		return
	}
	if datalive == nil {
		return
	}

	// Tái sử dụng Adapter Upload SeaweedFS tuyệt vời của bạn
	uploadInput := &dto.FileUploadInput{
		OwnerID:     "system", // Hoặc query DB lấy OwnerID từ sessionID
		Storage:     utils.BucketLive,
		FileName:    sessionID + "_vod.mp4",
		Content:     file,
		Size:        fileInfo.Size(),
		ContentType: "video/mp4",
	}
	if datalive.PageID != "" {
		uploadInput.OwnerID = datalive.PageID
		uploadInput.Storage = utils.BucketPageLiveStream
	}
	if datalive.GroupID != "" {
		uploadInput.OwnerID = datalive.GroupID
		uploadInput.Storage = utils.BucketGroupStream
	}
	if datalive.GroupID == "" && datalive.PageID == "" {
		uploadInput.OwnerID = datalive.HostUserID
		uploadInput.Storage = utils.BucketLive
	}
	uploadResult, err := c.seaweedfs.Upload(ctx, uploadInput)
	if err != nil {
		log.Printf("[VOD Error] Đẩy file lên SeaweedFS thất bại: %v", err)
		return
	}

	log.Printf("[VOD Success] Đã đẩy VOD lên SeaweedFS: %s", uploadResult.PublicURL)
	data, err := c.livesession.GetLiveSessionByID(ctx, sessionID)
	if err != nil {
		log.Printf("[VOD Error] Lỗi khi lấy live session từ DB: %v", err)
		return
	}
	if data == nil {
		log.Printf("[VOD Error] Không tìm thấy live session trong DB cho ID: %s", sessionID)
		return
	}
	// Cập nhật link VOD vào Database
	// _ = c.livesession.UpdatePlaybackURL(ctx, sessionID, uploadResult.PublicURL)

	// Dọn rác: Xóa file mp4 local sau khi đã đẩy an toàn lên SeaweedFS
	_ = os.Remove(localMP4Path)
	_ = os.Remove(fmt.Sprintf("/tmp/lives/%s.xml", sessionID)) // Xóa luôn file info của OME
}
