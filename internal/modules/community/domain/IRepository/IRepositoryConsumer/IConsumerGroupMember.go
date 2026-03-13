package IRepositoryConsumer

import "context"

type IConsumerGroupMember interface {
	ConsumerGroupMember(ctx context.Context) error
	ConsumerFailedGroupMember(ctx context.Context) error
}
