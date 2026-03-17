package IConsumer

import "context"

type IConsumerGraphSharedPost interface {
	ConsumerGraphSharedPost(ctx context.Context) error
	ConsumerFailedGraphSharedPost(ctx context.Context) error
}
