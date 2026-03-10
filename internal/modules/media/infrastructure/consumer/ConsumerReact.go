package consumer

import (
	"context"
	"errors"
	"log"
	"math"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/utils"
	"github.com/gocql/gocql"
)

type ConsumerReact struct {
	// Define any dependencies needed for consuming react events here
	ablumRepo       IRepositoryMongodb.IAlbumsRepository
	mediaAssetsRepo IRepositoryMongodb.IMediaAssetsRepository
	storyRepo       IRepositoryMongodb.IStoryRepository
	storyview       IRepositoryCassandra.IStoryViewRepository
	reelRepo        IRepositoryMongodb.IReelRepository
	livesessionRepo IRepositoryMongodb.ILiveSessionRepository
	liveCommentRepo IRepositoryCassandra.ILiveCommentsRepository
	pool            IRepositoryShare.IWorkerPool
	eventbus        events.EventBus
}

func NewConsumerReact(ablumRepo IRepositoryMongodb.IAlbumsRepository,
	mediaAssetsRepo IRepositoryMongodb.IMediaAssetsRepository,
	storyRepo IRepositoryMongodb.IStoryRepository,
	storyview IRepositoryCassandra.IStoryViewRepository,
	livesessionRepo IRepositoryMongodb.ILiveSessionRepository,
	liveCommentRepo IRepositoryCassandra.ILiveCommentsRepository,
	pool IRepositoryShare.IWorkerPool,
	eventbus events.EventBus,
	reelRepo IRepositoryMongodb.IReelRepository) *ConsumerReact {
	return &ConsumerReact{
		ablumRepo:       ablumRepo,
		mediaAssetsRepo: mediaAssetsRepo,
		storyRepo:       storyRepo,
		storyview:       storyview,
		reelRepo:        reelRepo,
		livesessionRepo: livesessionRepo,
		liveCommentRepo: liveCommentRepo,
		pool:            pool,
		eventbus:        eventbus,
	}
}
func (c *ConsumerReact) CosumerReactStory(ctx context.Context) {

	workerCount := 2 // Adjust the number of workers as needed
	resultChan := make(chan res.FailedConsumerReactStoryResponse, workerCount)
	// Implement the logic for consuming react story events here
	err := c.eventbus.Subscribe(ctx, constants.TopicReactStory.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.ReactStoryPayload)
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
						storyid, err := gocql.ParseUUID(data.StoryId)
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
						viewerid, err := gocql.ParseUUID(data.UserId)
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
						newStoryView := &entity.StoryView{
							StoryID:         storyid,
							ViewerID:        viewerid,
							ViewerName:      data.Name,
							ViewerAvatarURL: data.Avatar,
							ViewedAt:        data.ViewedAt,
							InteractionType: data.InteractionType,
							ReactionCode:    data.ReactionCode,
							PollOptionIndex: data.PollOptionIndex,
							Content:         data.Content,
						}
						err = c.storyview.CreateStoryView(ctx, newStoryView)
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
					default:
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
							ErrorMessage:    "invalid worker index",
						}
						return
						// Handle unknown worker index
						// You can choose to log this or return an error as needed
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
func (c *ConsumerReact) ConsumerReactReel(ctx context.Context) {
	// Implement the logic for consuming react reel events here
	err := c.eventbus.Subscribe(ctx, constants.TopicReactReel.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.ReactReelPayload)
		if !ok {
			// Handle type assertion error
			return errors.New("invalid event payload")
		}
		switch event.Type {
		case constants.Created.String():
			datareel, err := c.reelRepo.GetReelByID(ctx, data.ReelID)
			if err != nil {
				// Handle error
				return nil
			}
			datareel.Stats.Total = data.Total + 1
			utils.CreateReactReelRequestToReelReactionStats(datareel, data.ReactionCode)
			err = c.reelRepo.UpdateReel(ctx, datareel)
			if err != nil {
				// Handle error
				return nil
			}
		case constants.Deleted.String():
		}
		return nil
	})
	if err != nil {
		// Handle subscription error
		log.Printf("Error subscribing to topic: %v", err)
		return
	}
}
func (c *ConsumerReact) ConsumerFailedReactReel(ctx context.Context)

