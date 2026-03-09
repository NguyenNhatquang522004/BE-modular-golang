package consumer

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/gocql/gocql"
)

type ConsumerStats struct {
	events   events.EventBus
	pageRepo IRepositoryMongodb.IPagesRepository
}

func NewConsumerStats(events events.EventBus, pageRepo IRepositoryMongodb.IPagesRepository) *ConsumerStats {
	return &ConsumerStats{
		events:   events,
		pageRepo: pageRepo,
	}
}

func (c *ConsumerStats) ConsumerStatsPage(ctx context.Context) error {
	err := c.events.Subscribe(ctx, constants.TopicStatsPage.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*businessEvent.StatsPagePayload)

		if !ok {
			return errors.New("Invalid payload type for StatsPage event")
		}
		switch event.Type {
		case constants.Created.String():
			datapage, err := c.pageRepo.GetPageByID(ctx, data.PageID)
			if err != nil {
				return err
			}
			if datapage == nil {
				return errors.New("Page not found for StatsPage event")
			}
			// Cập nhật thống kê cho trang ở đây, ví dụ:
			datapage.Stats.FollowersCount = datapage.Stats.FollowersCount + data.FollowersCount
			datapage.Stats.LikesCount = datapage.Stats.LikesCount + data.LikesCount
			datapage.Stats.RatingScore = datapage.Stats.RatingScore + data.RatingScore
			datapage.Stats.ReviewCount = datapage.Stats.ReviewCount + data.ReviewCount
			err = c.pageRepo.UpdatePage(ctx, datapage)
			if err != nil {
				return err
			}
			convertuserIDcql, err := gocql.ParseUUID(data.UserID)
			if err != nil {
				return err
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.PageID,
				UserID:       convertuserIDcql,
				TargetType:   sharedEnums.ReactionTargetFollowPage,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    data.CreatedAt,
				Type:         constants.Created,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, constants.Created.String(), payload)
			if err != nil {
				return err
			}
			// Xử lý event ở đây
			return nil
		case constants.Updated.String():
			datapage, err := c.pageRepo.GetPageByID(ctx, data.PageID)
			if err != nil {
				return err
			}
			if datapage == nil {
				return errors.New("Page not found for StatsPage event")
			}
			// Cập nhật thống kê cho trang ở đây, ví dụ:
			datapage.Stats.FollowersCount = datapage.Stats.FollowersCount + data.FollowersCount
			datapage.Stats.LikesCount = datapage.Stats.LikesCount + data.LikesCount
			datapage.Stats.RatingScore = datapage.Stats.RatingScore + data.RatingScore
			datapage.Stats.ReviewCount = datapage.Stats.ReviewCount + data.ReviewCount
			err = c.pageRepo.UpdatePage(ctx, datapage)
			if err != nil {
				return err
			}
			convertuserIDcql, err := gocql.ParseUUID(data.UserID)
			if err != nil {
				return err
			}
			payload := &interactionEvent.EntityReactionPayload{
				TargetID:     data.PageID,
				UserID:       convertuserIDcql,
				TargetType:   sharedEnums.ReactionTargetUnFollowPage,
				ReactionCode: sharedEnums.ReactionUnknown,
				CreatedAt:    data.CreatedAt,
				Type:         constants.Created,
			}
			err = c.events.Publish(ctx, constants.TopicEntityReaction.String(), data.UserID, constants.Created.String(), payload)
			if err != nil {
				return err
			}
			// Xử lý event ở đây
			return nil

		case constants.Deleted.String():
		default:
			return errors.New("Unsupported event type for StatsPage event")
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerStats) ConsumerFailedStatsPage(ctx context.Context) error {
	return nil
}
