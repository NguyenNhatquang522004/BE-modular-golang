package IConsumer

import "context"

type IConsumerGraphInteractionsRecent interface {
	ConsumerGraphInteractionsRecent(ctx context.Context) error
	ConsumerFailedGraphInteractionsRecent(ctx context.Context) error
}