func (c *ConsumerReact) ConsumerFailedCounterReel(ctx context.Context) {

}
func (c *ConsumerReact) ConsumerCounterReel(ctx context.Context) {
	err := c.eventbus.Subscribe(ctx, constants.TopicCounterReel.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.ReactCounterReelPayload)
		if !ok {
			// Handle type assertion error
			return errors.New("invalid event payload")
		}
		switch event.Type {
		case constants.Created.String():
			datareel, err := c.reelRepo.GetReelByID(ctx, data.ReelID)
			if err != nil {
				// Handle error
				return nil
			}
			datareel.Stats.Comments = data.Comments + datareel.Stats.Comments
			datareel.Stats.Saves = data.Saves + datareel.Stats.Saves
			datareel.Stats.Shares = data.Shares + datareel.Stats.Shares
			err = c.reelRepo.UpdateReel(ctx, datareel)
			if err != nil {
				// Handle error
				return nil
			}
			return nil
		}
		return nil
	})
	if err != nil {
		log.Printf("Error subscribing to topic: %v", err)
		return
	}
	return
}
func (c *ConsumerReact) ConsumerReactLive(ctx context.Context) {
	err := c.eventbus.Subscribe(ctx, constants.TopicReactLive.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.ReactLiveStreamPayload)
		if !ok {
			// Handle type assertion error
			return errors.New("invalid event payload")
		}
		switch event.Type {
		case constants.Created.String():
			// Process the react live stream request here
			// You can implement the logic to update the live stream's reaction count or perform other actions based on the request data
			dataLive, err := c.livesessionRepo.GetLiveSessionByID(ctx, data.LiveSessionID)
			if err != nil {
				// Handle error
				return nil
			}
			utils.CreateReactLiveSessionRequestToLiveSessionReactionStats(dataLive, data.ReactionCode)
			err = c.livesessionRepo.UpdateLiveSession(ctx, dataLive)
			if err != nil {
				// Handle error
				return nil
			}
			// You can also choose to publish an event or perform other actions as needed
			convertusercql, err := gocql.ParseUUID(data.UserID)
			if err != nil {
				// Handle error
				return nil
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.LiveSessionID,
				UserID:       convertusercql,
				TargetType:   data.TargetType,
				ReactionCode: data.ReactionCode,
				CreatedAt:    data.CreatedAt,
				Type:         constants.Created,
			}
			err = c.eventbus.Publish(ctx, constants.TopicEntityReaction.String(), data.LiveSessionID, constants.Created.String(), payload)
			if err != nil {
				// Handle publish error
				log.Printf("Error publishing entity reaction event: %v", err)
			}
		case constants.Deleted.String():
			// Process the react live stream request here for deleted event type
		default:
			// Handle unknown event type
		}
		return nil
	})
	if err != nil {
		log.Printf("Error subscribing to topic: %v", err)
		return
	}

}
func (c *ConsumerReact) ConsumerFailedReactLive(ctx context.Context)

func (c *ConsumerReact) ConsumerCounterLive(ctx context.Context) {
	err := c.eventbus.Subscribe(ctx, constants.TopicCounterLive.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.CoutnerLiveStreamPayload)
		if !ok {
			// Handle type assertion error
			return errors.New("invalid event payload")
		}
		switch event.Type {
		case constants.Created.String():
			dataLive, err := c.livesessionRepo.GetLiveSessionByID(ctx, data.LiveSessionID)
			if err != nil {
				// Handle error
				return nil
			}
			pek := math.Max(float64(data.Views+dataLive.Stats.PeakViewers), float64(dataLive.Stats.PeakViewers))
			dataLive.Stats.PeakViewers = int(pek)
			dataLive.Stats.TotalComments = data.Comments + dataLive.Stats.TotalComments
			dataLive.Stats.TotalViews = data.Views + dataLive.Stats.TotalViews
			err = c.livesessionRepo.UpdateLiveSession(ctx, dataLive)
			if err != nil {
				// Handle error
				return nil
			}
			return nil
		default:
			// Handle unknown event type
		}

		return nil
	})

	if err != nil {
		log.Printf("Error subscribing to topic: %v", err)
		return
	}
}
func (c *ConsumerReact) ConsumerFailedCounterLive(ctx context.Context)

func (c *ConsumerReact) ConsumerCounterReplyStory(ctx context.Context) {
	err := c.eventbus.Subscribe(ctx, constants.TopicReplyStory.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*mediaEvent.ReplyStoryPayload)
		if !ok {
			// Handle type assertion error
			return errors.New("invalid event payload")
		}
		switch event.Type {
		case constants.Created.String():
			dataStory, err := c.storyRepo.GetStoryByID(ctx, data.StoryID)
			if err != nil {
				// Handle error
				return nil
			}
			dataStory.Stats.ReplyCount = dataStory.Stats.ReplyCount + 1
			err = c.storyRepo.UpdateStory(ctx, dataStory)
			if err != nil {
				// Handle error
				return nil
			}
		case constants.Deleted.String():
		default:
			// Handle unknown event type
		}
		return nil
	})
	if err != nil {
		log.Printf("Error subscribing to topic: %v", err)
		return
	}
}

func (c *ConsumerReact) ConsumerFailedCounterReplyStory(ctx context.Context)
