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
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type ConsumerArtist struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	artisRepo IRepositoryMongodb.IArtistRepository
	socicalpb pb.SocialServiceClient
}

func NewConsumerArtist(events events.EventBus,
	pool IRepositoryShare.IWorkerPool,
	redisRepo IRepositoryShare.IRedis,
	artisRepo IRepositoryMongodb.IArtistRepository,
	socicalpb pb.SocialServiceClient) *ConsumerArtist {
	return &ConsumerArtist{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		artisRepo: artisRepo,
		socicalpb: socicalpb,
	}
}
func (c *ConsumerArtist) ConsumerArtist(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicArtist.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				wg.Done() // Ensure we mark the worker as done if we fail to start it
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
func (c *ConsumerArtist) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.CreateArtistPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse CreateArtistPayload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("received nil CreateArtistPayload for event %s", event.ID))
	}
	datainfo, err := c.socicalpb.GetInfoUserByID(ctx, &pb.UserSocialIDRequest{UserId: data.ArtistID})
	if err != nil {
		return fmt.Errorf("failed to get user info for ArtistID %s for event %s: %w", data.ArtistID, event.ID, err)
	}
	if datainfo == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get user info for ArtistID %s for event %s", data.ArtistID, event.ID))
	}
	data.Name = datainfo.AuthorName
	data.AvatarURL = datainfo.AuthorAvatar
	artistEntity := mapper.ToArtistEntity(data)
	err = c.artisRepo.CreatedArtist(ctx, artistEntity)
	if err != nil {
		return fmt.Errorf("failed to create artist in database for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerArtist) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.UpdateArtistPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse UpdateArtistPayload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("received nil UpdateArtistPayload for event %s", event.ID))
	}
	dataartis, err := c.artisRepo.GetArtistByID(ctx, data.ArtisID)
	if err != nil {
		return fmt.Errorf("failed to get existing artist from database for event %s: %w", event.ID, err)
	}
	if dataartis == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("artist with ID %s not found for event %s", data.ArtisID, event.ID))
	}
	entity := mapper.ApplyUpdateToArtist(dataartis, *data)
	err = c.artisRepo.UpdatedArtist(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to update artist in database for event %s: %w", event.ID, err)
	}
	return nil
}
func (c *ConsumerArtist) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.DeleteArtistPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse DeleteArtistPayload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("received nil DeleteArtistPayload for event %s", event.ID))
	}
	err = c.artisRepo.DeletedArtist(ctx, data.ArtisID)
	if err != nil {
		return fmt.Errorf("failed to delete artist from database for event %s: %w", event.ID, err)
	}
	payload := mediaEvent.DeleteMusicLibraryPayload{
		MusicID: "",
		ArtisID: data.ArtisID,
	}
	err = c.events.Publish(ctx, constants.TopicMusic.String(), data.ArtisID, constants.Deleted.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish music deletion event for artist %s after deleting artist for event %s: %w", data.ArtisID, event.ID, err)
	}
	return nil
}
func (c *ConsumerArtist) ConsumerFailedArtist(ctx context.Context) error {

	return nil
}
