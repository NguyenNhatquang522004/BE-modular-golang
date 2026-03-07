package IConsumer

import "context"

type IConsumerPost interface {
	ConsumerPost(ctx context.Context) error
	ConsumerFailedPost(ctx context.Context) error
}
