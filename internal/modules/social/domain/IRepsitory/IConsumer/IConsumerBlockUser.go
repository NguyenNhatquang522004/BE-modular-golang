package IConsumer

import "context"

type IConsumerBlockUser interface {
	ConsumerBlockUser(ctx context.Context) error
	ConsumerFailedBlockUser(ctx context.Context) error
}
