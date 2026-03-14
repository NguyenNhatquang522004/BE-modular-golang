package IRepositoryConsumer

import "context"

type IConsumerPage interface {
	ConsumerPage(ctx context.Context) error
	ConsumerFailedPage(ctx context.Context) error
}
