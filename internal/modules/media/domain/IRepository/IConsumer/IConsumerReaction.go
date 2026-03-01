package IConsumer

import (
	"context"
)

type IConsumerReact interface {
	ConsumerReactAlbum(ctx context.Context)
	ConsumerFailedReactAlbum(ctx context.Context)
	CosumerReactStory(ctx context.Context)
	ConsumerFailedReactStory(ctx context.Context)
	ConsumerReactReel(ctx context.Context)
	ConsumerFailedReactReel(ctx context.Context)
	ConsumerEntityReaction(ctx context.Context)
	ConsumerFailedEntityReaction(ctx context.Context)
	ConsumerCounterReel(ctx context.Context)
	ConsumerFailedCounterReel(ctx context.Context)
	ConsumerReactLive(ctx context.Context)
	ConsumerFailedReactLive(ctx context.Context)
	ConsumerCounterLive(ctx context.Context)
	ConsumerFailedCounterLive(ctx context.Context)
	ConsumerCounterReplyStory(ctx context.Context)
	ConsumerFailedCounterReplyStory(ctx context.Context)
	
}
