package IConsumer

import "context"

type IConsumerGraphPost interface {
	ConsumerGraphPost(ctx context.Context) error
	ConsumerFailedGraphPost(ctx context.Context) error
}
