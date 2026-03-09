package IConsumer

import "context"

type IConsumerEditTarget interface {
	ConsumerEditTarget(ctx context.Context) error
	ConsumerFailEditTarget(ctx context.Context) error
}
