package IConsumer

import "context"

type IConsumerMusic interface {
	ConsumerMusic(ctx context.Context) error
	ConsumerFailedMusic(ctx context.Context) error
}
