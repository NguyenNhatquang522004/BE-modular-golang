package IRepositoryConsumer

import "context"

type IConsumerStats interface {
	ConsumerGroupStats(ctx context.Context) error
	ConsumerFailedGroupStats(ctx context.Context) error


}
