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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerAlbum struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	albumRepo IRepositoryMongodb.IAlbumsRepository
}

func NewConsumerAlbum(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, albumRepo IRepositoryMongodb.IAlbumsRepository) *ConsumerAlbum {
	return &ConsumerAlbum{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		albumRepo: albumRepo,
	}
}

func (c *ConsumerAlbum) ConsumerAlbum(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicAlbum.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		wg.Wait()
		close(errchan)
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
				// Process the event here
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreatedAlbum(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedAlbum(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedAlbum(ctx, event)
				default:
					processErr = errors.New("unknown event type " + event.Type + " for event " + event.ID)
					return
				}
			})
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
			if err != nil {
				errchan <- errors.New("failed to submit task for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Decrement the WaitGroup counter since the task won't be processed
			}
		}
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
func (c *ConsumerAlbum) handleCreatedAlbum(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.CreateAlbumPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" event payload is nil for event %s", event.ID))
	}
	entity, err := mapper.ToAlbumEntity(data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to map event payload to album entity for event %s: %w", event.ID, err))
	}
	err = c.albumRepo.CreateAlbum(ctx, entity)
	if err != nil {
		return fmt.Errorf(" failed to create album in repository for event %s: %w", event.ID, err)
	}
	if len(data.ItemMediaIDs) > 0 {
		for _, mediaID := range data.ItemMediaIDs {
			payload := mediaEvent.UpdateMediaAssetsPayload{
				MediaID: mediaID,
				AlbumID: entity.ID.Hex(),
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), mediaID, constants.Updated.String(), payload)
		}
	}
	return nil
}
func (c *ConsumerAlbum) handleUpdatedAlbum(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.UpdateAlbumPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" event payload is nil for event %s", event.ID))
	}
	existing, err := c.albumRepo.GetAlbumByID(ctx, data.ID)
	if err != nil {
		return fmt.Errorf(" failed to retrieve existing album from repository for event %s: %w", event.ID, err)
	}
	if existing == nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" album with ID %s not found for event %s", data.ID, event.ID))
	}
	updatedEntity, err := mapper.ApplyAlbumUpdate(*existing, data)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to apply updates to album entity for event %s: %w", event.ID, err))
	}
	err = c.albumRepo.UpdateAlbum(ctx, updatedEntity)
	if err != nil {
		return fmt.Errorf(" failed to update album in repository for event %s: %w", event.ID, err)
	}
	if len(data.ItemMediaDeleteIDs) > 0 {
		for _, mediaID := range data.ItemMediaDeleteIDs {
			payload := mediaEvent.UpdateMediaAssetsPayload{
				MediaID: mediaID,
				AlbumID: "",
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), mediaID, constants.Updated.String(), payload)
		}
	}
	if len(data.ItemMediaAddIDs) > 0 {
		for _, mediaID := range data.ItemMediaAddIDs {
			payload := mediaEvent.UpdateMediaAssetsPayload{
				MediaID: mediaID,
				AlbumID: existing.ID.Hex(),
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), mediaID, constants.Updated.String(), payload)
		}
	}
	return nil
}
func (c *ConsumerAlbum) handleDeletedAlbum(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.DeleteAlbumPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" event payload is nil for event %s", event.ID))
	}
	err = c.albumRepo.DeleteAlbum(ctx, data.AlbumID)
	if err != nil {
		return fmt.Errorf(" failed to delete album in repository for event %s: %w", event.ID, err)
	}
	if len(data.ItemMediaIDs) > 0 {
		for _, mediaID := range data.ItemMediaIDs {
			payload := mediaEvent.UpdateMediaAssetsPayload{
				MediaID: mediaID,
				AlbumID: "",
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), mediaID, constants.Updated.String(), payload)
		}
	}
	return nil
}
func (c *ConsumerAlbum) ConsumerFailedAlbum(ctx context.Context) error {

	return nil
}
