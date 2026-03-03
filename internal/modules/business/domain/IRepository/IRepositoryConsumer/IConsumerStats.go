package IRepositoryConsumer

import "context"

type IConsumerStats interface {
	ConsumerStatsPage(ctx context.Context) error
	ConsumerFailedStatsPage(ctx context.Context) error
}
