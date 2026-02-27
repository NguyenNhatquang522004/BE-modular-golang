package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type CreateLiveStreamUseCase struct {
	eventsbus       events.EventBus
	liveSessionRepo IRepositoryMongodb.ILiveSessionRepository
}

func NewCreateLiveStreamUseCase(eventsbus events.EventBus, liveSessionRepo IRepositoryMongodb.ILiveSessionRepository) *CreateLiveStreamUseCase {
	return &CreateLiveStreamUseCase{
		eventsbus:       eventsbus,
		liveSessionRepo: liveSessionRepo,
	}
}

func (uc *CreateLiveStreamUseCase) Execute(ctx context.Context, req *req.CreateLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	// Implement the logic to create a live stream here
	// You can interact with repositories, perform validations, etc.

	// For demonstration, let's assume the creation is successful and return a response
	entity, err := mapper.ToEntityLiveSession(req.LiveSessionReq)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: entity.ID.Hex(),
			UserID:        req.HostUserID,
			ErrorMessage:  "Failed to map request to entity: " + err.Error(),
		}, err
	}
	err = uc.liveSessionRepo.CreateLiveSession(ctx, entity)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: entity.ID.Hex(),
			UserID:        req.HostUserID,
			ErrorMessage:  "Failed to create live session: " + err.Error(),
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: entity.ID.Hex(),
		UserID:        req.HostUserID,
		ErrorMessage:  "",
	}, nil
}
