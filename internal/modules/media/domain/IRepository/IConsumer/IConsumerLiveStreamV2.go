package IConsumer

import "context"

type IConsumerLiveStreamV2 interface {
	ConsumeStartLiveStreamV2(ctx context.Context) error
	ConsumerFailedTopicLiveStreamV2(ctx context.Context) error
}
