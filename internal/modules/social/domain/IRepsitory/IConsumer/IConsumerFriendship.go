package IConsumer

import "context"

type IConsumerFriendship interface {
	ConsumerFriendUser(ctx context.Context) error
	ConsumerFailedFriendUser(ctx context.Context) error
}
