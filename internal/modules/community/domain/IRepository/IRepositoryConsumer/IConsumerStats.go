package IRepositoryConsumer

import "context"

type IConsumerStats interface {
	ConsumerGroupStats(ctx context.Context)
	ConsumerFailedGroupStats(ctx context.Context)
	ConsumerEventStats(ctx context.Context)
	ConsumerFailedEventStats(ctx context.Context)
	ConsumerDownloadGroupFile(ctx context.Context)
	ConsumerFailedDownloadGroupFile(ctx context.Context)
}
