package IRepositoryConsumer

import "context"

type IConsumerGroupQA interface {
	ConsumerGroupQA(ctx context.Context) error
	ConsumerFailedGroupQA(ctx context.Context) error
}
