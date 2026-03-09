package IConsumer

import "context"

type IConsumerEntityReaction interface {
	ConsumerEntityReaction(ctx context.Context) error
	ConsumerFailedEntityReaction(ctx context.Context) error
}
