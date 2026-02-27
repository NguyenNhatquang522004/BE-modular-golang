package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type DeleteLiveStreamUseCase struct {
	eventsbus       events.EventBus
	liveSessionRepo IRepositoryMongodb.ILiveSessionRepository
}

func NewDeleteLiveStreamUseCase(eventsbus events.EventBus, liveSessionRepo IRepositoryMongodb.ILiveSessionRepository) *DeleteLiveStreamUseCase {
	return &DeleteLiveStreamUseCase{
		eventsbus:       eventsbus,
		liveSessionRepo: liveSessionRepo,
	}
}

func (uc *DeleteLiveStreamUseCase) Execute(ctx context.Context, req *req.DeleteLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	// Implement the logic to delete a live stream here
	// You can interact with repositories, perform validations, etc.
	// For demonstration, let's assume the deletion is successful and return a response
	err := uc.liveSessionRepo.DeleteLiveSession(ctx, req.LiveSessionID)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: req.LiveSessionID,
			UserID:        req.UserID,
			ErrorMessage:  "Failed to delete live session: " + err.Error(),
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: req.LiveSessionID,
		UserID:        req.UserID,
		ErrorMessage:  "",
	}, nil
}
