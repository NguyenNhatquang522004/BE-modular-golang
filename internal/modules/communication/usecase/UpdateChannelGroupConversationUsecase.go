package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
)

type UpdateChannelGroupConversationUsecase struct {
	conversationRepo IRepositoryMongodb.IConversationsRepository
}

func NewUpdateChannelGroupConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository) *UpdateChannelGroupConversationUsecase {
	return &UpdateChannelGroupConversationUsecase{
		conversationRepo: conversationRepo,
	}
}
func (u *UpdateChannelGroupConversationUsecase) Execute(ctx context.Context, req *req.UpdateChannelConversationRequest) (*res.FailedChannelGroupConversationResponse, error) {
	data, err := u.conversationRepo.GetConversationByID(ctx, req.ChannelID)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ChannelID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   err,
		}, err
	}
	if data == nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ChannelID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   err,
		}, err
	}
	if data.Type != enum.TypeChannel {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ChannelID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   err,
		}, err
	}
	mapper.UpdateToEntityConversation(req.ConversationReq, data)
	err = u.conversationRepo.UpdateConversation(ctx, data)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ChannelID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   err,
		}, err
	}
	return &res.FailedChannelGroupConversationResponse{
		ChannelID:      req.ChannelID,
		ConversationID: req.ID,
		UserID:         "unknown",
		ErrorMessage:   nil,
	}, nil
}
