package IConsumer

import "context"

type IConsumerMessage interface {
	ConsumerMessage(ctx context.Context) error
	FailedMessage(ctx context.Context) error
	ConsumerStatsMessage(ctx context.Context) error
	FailedStatsMessage(ctx context.Context) error
}
