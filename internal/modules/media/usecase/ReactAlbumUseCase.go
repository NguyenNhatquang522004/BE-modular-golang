package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ReactAlbumUseCase struct {
	// Define any dependencies needed for reacting to an album here
	eventbus events.EventBus
}

func NewReactAlbumUseCase(eventbus events.EventBus) *ReactAlbumUseCase {
	return &ReactAlbumUseCase{
		eventbus: eventbus,
	}
}

func (uc *ReactAlbumUseCase) Execute(ctx context.Context, req *req.ReactAlbumRequest) (*res.ReactAlbumResponse, error) {
	// Implement the logic for reacting to an album here
	payload := &mediaEvent.AlbumStatsPayload{
		AlbumID:    req.AlbumID,
		UserID:     req.UserID,
		AssetCount: req.AssetCount,
		Like:       req.Like,
		Love:       req.Love,
		Haha:       req.Haha,
		Wow:        req.Wow,
		Sad:        req.Sad,
		Angry:      req.Angry,
		EventType:  req.EventType,
	}
	err := uc.eventbus.Publish(ctx, constants.TopicAlbumStats.String(), req.AlbumID, req.EventType.String(), payload)
	if err != nil {
		return &res.ReactAlbumResponse{
			AlbumID:            req.AlbumID,
			AlbumsErrorMessage: err.Error(),
			EventType:          req.EventType.String(),
		}, err
	}
	return &res.ReactAlbumResponse{
		AlbumID:   req.AlbumID,
		EventType: req.EventType.String(),
	}, nil
}
