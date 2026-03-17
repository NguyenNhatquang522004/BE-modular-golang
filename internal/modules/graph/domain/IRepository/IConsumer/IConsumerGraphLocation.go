package IConsumer

import "context"

type IConsumerGraphLocation interface {
	ConsumerGraphLocation(ctx context.Context) error
	ConsumerFailedGraphLocation(ctx context.Context) error
}
