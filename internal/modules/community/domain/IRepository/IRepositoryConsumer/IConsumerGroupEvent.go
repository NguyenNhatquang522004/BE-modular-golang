package IRepositoryConsumer

import "context"

type IConsumerGroupEvent interface {
	ConsumerGroupEvent(ctx context.Context) error
	ConsumerFailedGroupEvent(ctx context.Context) error
}
