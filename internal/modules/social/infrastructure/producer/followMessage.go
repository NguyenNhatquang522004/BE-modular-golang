package producer

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
func (r *FollowMessage) PublishFollowCreateMessage(ctx context.Context, payload ...*socialEvent.FollowCreatePayload) error {
	for _, p := range payload {
		err := r.eventBus.Publish(ctx, constants.TopicFollow.String(), p.Follower_UserID, constants.Created.String(), p)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *FollowMessage) PublishFollowDeleteMessage(ctx context.Context, payload ...*socialEvent.FollowDeletePayload) error {
	for _, p := range payload {
		err := r.eventBus.Publish(ctx, constants.TopicFollow.String(), p.Follower_UserID, constants.Deleted.String(), p)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *FollowMessage) PublishFollowUpdateMessage(ctx context.Context, payload ...*socialEvent.FollowUpdatePayload) error {
	for _, p := range payload {
		err := r.eventBus.Publish(ctx, constants.TopicFollow.String(), p.Follower_UserID, constants.Updated.String(), p)
		if err != nil {
			return err
		}
	}
	return nil
}
