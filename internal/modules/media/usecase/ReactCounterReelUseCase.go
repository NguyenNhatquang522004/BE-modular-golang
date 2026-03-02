package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ReactCounterReelUseCase struct {
	events events.EventBus
}

func NewReactCounterReelUseCase(events events.EventBus) *ReactCounterReelUseCase {
	return &ReactCounterReelUseCase{
		events: events,
	}
}

func (u *ReactCounterReelUseCase) Execute(ctx context.Context, req *req.ReactCounterReelRequest) (*res.FailedReelResponse, error) {
	payload := &mediaEvent.ReactCounterReelPayload{
		ReelID:    req.ReelID,
		UserID:    req.UserID,
		Comments:  req.Comments,
		Saves:     req.Saves,
		Shares:    req.Shares,
		EventType: req.EventType,
	}
	err := u.events.Publish(ctx, constants.TopicCounterReel.String(), req.ReelID, constants.Created.String(), payload)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to publish counter reel event",
		}, err
	}
	return &res.FailedReelResponse{
		ReelID:       req.ReelID,
		UserID:       req.UserID,
		ErrorMessage: "",
	}, nil
}
