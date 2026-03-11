package IConsumer

import "context"

type IConsumerComment interface {
	ConsumerComment(ctx context.Context) error
	ConsumerFailedComment(ctx context.Context) error
}
