package IConsumer

import "context"

type IConsumerCommentLive interface {
	ConsumerStartLiveStream(ctx context.Context) error
	ConsumerFailedStartLiveStream(ctx context.Context) error
}
