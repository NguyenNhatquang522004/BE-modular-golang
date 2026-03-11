package IConsumer

import "context"

type ICommentLiveStream interface {
	ConsumerLiveComment(ctx context.Context) error
	ConsumerFailedLiveComment(ctx context.Context) error
}
