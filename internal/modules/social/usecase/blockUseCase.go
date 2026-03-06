package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
)

type BlockUseCase struct {
	events events.EventBus
}

func NewBlockUseCase(events events.EventBus) *BlockUseCase {
	return &BlockUseCase{
		events: events,
	}
}

func (c *BlockUseCase) Execute(ctx context.Context, req *req.BlockRequest) error {
	// Implement the logic for executing the block use case
	payload := &socialEvent.BlockUserPayload{
		ID:            req.ID,
		BlockerUserID: req.BlockerUserID,
		BlockedUserID: req.BlockedUserID,
		TypeBlock:     req.TypeBlock,
		EventType:     req.EventType,
		CreatedAt:     req.CreatedAt,
		UpdatedAt:     req.UpdatedAt,
	}
	err := c.events.Publish(ctx, constants.TopicBlockUser.String(), req.BlockerUserID, req.EventType.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
