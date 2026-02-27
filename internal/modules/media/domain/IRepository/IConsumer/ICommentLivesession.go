package IConsumer

import "context"

type ICommentLiveStreamUseCase interface {
	ConsumerLiveComment(ctx context.Context)
}
