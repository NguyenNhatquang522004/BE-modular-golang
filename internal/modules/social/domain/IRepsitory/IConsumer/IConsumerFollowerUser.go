package IConsumer

import "context"

type IConsumerFollowerUser interface {
	ConsumerFollowUser(ctx context.Context) error
	ConsumerFailedFollowUser(ctx context.Context) error
}
