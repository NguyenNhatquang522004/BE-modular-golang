package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type CounterLiveStreamUseCase struct {
	eventbus events.EventBus
}

func NewCounterLiveStreamUseCase(eventbus events.EventBus) *CounterLiveStreamUseCase {
	return &CounterLiveStreamUseCase{
		eventbus: eventbus,
	}
}

func (c *CounterLiveStreamUseCase) Execute(ctx context.Context, req *req.CounterLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	payload := &mediaEvent.CoutnerLiveStreamPayload{
		LiveSessionID: req.LiveSessionID,
		Comments:      req.Comments,
		Views:         req.Views,
		EventType:     req.EventType,
	}
	err := c.eventbus.Publish(ctx, constants.TopicCounterLive.String(), req.LiveSessionID, req.EventType.String(), payload)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			ErrorMessage:  "Failed to publish counter event to event bus",
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: req.LiveSessionID,
		ErrorMessage:  "Failed to update counter for live stream",
	}, nil
}
