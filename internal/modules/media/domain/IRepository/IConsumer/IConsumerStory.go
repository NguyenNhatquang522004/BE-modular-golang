package IConsumer

import "context"

type IConsumerStory interface {
	ConsumerStory(ctx context.Context) error
	ConsumerFailedStory(ctx context.Context) error
}
