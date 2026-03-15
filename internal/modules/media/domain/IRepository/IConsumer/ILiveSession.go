package IConsumer

import "context"

type IConsumerTopicLiveSession interface {
	ConsumerTopicLiveSession(ctx context.Context) error
	ConsumerFailedTopicLiveSession(ctx context.Context) error
}
