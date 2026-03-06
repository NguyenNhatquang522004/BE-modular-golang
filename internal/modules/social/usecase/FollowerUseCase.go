package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
)

type FollowUseCase struct {
	events events.EventBus
}

func NewFollowUseCase(events events.EventBus) *FollowUseCase {
	return &FollowUseCase{
		events: events,
	}
}

func (c *FollowUseCase) Execute(ctx context.Context, req *req.FollowerRequest) error {
	payload := &socialEvent.FollowerUserPayload{
		ID:             req.ID,
		FollowerUserID: req.Follower_UserID,
		FollowedUserID: req.Followed_UserID,
		IsMuted:        req.IsMuted,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      req.UpdatedAt,
		EventType:      req.EventType,
	}
	err := c.events.Publish(ctx, constants.TopicFollowUser.String(), req.Follower_UserID, req.EventType.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
