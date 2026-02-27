package IConsumer

import "context"

type IConsumerView interface {
	ConsumerViewCountStory(ctx context.Context)
	ConsumerFailedViewCountStory(ctx context.Context)
}
