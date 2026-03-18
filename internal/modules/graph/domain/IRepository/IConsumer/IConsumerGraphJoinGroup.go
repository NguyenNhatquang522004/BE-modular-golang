package IConsumer

import "context"

type IConsumerGraphJoinGroup interface {
	ConsumerGraphJoinGroup(ctx context.Context) error
	ConsumerFailedGraphJoinGroup(ctx context.Context) error
}
