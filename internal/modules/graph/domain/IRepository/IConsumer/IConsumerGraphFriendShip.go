package IConsumer

import "context"

type IConsumerGraphFriendShip interface {
	ConsumerGraphFriendShip(ctx context.Context) error
	ConsumerFailedGraphFriendShip(ctx context.Context) error
}
