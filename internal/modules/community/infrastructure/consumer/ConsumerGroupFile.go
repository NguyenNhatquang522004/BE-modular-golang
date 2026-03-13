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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerGroupFile struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	groupFile IRepositoryMongodb.IGroupFilesRepository
}

func NewConsumerGroupFile(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupFile IRepositoryMongodb.IGroupFilesRepository) *ConsumerGroupFile {
	return &ConsumerGroupFile{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		groupFile: groupFile,
	}
}
func (c *ConsumerGroupFile) ConsumerGroupFile(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGroupFile.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedGroupFile(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedGroupFile(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedGroupFile(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
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
func (c *ConsumerGroupFile) handleCreatedGroupFile(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi nhận được sự kiện tạo mới group file
	data, err := utils.ParsePayload[communityEvent.CreateGroupFilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity, err := mapper.MapCreatePayloadToEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	if entity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("mapped entity is nil"))
	}
	err = c.groupFile.CreateGroupFile(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create group file in database for event %s: %w", event.ID, err)
	}
	payload := &mediaEvent.CreateMediaAssetsPayload{
		UserID: data.UploaderID,
		Items: []mediaEvent.MediaItemPayload{
			{
				MediaID:      entity.ID.Hex(),
				PostID:       *data.PostID,
				GroupID:      data.GroupID,
				CommentID:    *data.CommentID,
				MessageID:    *data.MessageID,
				MediaType:    entity.FileType,
				URL:          entity.OriginalURL,
				ThumbnailURL: entity.ThumbnailURL,
				Metadata: mediaEvent.MetadataPayload{
					Width:     *entity.Metadata.Width,
					Height:    *entity.Metadata.Height,
					Duration:  *entity.Metadata.Duration,
					SizeBytes: entity.Metadata.SizeBytes,
					MimeType:  entity.Metadata.MimeType,
				},
				Order:       0,                                // Có thể điều chỉnh nếu cần thiết
				Hashtags:    []string{},                       // Có thể trích xuất từ caption nếu cần thiết
				TaggedUsers: []mediaEvent.TaggedUserPayload{}, // Có thể trích xuất từ caption nếu cần thiết
			},
		},
	}
	err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.GroupID, constants.Created.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish media asset creation event for group file %s: %w", entity.ID.Hex(), err)
	}

	// Xử lý dữ liệu từ sự kiện
	return nil
}
func (c *ConsumerGroupFile) handleUpdatedGroupFile(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi nhận được sự kiện cập nhật group file
	data, err := utils.ParsePayload[communityEvent.UpdateGroupFilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	existingFile, err := c.groupFile.GetGroupFileByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing group file from database for event %s: %w", event.ID, err)
	}
	if existingFile == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("group file with ID %s not found for event %s", data.ID, event.ID))
	}
	updatedEntity := mapper.MapUpdatePayloadToEntityGroupFile(existingFile, data)
	err = c.groupFile.UpdateGroupFile(ctx, updatedEntity)
	if err != nil {
		return fmt.Errorf("failed to update group file in database for event %s: %w", event.ID, err)
	}
	// Xử lý dữ liệu từ sự kiện
	return nil
}
func (c *ConsumerGroupFile) handleDeletedGroupFile(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi nhận được sự kiện xóa group file
	data, err := utils.ParsePayload[communityEvent.DeleteGroupFilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.DeleteAll == true {
		err = c.groupFile.DeleteGroupFilesByGroupID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to delete group files by group ID in database for event %s: %w", event.ID, err)
		}
	} else {
		err = c.groupFile.DeleteGroupFile(ctx, data.ID)
		if err != nil {
			return fmt.Errorf("failed to delete group file in database for event %s: %w", event.ID, err)
		}
	}
	payload := &mediaEvent.DeleteMediaAssetsPayload{
		MediaID:   data.ID,
		MessageID: "",
	}
	err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.GroupID, constants.Deleted.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish media asset creation event for group file %s: %w", data.ID, err)
	}

	return nil
}
func (c *ConsumerGroupFile) ConsumerFailedGroupFile(ctx context.Context) error {
	return nil
}
