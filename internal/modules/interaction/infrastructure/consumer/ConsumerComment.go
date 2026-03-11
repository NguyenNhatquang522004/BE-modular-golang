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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type ConsumerComment struct {
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	commentRepo IRepositoryMongoDB.ICommentRepository
}

func NewConsumerComment(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, commentRepo IRepositoryMongoDB.ICommentRepository) *ConsumerComment {
	return &ConsumerComment{
		events:      events,
		pool:        pool,
		redisRepo:   redisRepo,
		commentRepo: commentRepo,
	}
}

func (c *ConsumerComment) ConsumerComment(ctx context.Context) error {
	// Implement the logic for consuming comments here
	err := c.events.SubscribeBatch(ctx, constants.TopicComment.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedComment(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedComment(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedComment(ctx, event)
				default:
					processErr = errors.New("unknown event type: " + event.Type)
				}
			})
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
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
		wg.Done()
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
func (c *ConsumerComment) handleCreatedComment(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling created comments here
	data, err := utils.ParsePayload[interactionEvent.CreatedCommentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	entitycomment, err := mapper.MapCreatedCommentPayloadToEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	err = c.commentRepo.CreateComment(ctx, entitycomment)
	if err != nil {
		return fmt.Errorf("failed to create comment in repository for event %s: %w", event.ID, err)
	}
	if data.Media != nil {
		var itemmedia = []mediaEvent.MediaItemPayload{}
		item := mediaEvent.MediaItemPayload{
			MediaID:      entitycomment.AssetID.Hex(),
			CommentID:    entitycomment.ID.Hex(),
			PostID:       "",
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			MediaType:    data.Media.Type,
			URL:          data.Media.URL,
			ThumbnailURL: "",
			Order:        0,
			Hashtags:     nil,
			Metadata: mediaEvent.MetadataPayload{
				Width:     data.Media.DisplayMeta.Width,
				Height:    data.Media.DisplayMeta.Height,
				Duration:  data.Media.DisplayMeta.Duration,
				SizeBytes: data.Media.DisplayMeta.SizeBytes,
				MimeType:  data.Media.DisplayMeta.MimeType,
			},
			TaggedUsers: func() []mediaEvent.TaggedUserPayload {
				var taggedUsers []mediaEvent.TaggedUserPayload
				if data.Mentions != nil {
					for _, mention := range *data.Mentions {
						taggedUsers = append(taggedUsers, mediaEvent.TaggedUserPayload{
							UserID: mention,
							Name:   "", // Nếu có tên người dùng, bạn có thể thêm vào đây
							X:      0,  // Vị trí X nếu có
							Y:      0,  // Vị trí Y nếu có
						})
					}
				}
				return taggedUsers
			}(),
		}
		itemmedia = append(itemmedia, item)
		payload := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
			Items:  itemmedia,
		}
		err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), entitycomment.ID.Hex(), constants.Created.String(), payload)
		if err != nil {
			return fmt.Errorf("failed to publish media asset event for comment %s: %w", entitycomment.ID.Hex(), err)
		}
	}
	return nil
}

func (c *ConsumerComment) handleUpdatedComment(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling updated comments here
	data, err := utils.ParsePayload[interactionEvent.UpdatedCommentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataComment, err := c.commentRepo.GetCommentByID(ctx, data.CommentID)
	if err != nil {
		return fmt.Errorf("failed to get comment by ID %s for event %s: %w", data.CommentID, event.ID, err)
	}
	if dataComment == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("comment with ID %s not found for event %s", data.CommentID, event.ID))
	}
	entitycomment, err := mapper.MapUpdatedCommentPayloadToEntity(data, dataComment)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to map payload to entity for event %s: %w", event.ID, err))
	}
	err = c.commentRepo.UpdateComment(ctx, entitycomment)
	if err != nil {
		return fmt.Errorf("failed to update comment in repository for event %s: %w", event.ID, err)
	}
	if data.Media != nil {
		itemupdate := &mediaEvent.UpdateMediaAssetsPayload{}
		itemupdate.MediaID = entitycomment.AssetID.Hex()
		itemupdate.URL = &data.Media.URL
		itemupdate.ThumbnailURL = nil
		itemupdate.Order = nil
		itemupdate.Hashtags = nil
		itemupdate.TaggedUsers = func() *[]mediaEvent.TaggedUserPayload {
			if data.Mentions != nil {
				var taggedUsers []mediaEvent.TaggedUserPayload
				for _, mention := range *data.Mentions {
					taggedUsers = append(taggedUsers, mediaEvent.TaggedUserPayload{
						UserID: mention,
						Name:   "", // Nếu có tên người dùng, bạn có thể thêm vào đây
						X:      0,  // Vị trí X nếu có
						Y:      0,  // Vị trí Y nếu có
					})
				}
				return &taggedUsers
			}
			return nil
		}()
		itemupdate.Metadata = &mediaEvent.MetadataPayload{
			Width:     data.Media.DisplayMeta.Width,
			Height:    data.Media.DisplayMeta.Height,
			Duration:  data.Media.DisplayMeta.Duration,
			SizeBytes: data.Media.DisplayMeta.SizeBytes,
			MimeType:  data.Media.DisplayMeta.MimeType,
		}
		err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), entitycomment.ID.Hex(), constants.Updated.String(), itemupdate)
		if err != nil {
			return fmt.Errorf("failed to publish media asset update event for comment %s: %w", entitycomment.ID.Hex(), err)
		}
	}
	return nil
}

func (c *ConsumerComment) handleDeletedComment(ctx context.Context, event events.IntegrationEvent) error {
	// Implement the logic for handling deleted comments here
	data, err := utils.ParsePayload[interactionEvent.UpdatedCommentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil for event %s", event.ID))
	}
	dataComment, err := c.commentRepo.GetCommentByID(ctx, data.CommentID)
	if err != nil {
		return fmt.Errorf("failed to get comment by ID %s for event %s: %w", data.CommentID, event.ID, err)
	}
	if dataComment == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("comment with ID %s not found for event %s", data.CommentID, event.ID))
	}
	err = c.commentRepo.DeleteComment(ctx, data.CommentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment in repository for event %s: %w", event.ID, err)
	}
	if data.Media != nil && dataComment.AssetID != nil {
		err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.CommentID, constants.Deleted.String(), mediaEvent.DeleteMediaRelationTargetPayload{
			TargetID: data.CommentID,
		})
		if err != nil {
			return fmt.Errorf("failed to publish media asset delete event for comment %s: %w", data.CommentID, err)
		}
	}
	return nil
}
func (c *ConsumerComment) ConsumerFailedComment(ctx context.Context) error {
	// Implement the logic for handling failed comments here
	return nil
}
