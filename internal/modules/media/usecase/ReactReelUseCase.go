package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/gocql/gocql"
)

type ReactReelUseCase struct {
	eventBus events.EventBus
}

func NewReactReelUseCase(eventBus events.EventBus) *ReactReelUseCase {
	return &ReactReelUseCase{
		eventBus: eventBus,
	}
}

func (uc *ReactReelUseCase) Execute(ctx context.Context, req *req.ReactReelRequest) (*res.FailedReelResponse, error) {
	converUserID, err := gocql.ParseUUID(req.UserID)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       req.UserID,
			ErrorMessage: "Invalid UserID format",
		}, err
	}
	payload := &interactionEvent.EntityReactionPayload{
		TargetID:     req.ReelID,
		UserID:       converUserID,
		TargetType:   req.TargetType,
		ReactionCode: req.ReactionCode,
		CreatedAt:    req.CreatedAt,
		Topic:        req.EventType,
	}
	payload2 := &mediaEvent.ReactReelPayload{
		ReelID:       req.ReelID,
		UserID:       req.UserID,
		Total:        req.Total,
		TargetType:   req.TargetType,
		ReactionCode: req.ReactionCode,
		CreatedAt:    req.CreatedAt,
		EventType:    req.EventType,
	}
	err = uc.eventBus.Publish(ctx, constants.TopicReactReel.String(), req.ReelID, req.EventType.String(), payload2)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to publish event",
		}, err
	}
	err = uc.eventBus.Publish(ctx, constants.TopicEntityReaction.String(), req.ReelID, req.EventType.String(), payload)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to publish event",
		}, err
	}
	return nil, nil
}
