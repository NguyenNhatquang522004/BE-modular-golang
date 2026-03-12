package IConsumer

import "context"

type IConsumerStatsMessage interface {
	ConsumerStatsMessage(ctx context.Context) error
	ConsumerFailedStatsMessage(ctx context.Context) error
}
