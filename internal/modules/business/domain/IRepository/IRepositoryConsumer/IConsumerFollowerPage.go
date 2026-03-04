package IRepositoryConsumer

import "context"

type IConsumerFollowerPage interface {
	ConsumerFollowerPage(ctx context.Context) error
	ConsumerFailedFollowerPage(ctx context.Context) error
}
