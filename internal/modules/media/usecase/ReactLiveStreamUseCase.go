package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ReactLiveStreamUseCase struct {
	eventsbus events.EventBus
}

func NewReactLiveStreamUseCase(eventsbus events.EventBus) *ReactLiveStreamUseCase {
	return &ReactLiveStreamUseCase{
		eventsbus: eventsbus,
	}
}

func (u *ReactLiveStreamUseCase) Execute(ctx context.Context, req *req.ReactLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	payload := &mediaEvent.ReactLiveStreamPayload{
		LiveSessionID: req.LiveSessionID,
		UserID:        req.UserID,
		Total:         req.Total,
		TargetType:    req.TargetType,
		ReactionCode:  req.ReactionCode,
		CreatedAt:     req.CreatedAt,
		EventType:     req.EventType,
	}
	err := u.eventsbus.Publish(ctx, constants.TopicReactLive.String(), req.LiveSessionID, string(req.EventType), payload)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        req.UserID,
			ErrorMessage:  err.Error(),
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: req.LiveSessionID,
		UserID:        req.UserID,
		ErrorMessage:  "",
	}, nil
}
