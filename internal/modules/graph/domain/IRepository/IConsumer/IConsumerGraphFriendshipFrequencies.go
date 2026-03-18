package IConsumer

import "context"

type IConsumerGraphFriendshipFrequencies interface {
	ConsumerGraphFriendshipFrequencies(ctx context.Context) error
	ConsumerFailedGraphFriendshipFrequencies(ctx context.Context) error
}
