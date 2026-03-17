package IConsumer

import "context"

type IConsumerGraphInteractions interface {
	ConsumerGraphInteractions(ctx context.Context) error
	ConsumerFailedGraphInteractions(ctx context.Context) error
}
