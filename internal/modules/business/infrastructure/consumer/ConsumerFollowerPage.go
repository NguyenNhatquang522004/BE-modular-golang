package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ConsumerFollowerPage struct {
	events           events.EventBus
	pageRepo         IRepositoryMongodb.IPagesRepository
	pageFollowerRepo IRepositoryMongodb.IPageFollowersRepository
	pool             IRepositoryShare.IWorkerPool
}

func NewConsumerFollowerPage(events events.EventBus, pageRepo IRepositoryMongodb.IPagesRepository, pageFollowerRepo IRepositoryMongodb.IPageFollowersRepository, pool IRepositoryShare.IWorkerPool) *ConsumerFollowerPage {
	return &ConsumerFollowerPage{
		events:           events,
		pageRepo:         pageRepo,
		pageFollowerRepo: pageFollowerRepo,
		pool:             pool,
	}
}

func (c *ConsumerFollowerPage) ConsumerFollowerPage(ctx context.Context) error {
	err := c.events.Subscribe(ctx, constants.TopicFollowerPage.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*businessEvent.FollowerPagePayload)
		if !ok {
			return errors.New("invalid event payload type")
		}
		if data == nil {
			return errors.New("event payload is nil")
		}
		switch event.Type {
		case constants.Created.String():
			datapage, err := c.pageRepo.GetPageByID(ctx, data.PageID)
			if err != nil {
				return errors.New("page not found: " + data.PageID)
			}
			if datapage == nil {
				return errors.New("page not found: " + data.PageID)
			}
			entity := &entity.PageFollower{
				ID:     primitive.NewObjectID(),
				PageID: datapage.ID,
				UserID: data.UserID,
				Settings: &entity.FollowerSettings{
					NotificationLevel: data.Settings.NotificationLevel,
					IsFavorite:        data.Settings.IsFavorite,
				},
				CreatedAt: time.Now(),
			}
			err = c.pageFollowerRepo.CreateFollower(ctx, entity)
			if err != nil {
				return errors.New("failed to create page follower: " + err.Error())
			}
			payload := &businessEvent.StatsPagePayload{
				UserID:         data.UserID,
				PageID:         data.PageID,
				FollowersCount: 1,
				LikesCount:     1,
				RatingScore:    0,
				ReviewCount:    0,
				CreatedAt:      time.Now(),
			}
			err = c.events.Publish(ctx, constants.TopicStatsPage.String(), data.UserID, constants.Created.String(), payload)
			if err != nil {
				return errors.New("failed to publish stats page event: " + err.Error())
			}
			// Xử lý event ở đây
			return nil
		case constants.Updated.String():
			datapage, err := c.pageRepo.GetPageByID(ctx, data.PageID)
			if err != nil {
				return errors.New("page not found: " + data.PageID)
			}
			if datapage == nil {
				return errors.New("page not found: " + data.PageID)
			}
			datafollower, err := c.pageFollowerRepo.GetFollowerByPageIDAndUserID(ctx, data.PageID, data.UserID)
			if err != nil {
				return errors.New("failed to get page follower: " + err.Error())
			}
			if datafollower == nil {
				return errors.New("page follower not found for pageID: " + data.PageID + " and userID: " + data.UserID)
			}
			req := &req.PageFollowerReq{
				ID:     datafollower.ID.Hex(),
				PageID: data.PageID,
				UserID: data.UserID,
				Settings: &req.FollowerSettingsReq{
					NotificationLevel: data.Settings.NotificationLevel,
					IsFavorite:        data.Settings.IsFavorite,
				},
			}

			mapper.UpdateToEntityPageFollower(req, datafollower)
			err = c.pageFollowerRepo.UpdateFollower(ctx, datafollower)
			if err != nil {
				return errors.New(" failed to update page follower: " + err.Error())
			}
			payload := &businessEvent.StatsPagePayload{
				UserID:         data.UserID,
				PageID:         data.PageID,
				FollowersCount: 0,
				LikesCount:     0,
				RatingScore:    0,
				ReviewCount:    0,
				CreatedAt:      time.Now(),
			}
			err = c.events.Publish(ctx, constants.TopicStatsPage.String(), data.UserID, constants.Created.String(), payload)
			if err != nil {
				return errors.New("failed to publish stats page event: " + err.Error())
			}
			// Xử lý event ở đây
			return nil
		case constants.Deleted.String():
			datapage, err := c.pageRepo.GetPageByID(ctx, data.PageID)
			if err != nil {
				return errors.New("page not found: " + data.PageID)
			}
			if datapage == nil {
				return errors.New("page not found: " + data.PageID)
			}
			datafollower, err := c.pageFollowerRepo.GetFollowerByPageIDAndUserID(ctx, data.PageID, data.UserID)
			if err != nil {
				return errors.New("failed to get page follower: " + err.Error())
			}
			if datafollower == nil {
				return errors.New("page follower not found for pageID: " + data.PageID + " and userID: " + data.UserID)
			}
			err = c.pageFollowerRepo.DeleteFollower(ctx, data.PageID, data.UserID)
			if err != nil {
				return errors.New("failed to delete page follower: " + err.Error())
			}
			payload := &businessEvent.StatsPagePayload{
				UserID:         data.UserID,
				PageID:         data.PageID,
				FollowersCount: -1,
				LikesCount:     -1,
				RatingScore:    0,
				ReviewCount:    0,
				CreatedAt:      time.Now(),
			}
			err = c.events.Publish(ctx, constants.TopicStatsPage.String(), data.UserID, constants.Updated.String(), payload)
			if err != nil {
				return errors.New("failed to publish stats page event: " + err.Error())
			}

			// Xử lý event ở đây
		default:
			return errors.New("unsupported event type: " + cases.Title(language.English).String(event.Type))
		}
		// Xử lý logic khi nhận được sự kiện FollowerPagePayload
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerFollowerPage) ConsumerFailedFollowerPage(ctx context.Context) error
