package IConsumer

import "context"

type IConsumerDeleteRelationTarget interface {
	ConsumerDeleteRelationTarget(ctx context.Context) error
	FailedDeleteRelationTarget(ctx context.Context) error
}
