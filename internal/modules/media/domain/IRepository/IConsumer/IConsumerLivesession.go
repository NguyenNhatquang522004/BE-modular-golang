package IConsumer

import "context"

type IConsumerLiveSession interface {
	ConsumeStartLiveStream(ctx context.Context)
}
