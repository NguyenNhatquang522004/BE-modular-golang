package kafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type FriendshipMessage struct {
	eventBus events.EventBus
}

func NewFriendshipMessage(eventBus events.EventBus) *FriendshipMessage {
	return &FriendshipMessage{
		eventBus: eventBus,
	}
}
func (r *FriendshipMessage) PublishFriendshipCreateMessage(ctx context.Context, payload socialEvent.FriendshipsCreatePayload) error {
	err := r.eventBus.Publish(ctx, constants.TopicFriendship.String(), payload.Requester_ID, constants.Created.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
func (r *FriendshipMessage) PublishFriendshipDeleteMessage(ctx context.Context, payload socialEvent.FriendshipsDeletePayload) error {
	err := r.eventBus.Publish(ctx, constants.TopicFriendship.String(), payload.Requester_ID, constants.Deleted.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
func (r *FriendshipMessage) PublishFriendshipUpdateMessage(ctx context.Context, payload socialEvent.FriendshipsUpdatePayload) error {
	err := r.eventBus.Publish(ctx, constants.TopicFriendship.String(), payload.Requester_ID, constants.Updated.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
