package kafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type FollowMessage struct {
	eventBus events.EventBus
}

func NewFollowMessage(eventBus events.EventBus) *FollowMessage {
	return &FollowMessage{
		eventBus: eventBus,
	}
}
func (r *FollowMessage) PublishFollowCreateMessage(ctx context.Context, payload socialEvent.FollowCreatePayload) error {
	err := r.eventBus.Publish(ctx, constants.TopicFollow.String(), payload.Follower_UserID, constants.Created.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
func (r *FollowMessage) PublishFollowDeleteMessage(ctx context.Context, payload socialEvent.FollowDeletePayload) error {
	err := r.eventBus.Publish(ctx, constants.TopicFollow.String(), payload.Follower_UserID, constants.Deleted.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
func (r *FollowMessage) PublishFollowUpdateMessage(ctx context.Context, payload socialEvent.FollowUpdatePayload) error {
	err := r.eventBus.Publish(ctx, constants.TopicFollow.String(), payload.Follower_UserID, constants.Updated.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
