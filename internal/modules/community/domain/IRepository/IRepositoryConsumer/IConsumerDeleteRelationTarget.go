package IRepositoryConsumer

import (
	"context"
)

type IConsumerDeleteRelationTarget interface {
	ConsumerDeleteRelationTarget(ctx context.Context) error
	ConsumerFailedDeleteRelationTarget(ctx context.Context) error
}
