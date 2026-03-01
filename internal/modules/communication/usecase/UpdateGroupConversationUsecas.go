package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
)

type UpdateGroupConversationUsecase struct {
	conversationRepo IRepositoryMongodb.IConversationsRepository
}

func NewUpdateGroupConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository) *UpdateGroupConversationUsecase {
	return &UpdateGroupConversationUsecase{
		conversationRepo: conversationRepo,
	}
}

func (u *UpdateGroupConversationUsecase) Execute(ctx context.Context, req *req.UpdateGroupConversationRequest) (*res.FailedChannelGroupConversationResponse, error) {
	data, err := u.conversationRepo.GetConversationByID(ctx, req.ID)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to retrieve conversation for update"),
		}, err
	}
	if data == nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Conversation not found for update"),
		}, nil
	}
	if data.Type != enum.TypeGroup {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Conversation is not a group conversation"),
		}, nil
	}
	mapper.UpdateToEntityConversation(req.ConversationReq, data)
	err = u.conversationRepo.UpdateConversation(ctx, data)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ID,
			ConversationID: req.ID,
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to update group conversation"),
		}, err
	}
	// Implement the logic to update a group conversation here
	// This may involve validating the request, checking if the conversation exists and is a group conversation, updating the conversation record in the database, etc.
	return &res.FailedChannelGroupConversationResponse{
		ChannelID:      "unknown",
		ConversationID: "unknown",
		UserID:         "unknown",
		ErrorMessage:   errors.New("Update group conversation use case not implemented yet"),
	}, nil
}
