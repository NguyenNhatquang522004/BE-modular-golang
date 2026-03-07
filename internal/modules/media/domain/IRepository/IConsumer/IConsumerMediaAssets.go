package IConsumer

import "context"

type IConsumerMediaAssets interface {
	ConsumeMediaAsset(ctx context.Context) error
	ConsumerFailedMediaAsset(ctx context.Context) error
}
