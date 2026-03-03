package IConsumer

import "context"

type IConsumerDeleteMediaRelationTarget interface {
	ConsumerDeleteMediaRelationTarget(ctx context.Context) error
	ConsumerFailedDeleteMediaRelationTarget(ctx context.Context) error
}
