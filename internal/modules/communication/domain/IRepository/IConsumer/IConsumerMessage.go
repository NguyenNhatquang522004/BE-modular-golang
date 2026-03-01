package IConsumer

import "context"

type IConsumerMessage interface {
	ConsumerMessage(ctx context.Context) error
	FailedMessage(ctx context.Context) error
	ConsumerStateMessage(ctx context.Context) error
	FailedStateMessage(ctx context.Context) error
	ConsumerReactMessage(ctx context.Context) error
	FailedReactMessage(ctx context.Context) error
}
