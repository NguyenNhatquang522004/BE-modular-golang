package consumer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerMediaAssets struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	mediaRepo IRepositoryMongodb.IMediaAssetsRepository
}

func NewConsumerMediaAssets(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, mediaRepo IRepositoryMongodb.IMediaAssetsRepository) *ConsumerMediaAssets {
	return &ConsumerMediaAssets{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		mediaRepo: mediaRepo,
	}
}

func (c *ConsumerMediaAssets) ConsumeMediaAsset(ctx context.Context) error {
	// Implement the logic for consuming a media asset
	err := c.events.SubscribeBatch(ctx, constants.TopicMediaAsset.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, acquired, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- errors.New("failed to acquire lock for event " + event.ID + ": " + err.Error())
				continue
			}
			if !acquired {
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
					processErr = c.handlerCreatedMediaAsset(ctx, ev)
				case constants.Updated.String():
					processErr = c.handlerUpdatedMediaAsset(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handlerDeletedMediaAsset(ctx, ev)
				default:
					errchan <- errors.New("unknown event type: " + ev.Type)
				}
			})
			if processErr != nil {
				errchan <- errors.New("failed to process event " + ev.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, ev.ID)

			} else {
				c.redisRepo.MarkCompleted(ctx, ev.ID)
				errchan <- nil
			}
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + ev.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, ev.ID)
				wg.Done()
			}
		}
		wg.Wait()
		close(errchan)
		var combinedErr error
		for err := range errchan {
			if err != nil {
				combinedErr = errors.Join(combinedErr, err) // Hoặc bạn có thể kết hợp nhiều lỗi lại với nhau
			}
		}
		return combinedErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerMediaAssets) handlerCreatedMediaAsset(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling a created media asset
	data, err := utils.ParsePayload[mediaEvent.CreateMediaAssetsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("failed to parse event payload: " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(errors.New("event payload is nil"))
	}
	entity, err := mapper.ToEntityCreateMediaAssetsPayload(data)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("failed to map event payload to entity: " + err.Error()))
	}
	if len(entity) == 0 {
		return kafka.NewNonRetryableError(errors.New("mapped entity is empty"))
	}
	_, _, err = c.mediaRepo.CreateBulkMediaAssets(ctx, entity)
	if err != nil {
		return errors.New("failed to create media assets in repository: " + err.Error())
	}
	return nil
}

func (c *ConsumerMediaAssets) handlerUpdatedMediaAsset(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling an updated media asset
	data, err := utils.ParsePayload[mediaEvent.UpdateMediaAssetsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("failed to parse event payload: " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(errors.New("event payload is nil"))
	}
	datamediasset, err := c.mediaRepo.GetMediaAssetByID(ctx, data.MediaID)
	if err != nil {
		return errors.New("failed to get media asset from repository: " + err.Error())
	}
	if datamediasset == nil {
		return kafka.NewNonRetryableError(errors.New("media asset not found with ID: " + data.MediaID))
	}
	err = mapper.UpdateEntityMediaAssetsFromPayload(datamediasset, data)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("failed to update media asset entity: " + err.Error()))
	}
	err = c.mediaRepo.UpdateMediaAsset(ctx, datamediasset)
	if err != nil {
		return errors.New("failed to update media asset in repository: " + err.Error())
	}
	// Nếu có logic liên quan đến việc cập nhật album hoặc group khi media asset thay đổi, bạn có thể thêm vào đây.
	return nil
}

func (c *ConsumerMediaAssets) handlerDeletedMediaAsset(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling a deleted media asset
	data, err := utils.ParsePayload[mediaEvent.DeleteMediaAssetsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(errors.New("failed to parse event payload: " + err.Error()))
	}
	if data == nil {
		return kafka.NewNonRetryableError(errors.New("event payload is nil"))
	}
	if data.MediaID != "" {
		err = c.mediaRepo.DeleteMediaAsset(ctx, data.MediaID)
		if err != nil {
			return errors.New("failed to delete media asset from repository: " + err.Error())
		}
	}
	if data.MessageID != "" {
		err = c.mediaRepo.DeleteMediaAssetsByMessageID(ctx, data.MessageID)
		if err != nil {
			return errors.New("failed to delete media assets by message ID from repository: " + err.Error())
		}
	}

	return nil
}
func (c *ConsumerMediaAssets) ConsumerFailedMediaAsset(ctx context.Context) error {
	// Implement the logic for handling a failed media asset
	return nil
}
