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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type ConsumerStory struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	storyRepo     IRepositoryMongodb.IStoryRepository
	storyViewRepo IRepositoryCassandra.IStoryViewRepository
	socialGrpc    pb.SocialServiceClient
}

func NewConsumerStory(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, storyRepo IRepositoryMongodb.IStoryRepository, storyViewRepo IRepositoryCassandra.IStoryViewRepository, socialGrpc pb.SocialServiceClient) *ConsumerStory {
	return &ConsumerStory{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		storyRepo:     storyRepo,
		storyViewRepo: storyViewRepo,
		socialGrpc:    socialGrpc,
	}
}
func (c *ConsumerStory) ConsumerStory(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicStory.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				}
			})
			if err != nil {
				errchan <- errors.New("failed to run task for event " + event.ID + ": " + err.Error())
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
func (c *ConsumerStory) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.CreateStoryPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	entity := mapper.ToStoryEntity(data)
	//gọi blocklist
	blockListResp, err := c.socialGrpc.GetListBlockByUserIDV2(ctx, &pb.UserblockIDRequestv2{UserId: data.UserID})
	if err != nil {
		return fmt.Errorf("failed to get block list from social service: %w", err)
	}
	if len(blockListResp.Typeblock) > 0 {
		entity.Privacy.BlockList = append(entity.Privacy.BlockList, blockListResp.Typeblock...)
	}
	err = c.storyRepo.CreateStory(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create story in repository: %w", err)
	}
	//gọi mediaasset
	payloaditemMedia := &mediaEvent.MediaItemPayload{
		MediaID:      entity.ID.Hex(),
		StoryID:      entity.ID.Hex(),
		PostID:       "",
		AlbumID:      "",
		CommentID:    "",
		ReelID:       "",
		GroupID:      "",
		PageID:       "",
		MediaType:    entity.Media.Type,
		URL:          entity.Media.URL,
		ThumbnailURL: entity.Media.ThumbnailURL,
		Metadata: mediaEvent.MetadataPayload{
			Width:     entity.Media.Width,  // Cần bổ sung nếu có thông tin
			Height:    entity.Media.Height, // Cần bổ sung nếu có thông tin
			Duration:  entity.Media.Duration,
			SizeBytes: entity.Media.SizeBytes,
			MimeType:  entity.Media.MimeType, // Cần bổ sung nếu có thông tin
		},
		Order:       0,
		Hashtags:    []string{},                       // Cần bổ sung nếu có thông tin
		TaggedUsers: []mediaEvent.TaggedUserPayload{}, // Cần bổ sung nếu có thông tin
	}
	payloadMediaAsset := &mediaEvent.CreateMediaAssetsPayload{
		UserID: data.UserID,
		Items:  []mediaEvent.MediaItemPayload{*payloaditemMedia},
	}
	err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), payloadMediaAsset)
	if err != nil {
		return fmt.Errorf("failed to publish media asset event: %w", err)
	}
	return nil

}
func (c *ConsumerStory) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.UpdateStoryPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datastory, err := c.storyRepo.GetStoryByID(ctx, data.StoryID)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to get story by ID: %w", err))
	}
	if datastory == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("story not found"))
	}
	datastorynew := mapper.ApplyUpdateStoryPayload(datastory, data)
	err = c.storyRepo.UpdateStory(ctx, datastorynew)
	if err != nil {
		return fmt.Errorf("failed to update story in repository: %w", err)
	}
	return nil
}
func (c *ConsumerStory) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.DeleteStoryPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.storyRepo.DeleteStory(ctx, data.StoryID)
	if err != nil {
		return fmt.Errorf("failed to delete story in repository: %w", err)
	}
	//gọi mediaasset
	payloadMediaAsset := &mediaEvent.DeleteMediaAssetsPayload{
		MediaID: data.StoryID,
	}
	err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), "", constants.Deleted.String(), payloadMediaAsset)
	return nil
}
func (c *ConsumerStory) ConsumerFailedStory(ctx context.Context) error {

	return nil
}
