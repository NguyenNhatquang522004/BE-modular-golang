package IConsumer

import "context"

type IConsumerStoryStats interface {
	ConsumerStoryStats(ctx context.Context) error
	ConsumerFailedStoryStats(ctx context.Context) error
}
