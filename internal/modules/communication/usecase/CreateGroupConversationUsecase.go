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

type CreateGroupConversationUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewCreateGroupConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository) *CreateGroupConversationUsecase {
	return &CreateGroupConversationUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
	}
}

func (u *CreateGroupConversationUsecase) Execute(ctx context.Context, req *req.CreateGroupConversationRequest) (*res.FailedChannelGroupConversationResponse, error) {

	entityConversation := mapper.ToEntityConversation(req.ConversationReq)
	entityConversation.Type = enum.TypeGroup
	entity := mapper.ToEntityBulkConversationParticipant(req.User, entityConversation.ID.Hex())
	entityConversation.ParticipantCount = len(entity)
	err := u.conversationRepo.CreateConversation(ctx, entityConversation)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      "unknown",
			ConversationID: "unknown",
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to create group conversation"),
		}, err
	}
	_, _, err = u.conversationParticipantRepo.CreateBulkConversationParticipants(ctx, entity)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      "unknown",
			ConversationID: entityConversation.ID.Hex(),
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to create group conversation"),
		}, nil
	}
	// For now, we will return a dummy response indicating failure
	// Optionally, you can also create conversation participants here using u.conversationParticipantRepo
	return &res.FailedChannelGroupConversationResponse{
		ChannelID:      "unknown",
		ConversationID: entityConversation.ID.Hex(),
		UserID:         "unknown",
		ErrorMessage:   errors.New("Create group conversation use case not implemented yet"),
	}, nil
}
