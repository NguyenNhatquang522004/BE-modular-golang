package IConsumer

import "context"

type IConsumerReelStats interface {
	ConsumerReelStats(ctx context.Context) error
	ConsumerFailedReelStats(ctx context.Context) error
}
