package kafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
)

type BlockMessage struct {
	eventBus events.EventBus
}

func NewBlockMessage(eventBus events.EventBus) *BlockMessage {
	return &BlockMessage{
		eventBus: eventBus,
	}
}
func (bm *BlockMessage) PublishBlockCreateMessage(ctx context.Context, payload socialEvent.BlockCreatePayload) error {
	// Implement the logic to publish a block create message
	err := bm.eventBus.Publish(ctx, constants.TopicBlock.String(), payload.Blocker_UserID, constants.Created.String(), payload)
	if err != nil {
		return err
	}
	return nil
}

func (bm *BlockMessage) PublishBlockUpdateMessage(ctx context.Context, payload socialEvent.BlockUpdatePayload) error {
	// Implement the logic to publish a block update message
	err := bm.eventBus.Publish(ctx, constants.TopicBlock.String(), payload.Blocker_UserID, constants.Updated.String(), payload)
	if err != nil {
		return err
	}
	return nil
}

func (bm *BlockMessage) PublishBlockDeleteMessage(ctx context.Context, payload socialEvent.BlockDeletePayload) error {
	// Implement the logic to publish a block delete message
	err := bm.eventBus.Publish(ctx, constants.TopicBlock.String(), payload.Blocker_UserID, constants.Deleted.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
