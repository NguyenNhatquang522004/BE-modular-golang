package IConsumer

import "context"

type IConsumerPostStats interface {
	ConsumerPostStats(ctx context.Context) error
	ConsumerFailedPostStats(ctx context.Context) error
}
