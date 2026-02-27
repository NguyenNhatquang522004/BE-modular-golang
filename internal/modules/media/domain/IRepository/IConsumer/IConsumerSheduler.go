package IConsumer

import "context"

type IConsumerSheduler interface {
	ConsumerDeleteStoryExpiresAt(ctx context.Context)
}
