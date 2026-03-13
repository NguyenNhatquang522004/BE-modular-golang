package IRepositoryConsumer

import "context"

type IConsumerGroup interface {
	ConsumerGroup(ctx context.Context) error
	ConsumerFailedGroup(ctx context.Context) error
}
