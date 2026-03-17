package IConsumer

import "context"

type IConsumerGraphPage interface {
	ConsumerGraphPage(ctx context.Context) error
	ConsumerFailedGraphPage(ctx context.Context) error
}
