package IConsumer

import "context"

type IConsumerGraphBlock interface {
	ConsumerGraphBlockEvent(ctx context.Context) error
	ConsumerFailedGraphBlockEvent(ctx context.Context) error
}

