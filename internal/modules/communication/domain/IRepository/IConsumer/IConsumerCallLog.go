package IConsumer

import "context"

type IConsumerCallLog interface {
	ConsumerCallLogEvents(ctx context.Context) error
	ConsumerFailedCallLogEvents(ctx context.Context) error
}
