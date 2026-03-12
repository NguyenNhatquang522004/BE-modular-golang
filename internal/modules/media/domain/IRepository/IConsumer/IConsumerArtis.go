package IConsumer

import "context"

type IConsumerArtist interface {
	ConsumerArtist(ctx context.Context) error
	ConsumerFailedArtist(ctx context.Context) error
}
