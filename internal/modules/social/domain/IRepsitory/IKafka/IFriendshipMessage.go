package IKafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type IFriendshipMessage interface {
	PublishFriendshipCreateMessage(ctx context.Context, payload socialEvent.FriendshipsCreatePayload) error
	PublishFriendshipDeleteMessage(ctx context.Context, payload socialEvent.FriendshipsDeletePayload) error
	PublishFriendshipUpdateMessage(ctx context.Context, payload socialEvent.FriendshipsUpdatePayload) error
}
