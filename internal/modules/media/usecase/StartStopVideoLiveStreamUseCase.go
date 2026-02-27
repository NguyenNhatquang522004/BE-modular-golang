package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type StartStopVideoLiveStreamUseCase struct {
	eventsbus events.EventBus
}

func NewStartStopVideoLiveStreamUseCase(eventsbus events.EventBus) *StartStopVideoLiveStreamUseCase {
	return &StartStopVideoLiveStreamUseCase{
		eventsbus: eventsbus,
	}
}

func (uc *StartStopVideoLiveStreamUseCase) Execute(ctx context.Context, req *req.StartStopVideoLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	// Implement logic to start a video live stream here
	// This may involve validating the request, interacting with repositories,
	// and returning an appropriate response or
	err := uc.eventsbus.Publish(ctx, constants.TopicStartStopLive.String(), req.LiveSessionID, req.EventType.String(), req)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        req.OwnerID,
			ErrorMessage:  err.Error(),
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: req.LiveSessionID,
		UserID:        req.OwnerID,
		ErrorMessage:  "",
	}, nil
}
