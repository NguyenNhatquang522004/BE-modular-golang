package consumer

import (
	"context"
	"errors"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepostitoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/utils"
)

type ConsumerReact struct {
	// Define any dependencies needed for consuming react events here
	ablumRepo       IRepostitoryMongodb.IAlbumsRepository
	mediaAssetsRepo IRepostitoryMongodb.IMediaAssetsRepository
	storyRepo       IRepostitoryMongodb.IStoryRepository
	storyview       IRepositoryCassandra.IStoryViewRepository
	pool            IRepositoryShare.IWorkerPool
	eventbus        events.EventBus
}

func NewConsumerReact(ablumRepo IRepostitoryMongodb.IAlbumsRepository, mediaAssetsRepo IRepostitoryMongodb.IMediaAssetsRepository, storyRepo IRepostitoryMongodb.IStoryRepository, storyview IRepositoryCassandra.IStoryViewRepository, pool IRepositoryShare.IWorkerPool, eventbus events.EventBus) *ConsumerReact {
	return &ConsumerReact{
		ablumRepo:       ablumRepo,
		mediaAssetsRepo: mediaAssetsRepo,
		storyRepo:       storyRepo,
		storyview:       storyview,
		pool:            pool,
		eventbus:        eventbus,
	}
}

// Implement the methods defined in the IConsumerReact interface here
func (c *ConsumerReact) ConsumerReactAlbum(ctx context.Context) {
	// Implement the logic for consuming react album events here
	err := c.eventbus.Subscribe(ctx, constants.TopicReactAlbum.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*req.ReactAlbumRequest)
		if !ok {
			// Handle type assertion error
			return nil
		}
		switch event.Type {
		case constants.Created.String():
			// Process the react album request here
			dataAlbum, err := c.ablumRepo.GetAlbumByID(ctx, data.AlbumID)
			if err != nil {
				// Handle error
				return nil
			}
			utils.CreateReactAlbumRequestToAlbumReactionStats(dataAlbum, data.ReactionType)
			err = c.ablumRepo.UpdateAlbum(ctx, dataAlbum)
			if err != nil {
				// Handle error
				return nil
			}

		case constants.Deleted.String():
			// Process the react album request here
			dataAlbum, err := c.ablumRepo.GetAlbumByID(ctx, data.AlbumID)
			if err != nil {
				// Handle error
				return nil
			}
			utils.DeleteReactAlbumRequestToAlbumReactionStats(dataAlbum, data.ReactionType)
			err = c.ablumRepo.UpdateAlbum(ctx, dataAlbum)
			if err != nil {
				// Handle error
				return nil
			}
		default:
			// Handle unknown event type
		}
		// Process the react album request here
		return nil
	})
	if err != nil {
		// Handle subscription error
	}

}
func (c *ConsumerReact) ConsumerFailedReactAlbum(ctx context.Context) {
	// Implement the logic for consuming failed react album events here
}
func (c *ConsumerReact) CosumerReactStory(ctx context.Context) {

	workerCount := 2 // Adjust the number of workers as needed
	resultChan := make(chan res.FailedConsumerReactStoryResponse, workerCount)
	// Implement the logic for consuming react story events here
	err := c.eventbus.Subscribe(ctx, constants.TopicReactStory.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*req.ReactStoryRequest)
		if !ok {
			// Handle type assertion error
			resultChan <- res.FailedConsumerReactStoryResponse{
				StoryId:         "",
				UserId:          "",
				Avatar:          "",
				Name:            "",
				ViewedAt:        data.ViewedAt,
				InteractionType: data.InteractionType,
				ReactionCode:    0,
				PollOptionIndex: nil,
				EventType:       data.EventType,
				ErrorMessage:    "invalid event payload",
			}
			return errors.New("invalid event payload")
		}
		if event.Type != constants.Created.String() {
			resultChan <- res.FailedConsumerReactStoryResponse{
				StoryId:         "",
				UserId:          "",
				Avatar:          "",
				Name:            "",
				ViewedAt:        data.ViewedAt,
				InteractionType: data.InteractionType,
				ReactionCode:    0,
				PollOptionIndex: nil,
				EventType:       data.EventType,
				ErrorMessage:    "invalid event payload",
			}
			return errors.New("unsupported event type")
		}

		switch event.Type {
		case constants.Created.String():
			// Process the react story request here
			// You can implement the logic to update the story's reaction count or perform other actions based on the request data
			for i := 0; i < workerCount; i++ {
				c.pool.Run(ctx, func() {
					switch i {
					case 0:
						dataStory, err := c.storyRepo.GetStoryByID(ctx, data.StoryId)
						if err != nil {
							// Handle error
							resultChan <- res.FailedConsumerReactStoryResponse{
								StoryId:         data.StoryId,
								UserId:          data.UserId,
								Avatar:          data.Avatar,
								Name:            data.Name,
								ViewedAt:        data.ViewedAt,
								InteractionType: data.InteractionType,
								ReactionCode:    data.ReactionCode,
								PollOptionIndex: data.PollOptionIndex,
								EventType:       data.EventType,
								ErrorMessage:    err.Error(),
							}
							return
						}
						utils.CreateReactStoryRequestToStoryReactionStats(dataStory, data.ReactionCode)
						err = c.storyRepo.UpdateStory(ctx, dataStory)
						if err != nil {
							// Handle error
							resultChan <- res.FailedConsumerReactStoryResponse{
								StoryId:         data.StoryId,
								UserId:          data.UserId,
								Avatar:          data.Avatar,
								Name:            data.Name,
								ViewedAt:        data.ViewedAt,
								InteractionType: data.InteractionType,
								ReactionCode:    data.ReactionCode,
								PollOptionIndex: data.PollOptionIndex,
								EventType:       data.EventType,
								ErrorMessage:    err.Error(),
							}
							return
						}
						resultChan <- res.FailedConsumerReactStoryResponse{
							StoryId:         data.StoryId,
							UserId:          data.UserId,
							Avatar:          data.Avatar,
							Name:            data.Name,
							ViewedAt:        data.ViewedAt,
							InteractionType: data.InteractionType,
							ReactionCode:    data.ReactionCode,
							PollOptionIndex: data.PollOptionIndex,
							EventType:       data.EventType,
							ErrorMessage:    "",
						}
					case 1:
						dataStoryView, err := c.storyview.GetStoryViewsByStoryIDAndUserID(ctx, data.StoryId, data.UserId)
						if err != nil {
							// Handle error
							resultChan <- res.FailedConsumerReactStoryResponse{
								StoryId:         data.StoryId,
								UserId:          data.UserId,
								Avatar:          data.Avatar,
								Name:            data.Name,
								ViewedAt:        data.ViewedAt,
								InteractionType: data.InteractionType,
								ReactionCode:    data.ReactionCode,
								PollOptionIndex: data.PollOptionIndex,
								EventType:       data.EventType,
								ErrorMessage:    err.Error(),
							}
							return
						}
						if dataStoryView != nil {
							dataStoryView.ReactionCode = data.ReactionCode
							dataStoryView.InteractionType = data.InteractionType
							dataStoryView.PollOptionIndex = data.PollOptionIndex
							err = c.storyview.UpdateStoryView(ctx, dataStoryView)
							if err != nil {
								// Handle error
								resultChan <- res.FailedConsumerReactStoryResponse{
									StoryId:         data.StoryId,
									UserId:          data.UserId,
									Avatar:          data.Avatar,
									Name:            data.Name,
									ViewedAt:        data.ViewedAt,
									InteractionType: data.InteractionType,
									ReactionCode:    data.ReactionCode,
									PollOptionIndex: data.PollOptionIndex,
									EventType:       data.EventType,
									ErrorMessage:    err.Error(),
								}
								return
							}
							resultChan <- res.FailedConsumerReactStoryResponse{
								StoryId:         data.StoryId,
								UserId:          data.UserId,
								Avatar:          data.Avatar,
								Name:            data.Name,
								ViewedAt:        data.ViewedAt,
								InteractionType: data.InteractionType,
								ReactionCode:    data.ReactionCode,
								PollOptionIndex: data.PollOptionIndex,
								EventType:       data.EventType,
								ErrorMessage:    "",
							}

						}
					}
				})
			}
		}
		return nil
	})
	for i := 0; i < workerCount; i++ {
		go func() {
			for res := range resultChan {
				if res.ErrorMessage != "" {
					err := c.eventbus.Publish(ctx, constants.TopicReactStory.String(), res.StoryId, constants.Failed.String(), res)
					if err != nil {
						// Handle publish error
						log.Printf("Error publishing failed react story event: %v", err)
					}
				}
			}
		}()
	}

	if err != nil {
		// Handle subscription error
		log.Printf("Error subscribing to topic: %v", err)
		return
	}

}
func (c *ConsumerReact) ConsumerFailedReactStory(ctx context.Context) {
	// Implement the logic for consuming failed react story events here
	err := c.eventbus.Subscribe(ctx, constants.TopicReactStory.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*res.FailedConsumerReactStoryResponse)
		if !ok {
			// Handle type assertion error
			return errors.New("invalid event payload")
		}
		if event.Type != constants.Failed.String() {
			return errors.New("unsupported event type")
		}
		// Process the failed react story response here
		log.Printf("Failed to process react story event for StoryID: %s, UserID: %s, Error: %s", data.StoryId, data.UserId, data.ErrorMessage)
		return nil
	})
	if err != nil {
		// Handle subscription error
		log.Printf("Error subscribing to topic: %v", err)
		return
	}
}
