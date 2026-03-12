package IConsumer

import "context"

type IConsumerAlbum interface {
	ConsumerAlbum(ctx context.Context) error
	ConsumerFailedAlbum(ctx context.Context) error
}
