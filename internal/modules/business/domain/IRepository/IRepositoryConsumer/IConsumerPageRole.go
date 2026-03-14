package IRepositoryConsumer

import "context"

type IConsumerPageRole interface {
	ConsumerPageRole(ctx context.Context) error
	ConsumerFailedPageRole(ctx context.Context) error
}
