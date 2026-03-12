package IConsumer

import "context"

type IConsumerReel interface {
	ConsumerReel(ctx context.Context) error
	ConsumerFailedReel(ctx context.Context) error
}
