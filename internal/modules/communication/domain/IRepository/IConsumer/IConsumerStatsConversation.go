package IConsumer

import "context"

type IConsumerStatsConversation interface {
	ConsumerStatsConversation(ctx context.Context) error
	ConsumerFailedStatsConversation(ctx context.Context) error
}
