package IRepositoryConsumer

import "context"

type IConsumerGroupFile interface {
	ConsumerGroupFile(ctx context.Context) error
	ConsumerFailedGroupFile(ctx context.Context) error
}
