package IConsumer

import (
	"context"
)

type IConsumerGraphUser interface {
	ConsumerGraphUser(ctx context.Context) error
	ConsumerFailedGraphUser(ctx context.Context) error
}
