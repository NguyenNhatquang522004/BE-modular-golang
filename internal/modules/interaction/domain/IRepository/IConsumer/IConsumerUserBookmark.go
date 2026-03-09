package IConsumer

import "context"

type IConsumerUserBookmark interface {
	ConsumerUserBookmark(ctx context.Context) error
	ConsumerFailedUserBookmark(ctx context.Context) error
}
