package IConsumer

import "context"

type IConsumerGraphGroup interface {
	ConsumerGraphGroup(ctx context.Context) error
	ConsumerFailedGraphGroup(ctx context.Context) error
}
