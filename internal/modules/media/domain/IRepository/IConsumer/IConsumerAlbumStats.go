package IConsumer

import "context"

type IConsumerAlbumStats interface {
	ConsumerAlbumStats(ctx context.Context) error
	ConsumerFailedAlbumStats(ctx context.Context) error
}
