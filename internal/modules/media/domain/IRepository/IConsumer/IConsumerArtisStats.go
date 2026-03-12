package IConsumer

import "context"

type IConsumerArtistStats interface {
	ConsumerArtistStats(ctx context.Context) error
	ConsumerFailedArtistStats(ctx context.Context) error
}
