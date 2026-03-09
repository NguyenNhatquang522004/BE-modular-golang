package IConsumer

import "context"

type IConsumerSharePost interface {
	ConsumerSharePost(ctx context.Context) error
	ConsumerFailedSharePost(ctx context.Context) error
}
