package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type UpdateLiveStreamUseCase struct {
	eventsbus       events.EventBus
	liveSessionRepo IRepositoryMongodb.ILiveSessionRepository
}

func NewUpdateLiveStreamUseCase(eventsbus events.EventBus, liveSessionRepo IRepositoryMongodb.ILiveSessionRepository) *UpdateLiveStreamUseCase {
	return &UpdateLiveStreamUseCase{
		eventsbus:       eventsbus,
		liveSessionRepo: liveSessionRepo,
	}
}
func (u *UpdateLiveStreamUseCase) Execute(ctx context.Context, req *req.UpdateLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	// Implement the logic to update a live stream here
	// You can use the req parameter to get the necessary information for updating the live stream
	// Return a FailedLiveStreamResponse if there is an error during the update processd
	data, err := u.liveSessionRepo.GetLiveSessionByID(ctx, req.LiveSessionID)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        *req.HostUserID,
			ErrorMessage:  "Failed to retrieve live session: " + err.Error(),
		}, err
	}
	if data == nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        *req.HostUserID,
			ErrorMessage:  "Live session not found",
		}, nil
	}
	if data.HostUserID != *req.HostUserID {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        *req.HostUserID,
			ErrorMessage:  "Unauthorized to update this live session",
		}, nil
	}
	mapper.UpdateToEntityLiveSession(req.UpdateLiveSessionReq, data)
	err = u.liveSessionRepo.UpdateLiveSession(ctx, data)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        *req.HostUserID,
			ErrorMessage:  "Failed to update live session: " + err.Error(),
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: req.LiveSessionID,
		UserID:        *req.HostUserID,
		ErrorMessage:  "Live session updated successfully",
	}, nil
}
