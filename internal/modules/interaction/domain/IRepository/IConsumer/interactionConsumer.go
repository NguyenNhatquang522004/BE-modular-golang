package IConsumer

import "context"

type InteractionConsumer interface {
	ConsumerReactionComment(ctx context.Context) error
	ConsumerFailedReactionComment(ctx context.Context) error
	ConsumerEntityReaction(ctx context.Context) error
	ConsumerFailedEntityReaction(ctx context.Context) error
}
