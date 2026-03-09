package IConsumer

import "context"

type IConsumerCommentStats interface {
	ConsumerCommentStats(ctx context.Context) error
	ConsumerFailedCommentStats(ctx context.Context) error
}
