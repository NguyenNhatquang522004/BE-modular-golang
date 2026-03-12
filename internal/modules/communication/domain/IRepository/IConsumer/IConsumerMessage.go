package IConsumer

import "context"

type IConsumerMessage interface {
	ConsumerMessage(ctx context.Context) error
	ConsumerFailedMessage(ctx context.Context) error
}
