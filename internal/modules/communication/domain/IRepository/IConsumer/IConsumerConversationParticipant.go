package IConsumer

import "context"

type IConsumerConversationParticipant interface {
	ConsumerConversationParticipantEvents(ctx context.Context) error
	ConsumerFailedConversationParticipantEvents(ctx context.Context) error
}
