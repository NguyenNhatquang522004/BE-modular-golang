package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
)

type CreatePageFollowerUsecase struct {
	events events.EventBus
}

func NewCreatePageFollowerUsecase(events events.EventBus) *CreatePageFollowerUsecase {
	return &CreatePageFollowerUsecase{
		events: events,
	}
}
func (u *CreatePageFollowerUsecase) Execute(ctx context.Context, req *req.CreatePageFollowerRequest) (*res.FailedPageFollowerResponse, error) {
	payload := &businessEvent.FollowerPagePayload{
		PageID: req.PageID,
		UserID: req.UserID,
	}
	if req.Settings != nil {
		payload.Settings.NotificationLevel = req.Settings.NotificationLevel
		payload.Settings.IsFavorite = req.Settings.IsFavorite
	}
	err := u.events.Publish(ctx, constants.TopicFollowerPage.String(), req.PageID, constants.Created.String(), payload)
	if err != nil {
		return &res.FailedPageFollowerResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to publish follower page event",
		}, err
	}
	return &res.FailedPageFollowerResponse{
		PageID:       req.PageID,
		UserID:       req.UserID,
		ErrorMessage: "Failed to create page follower",
	}, nil
}
