package IKafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type IBlockMessage interface {
	PublishBlockCreateMessage(ctx context.Context, payload socialEvent.BlockCreatePayload) error
	PublishBlockUpdateMessage(ctx context.Context, payload socialEvent.BlockUpdatePayload) error
	PublishBlockDeleteMessage(ctx context.Context, payload socialEvent.BlockDeletePayload) error
}
