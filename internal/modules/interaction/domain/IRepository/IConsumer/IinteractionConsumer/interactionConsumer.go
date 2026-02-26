package IinteractionConsumer

import "context"

type InteractionConsumer interface {
	ConsumerReactionComment(ctx context.Context) error
}
