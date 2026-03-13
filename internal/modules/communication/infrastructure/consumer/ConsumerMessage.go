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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
)

type ConsumerMessage struct {
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	messageRepo IRepositoryCassandra.IMessageRepository
}

func NewConsumerMessage(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, messageRepo IRepositoryCassandra.IMessageRepository) *ConsumerMessage {
	return &ConsumerMessage{
		events:      events,
		pool:        pool,
		redisRepo:   redisRepo,
		messageRepo: messageRepo,
	}
}
func (c *ConsumerMessage) ConsumerMessage(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicMessage.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = c.handleCreatedMessage(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedMessage(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedMessage(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
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
		wg.Wait()
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
func (c *ConsumerMessage) handleCreatedMessage(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi nhận được sự kiện Created
	data, err := utils.ParsePayload[communicationEvent.CreateMessagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	messageEntity, ids := mapper.ToMessageEntity(data)
	err = c.messageRepo.CreateMessage(ctx, messageEntity)
	if err != nil {
		return fmt.Errorf("failed to create message in repository: %w", err)
	}
	if len(ids) > 0 {
		for index, item := range ids {
			payload := &mediaEvent.CreateMediaAssetsPayload{
				UserID: data.SenderID.String(),
			}
			payload.Items = append(payload.Items, mediaEvent.MediaItemPayload{
				MediaID:      item,
				PostID:       "",
				AlbumID:      "",
				GroupID:      "",
				CommentID:    "",
				PageID:       "",
				StoryID:      "",
				ReelID:       "",
				MessageID:    data.MessageID,
				MediaType:    data.Attachments[index].Type,
				URL:          data.Attachments[index].URL,
				ThumbnailURL: data.Attachments[index].ThumbnailURL,
				Metadata: mediaEvent.MetadataPayload{
					Width:     data.Attachments[index].Width,
					Height:    data.Attachments[index].Height,
					SizeBytes: data.Attachments[index].SizeBytes,
					MimeType:  data.Attachments[index].MimeType,
				},
				Order:       0,
				Hashtags:    nil,
				TaggedUsers: nil,
			})
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.MessageID, constants.Created.String(), payload)
			if err != nil {
				return fmt.Errorf("failed to publish media asset event for message %s: %w", data.MessageID, err)
			}
		}

	}
	payload := &communicationEvent.ConversationStatsPayload{
		UserID:           nil,
		ConversationID:   data.ConversationID,
		ParticipantCount: nil,
		LastMessage: &communicationEvent.LastMessageCachePayload{
			MessageID: data.MessageID,
			Content:   data.Content,
			SenderID:  data.SenderID.String(),
			Type:      data.Attachments[0].Type,
			CreatedAt: time.Now(),
		},
		LastSeenAt:        nil,
		LastSeenMessageID: nil,
		EventType:         constants.Created,
	}
	err = c.events.Publish(ctx, constants.TopicStatsConversation.String(), data.ConversationID, constants.Created.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish conversation stats event for message %s: %w", data.MessageID, err)
	}
	return nil
}
func (c *ConsumerMessage) handleUpdatedMessage(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi nhận được sự kiện Updated
	data, err := utils.ParsePayload[communicationEvent.UpdateMessagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	existingMessage, err := c.messageRepo.GetMessagesByMessageID(ctx, data.ConversationID, data.Bucket, data.MessageID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing message: %w", err)
	}
	if existingMessage == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("message not found for update"))
	}
	entity, addedAssets, removedIDs := mapper.UpdateMessageMapper(existingMessage, data)
	err = c.messageRepo.UpdateMessage(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to update message in repository: %w", err)
	}
	// Xử lý media assets đã thêm
	for _, asset := range addedAssets {
		payload := &mediaEvent.UpdateMediaAssetsPayload{
			MediaID:      asset.AssetID,
			AlbumID:      "",
			URL:          &asset.URL,
			ThumbnailURL: &asset.ThumbnailURL,
			Metadata: &mediaEvent.MetadataPayload{
				Width:     asset.Width,
				Height:    asset.Height,
				SizeBytes: asset.SizeBytes,
				MimeType:  asset.MimeType,
			},
			Order:       nil,
			Hashtags:    nil,
			TaggedUsers: nil,
		}

		err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.MessageID, constants.Updated.String(), payload)
		if err != nil {
			return fmt.Errorf("failed to publish media asset event for added asset in message %s: %w", data.MessageID, err)
		}
	}
	// Xử lý media assets đã xóa
	for _, removedID := range removedIDs {
		payload := &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   removedID,
			MessageID: data.MessageID,
		}
		err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.MessageID, constants.Deleted.String(), payload)
		if err != nil {
			return fmt.Errorf("failed to publish media asset event for removed asset in message %s: %w", data.MessageID, err)
		}
	}

	return nil
}
func (c *ConsumerMessage) handleDeletedMessage(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi nhận được sự kiện Deleted
	data, err := utils.ParsePayload[communicationEvent.DeleteMessagePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.DeleteAlll {
		err = c.messageRepo.DeleteAllMessagesByMessageConversationID(ctx, data.ConversationID)
		if err != nil {
			return fmt.Errorf("failed to delete all messages in conversation %s: %w", data.ConversationID, err)
		}
		return nil
	}
	// Xử lý xóa message trong repository
	err = c.messageRepo.DeleteMessage(ctx, data.ConversationID, data.Bucket, data.MessageID)
	if err != nil {
		return fmt.Errorf("failed to delete message in repository: %w", err)
	}
	return nil
}
func (c *ConsumerMessage) ConsumerFailedMessage(ctx context.Context) error {
	return nil
}
