package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ReelStatsUseCase struct {
	events events.EventBus
}

func NewReelStatsUseCase(events events.EventBus) *ReelStatsUseCase {
	return &ReelStatsUseCase{
		events: events,
	}
}

func (u *ReelStatsUseCase) Execute(ctx context.Context, req *req.ReelStatsRequest) (*res.FailedReelResponse, error) {
	// Implement the logic to handle the reel stats request here
	// This may involve validating the request, processing the data, and returning an appropriate response
	payload := &mediaEvent.ReelStatsPayload{
		ReelID:    req.ReelID,
		UserID:    req.UserID,
		Views:     req.Views,
		Like:      req.Like,
		Love:      req.Love,
		Haha:      req.Haha,
		Wow:       req.Wow,
		Sad:       req.Sad,
		Angry:     req.Angry,
		Comments:  req.Comments,
		Saves:     req.Saves,
		Shares:    req.Shares,
		EventType: req.EventType,
	}
	err := u.events.Publish(ctx, constants.TopicReelStats.String(), req.ReelID, req.EventType.String(), payload)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       req.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedReelResponse{
		ReelID:       req.ReelID,
		UserID:       req.UserID,
		ErrorMessage: "",
	}, nil
}
