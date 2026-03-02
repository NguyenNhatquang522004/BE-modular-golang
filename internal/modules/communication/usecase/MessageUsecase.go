package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
)

type MessageUsecase struct {
	events events.EventBus
}

func NewMessageUsecase(events events.EventBus) *MessageUsecase {
	return &MessageUsecase{
		events: events,
	}
}

func (u *MessageUsecase) Execute(ctx context.Context, req *req.MessageRequest) (*res.FailedMessageResponse, error) {
	messagePayload := &communicationEvent.MessagePayload{
		ConversationID:   req.ConversationID,
		Bucket:           req.Bucket,
		MessageID:        req.MessageID,
		SenderID:         req.SenderID,
		Type:             req.Type,
		Content:          req.Content,
		Attachments:      req.Attachments,
		IsEdited:         req.IsEdited,
		ReplyToMessageID: req.ReplyToMessageID,
		StoryRefID:       req.StoryRefID,
		IsRevoked:        req.IsRevoked,
		CreatedAt:        req.CreatedAt,
		EventType:        req.EventType,
	}
	err := u.events.Publish(ctx, constants.TopicMessage.String(), req.ConversationID, req.EventType.String(), messagePayload)
	if err != nil {
		return &res.FailedMessageResponse{
			ConversationID: req.ConversationID,
			UserID:         req.SenderID.String(),
			ErrorMessage:   errors.New("failed to publish message"),
		}, err
	}
	return &res.FailedMessageResponse{
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID.String(),
		UserID:         req.SenderID.String(),
		ErrorMessage:   nil,
	}, nil
}
