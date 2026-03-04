package IConsumer

import "context"

type IConsumerProfile interface {
	ConsumerProfile(ctx context.Context) error
	ConsumerFailedProfile(ctx context.Context) error
}
