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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerUserBookmark struct {
	events       events.EventBus
	pool         IRepositoryShare.IWorkerPool
	redisRepo    IRepositoryShare.IRedis
	bookmarkRepo IRepositoryMongoDB.ISavedItemsRepository
}

func NewConsumerUserBookmark(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, bookmarkRepo IRepositoryMongoDB.ISavedItemsRepository) *ConsumerUserBookmark {
	return &ConsumerUserBookmark{
		events:       events,
		pool:         pool,
		redisRepo:    redisRepo,
		bookmarkRepo: bookmarkRepo,
	}
}
func (c *ConsumerUserBookmark) ConsumerUserBookmark(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicUserBookmark.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchannel := make(chan error, len(events))
		var finalErr error
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchannel <- err
				continue
			}
			if !can {
				if status == constants.StatusProcessing {
					errchannel <- errors.New("event " + event.ID + " is currently being processed by another worker. Skipping.")
				} else {
					errchannel <- errors.New("event " + event.ID + " has already been processed with status " + status.String() + ". Skipping.")
				}
				continue
			}
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch event.Type {
				case constants.Created.String():
					processErr = c.HandleBookmarkCreated(ctx, event)
				case constants.Updated.String():
					processErr = c.HandleBookmarkUpdated(ctx, event)
				case constants.Deleted.String():
					processErr = c.HandleBookmarkDeleted(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
				}
			})
			if processErr != nil {
				errchannel <- processErr
				c.redisRepo.Unlock(ctx, event.ID)

			} else {
				errchannel <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchannel <- err
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
			}
		}
		wg.Wait()
		close(errchannel)
		for err := range errchannel {
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
func (c *ConsumerUserBookmark) HandleBookmarkCreated(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.BookmarkPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	convertTargetID, err := primitive.ObjectIDFromHex(data.TargetID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to convert target ID: %w", err))
	}
	entity := &entity.UserSavedItem{
		UserID:     data.UserID,
		TargetID:   convertTargetID,
		TargetType: data.TargetType,
		Snapshot: entity.SavedItemSnapshot{
			AuthorName:     data.Snapshot.AuthorName,
			ContentPreview: data.Snapshot.ContentPreview,
			ThumbnailURL:   data.Snapshot.ThumbnailURL,
		},
		CollectionName: data.CollectionName,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
	err = c.bookmarkRepo.CreateSaveItem(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create bookmark in database: %w", err)
	}
	return nil
}
func (c *ConsumerUserBookmark) HandleBookmarkUpdated(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.BookmarkPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	databookmark, err := c.bookmarkRepo.GetSaveItemByTargetID(ctx, data.UserID, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get existing bookmark: %w", err)
	}
	if databookmark == nil {
		return fmt.Errorf("bookmark not found for user %s and target %s", data.UserID, data.TargetID)
	}
	databookmark.Snapshot = entity.SavedItemSnapshot{
		AuthorName:     data.Snapshot.AuthorName,
		ContentPreview: data.Snapshot.ContentPreview,
		ThumbnailURL:   data.Snapshot.ThumbnailURL,
	}
	databookmark.CollectionName = data.CollectionName
	databookmark.UpdatedAt = data.UpdatedAt

	err = c.bookmarkRepo.UpdateSaveItem(ctx, databookmark)
	if err != nil {
		return fmt.Errorf("failed to update bookmark in database: %w", err)
	}

	return nil
}
func (c *ConsumerUserBookmark) HandleBookmarkDeleted(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[interactionEvent.BookmarkPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.bookmarkRepo.DeleteSaveItemByUserIDAndTargetID(ctx, data.UserID, data.TargetID)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark in database: %w", err)
	}
	return nil
}
func (c *ConsumerUserBookmark) ConsumerFailedUserBookmark(ctx context.Context) error {
	return nil
}
