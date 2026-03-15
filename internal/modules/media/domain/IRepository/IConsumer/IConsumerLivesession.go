package IConsumer

import "context"

type IConsumerLiveSession interface {
	ConsumerTopicLiveSession(ctx context.Context) error
	ConsumerFailedTopicLiveSession(ctx context.Context) error
}
