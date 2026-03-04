package IRepositoryConsumer

import "context"

type IConsumerPageDailyMetrics interface {
	// Define the methods for the ConsumerPageDailyMetricsRepository interface here
	ConsumerPageDailyMetric(ctx context.Context) error
	ConsumerFailedPageDailyMetric(ctx context.Context) error
}
