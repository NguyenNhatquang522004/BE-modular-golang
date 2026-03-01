package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
)

type ReactMessageUsecase struct {
	events events.EventBus
}

func NewReactMessageUsecase(events events.EventBus) *ReactMessageUsecase {
	return &ReactMessageUsecase{
		events: events,
	}
}

func (u *ReactMessageUsecase) Execute(ctx context.Context, req *req.ReactMessageRequest) (*res.FailedReactMessageResponse, error) {
	err := u.events.Publish(ctx, string(constants.TopicReactMessage), req.ConversationID, req.EventType.String(), req)
	if err != nil {
		return &res.FailedReactMessageResponse{
			ConversationID: req.ConversationID,
			MessageID:      req.MessageID.String(),
			UserID:         req.UserID.String(),
			ReactionCode:   req.ReactionCode.String(),
			ErrorMessage:   err,
		}, err
	}
	return &res.FailedReactMessageResponse{
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID.String(),
		UserID:         req.UserID.String(),
		ReactionCode:   req.ReactionCode.String(),
		ErrorMessage:   nil,
	}, nil
}
