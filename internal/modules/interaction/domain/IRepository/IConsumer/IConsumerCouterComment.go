package IConsumer

import "context"

type IConsumerCounterComment interface {
	ConsumeCounterComment(ctx context.Context) error
	ConsumerFailedCounterComment(ctx context.Context) error
}
