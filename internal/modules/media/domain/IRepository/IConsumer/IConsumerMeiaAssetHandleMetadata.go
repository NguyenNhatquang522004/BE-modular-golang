package IConsumer

import "context"

type IConsumerMeiaAssetHandleMetadata interface {
	ConsumerMediaAssetHandleMetadata(ctx context.Context) error
	ConsumerFailedMediaAssetHandleMetadata(ctx context.Context) error
}
