package IConsumer

import "context"

type IConsumerLiveStream interface {
	ConsumeStartLiveStream(ctx context.Context)
}
