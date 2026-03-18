package IConsumer

import "context"

type IConsumerGraphLikePage interface {
	ConsumerGraphLikePageEvent(ctx context.Context) error
	ConsumerFailedGraphLikePageEvent(ctx context.Context) error
}
