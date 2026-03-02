package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
)

type MessageStatedUsecase struct {
	events events.EventBus
}

func NewMessageStatedUsecase(events events.EventBus) *MessageStatedUsecase {
	return &MessageStatedUsecase{
		events: events,
	}
}
func (u *MessageStatedUsecase) Execute(ctx context.Context, req *req.MessageStateRequest) (*res.FailedMessageResponse, error) {
	payload := &communicationEvent.MessageStatePayload{
		ConversationID: req.ConversationID,
		UserID:         req.UserID,
		MessageID:      req.MessageID,
		LastReadAt:     req.LastReadAt,
		EventType:      req.EventType,
	}
	err := u.events.Publish(ctx, string(constants.TopicStateMessage), req.ConversationID, req.EventType.String(), payload)
	if err != nil {
		return &res.FailedMessageResponse{
			ConversationID: req.ConversationID,
			UserID:         req.UserID,
			MessageID:      req.MessageID,
			ErrorMessage:   err,
		}, err
	}
	return &res.FailedMessageResponse{
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		UserID:         req.UserID,
		ErrorMessage:   nil,
	}, nil
}
