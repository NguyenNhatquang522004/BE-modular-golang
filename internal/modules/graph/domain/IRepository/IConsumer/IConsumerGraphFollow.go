package IConsumer

import "context"

type IConsumerGraphFollow interface {
	ConsumerGraphFollow(ctx context.Context) error
	ConsumerFailedGraphFollow(ctx context.Context) error
}
