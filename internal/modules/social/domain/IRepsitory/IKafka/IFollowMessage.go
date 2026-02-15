package IKafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type IFollowMessage interface {
	PublishFollowCreateMessage(ctx context.Context, payload socialEvent.FollowCreatePayload) error
	PublishFollowDeleteMessage(ctx context.Context, payload socialEvent.FollowDeletePayload) error
	PublishFollowUpdateMessage(ctx context.Context, payload socialEvent.FollowUpdatePayload) error
}
