package IProducer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type IFollowMessage interface {

	PublishFollowDeleteMessage(ctx context.Context, payload ...*socialEvent.FollowDeletePayload) error
	PublishFollowUpdateMessage(ctx context.Context, payload ...*socialEvent.FollowUpdatePayload) error
}
