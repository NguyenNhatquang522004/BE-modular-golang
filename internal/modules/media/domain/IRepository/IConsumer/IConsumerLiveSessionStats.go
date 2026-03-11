package IConsumer

import "context"

type IConsumerLiveSessionStats interface {
	ConsumerLiveSessionStats(ctx context.Context) error
	ConsumerFailedLiveSessionStats(ctx context.Context) error
}
