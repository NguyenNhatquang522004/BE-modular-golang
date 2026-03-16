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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerGroup struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewConsumerGroup(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupRepo IRepositoryMongodb.IGroupRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository) *ConsumerGroup {
	return &ConsumerGroup{
		events:          events,
		pool:            pool,
		redisRepo:       redisRepo,
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
	}
}
func (c *ConsumerGroup) ConsumerGroup(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGroup.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Đảm bảo giảm WaitGroup nếu không thể chạy worker
			}
			if processErr != nil {
				errchan <- processErr
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
func (c *ConsumerGroup) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.CreateGroupPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity := mapper.ToCreateGroupEntity(data)
	if entity == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity"))
	}
	if entity.Avatar.ID != primitive.NilObjectID {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.CreatorID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.Avatar.ID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.Avatar.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.CreatorID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	if entity.Cover.ID != primitive.NilObjectID {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.CreatorID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.Cover.ID.Hex(),
			AlbumID:      "",
			GroupID:      data.GroupID,
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.Cover.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.CreatorID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	err = c.groupRepo.CreateGroup(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	if entity.Avatar.ID != primitive.NilObjectID {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.CreatorID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.Avatar.ID.Hex(),
			AlbumID:      "",
			GroupID:      data.GroupID,
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.Avatar.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.CreatorID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	return nil
}
func (c *ConsumerGroup) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.UpdateGroupPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	// TODO: Implement update logic
	datagroup, err := c.groupRepo.GetGroupByID(ctx, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group by ID: %w", err)
	}
	if data.Avatar != nil && data.Avatar.ID != datagroup.Avatar.ID.Hex() {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserActionID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   datagroup.Avatar.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old avatar: %w", err))
		}
	}

	if data.Cover != nil && data.Cover.ID != datagroup.Cover.ID.Hex() {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserActionID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   datagroup.Cover.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old cover: %w", err))
		}
	}
	if datagroup == nil {
		return fmt.Errorf("group not found with ID: %s", data.GroupID)
	}
	mapper.MapUpdatePayloadToEntity(data, datagroup)
	err = c.groupRepo.UpdateGroup(ctx, datagroup)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}
	if datagroup.Avatar.ID != primitive.NilObjectID {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserActionID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      datagroup.Avatar.ID.Hex(),
			AlbumID:      "",
			GroupID:      data.GroupID,
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          datagroup.Avatar.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserActionID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	if datagroup.Cover.ID != primitive.NilObjectID {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserActionID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      datagroup.Cover.ID.Hex(),
			AlbumID:      "",
			GroupID:      data.GroupID,
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          datagroup.Cover.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserActionID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	return nil
}
func (c *ConsumerGroup) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[communityEvent.DeleteGroupPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	dataaction, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserAction, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group member by user ID and group ID: %w", err)
	}
	if dataaction == nil {
		return fmt.Errorf("group member not found with user ID: %s and group ID: %s", data.UserAction, data.GroupID)
	}
	if dataaction.Role != sharedEnums.RoleTypeAdmin {
		return fmt.Errorf("user with ID: %s is not an admin of group with ID: %s", data.UserAction, data.GroupID)
	}
	datagroup, err := c.groupRepo.GetGroupByID(ctx, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group by ID: %w", err)
	}
	if datagroup == nil {
		return fmt.Errorf("group not found with ID: %s", data.GroupID)
	}
	err = c.groupRepo.DeleteGroup(ctx, data.GroupID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	if datagroup.Avatar.ID != primitive.NilObjectID {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.GroupID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   datagroup.Avatar.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old avatar: %w", err))
		}
	}

	if datagroup.Cover.ID != primitive.NilObjectID {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.GroupID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   datagroup.Cover.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old cover: %w", err))
		}
	}
	return nil
}
func (c *ConsumerGroup) ConsumerFailedGroup(ctx context.Context) error {

	return nil
}
