package IConsumer

import "context"

type IConsumerConversation interface {
	ConsumerConversation(ctx context.Context) error
	ConsumerFailedConversation(ctx context.Context) error
}
