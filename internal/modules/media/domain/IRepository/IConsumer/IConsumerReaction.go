package IConsumer

import (
	"context"
)

type IConsumerReact interface {
	ConsumerReactAlbum(ctx context.Context)
	ConsumerFailedReactAlbum(ctx context.Context)
	CosumerReactStory(ctx context.Context)
	ConsumerFailedReactStory(ctx context.Context)
}
