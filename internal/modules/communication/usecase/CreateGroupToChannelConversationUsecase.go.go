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

type CreateGroupToChannelConversationUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewCreateGroupToChannelConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository) *CreateGroupToChannelConversationUsecase {
	return &CreateGroupToChannelConversationUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
	}
}

func (u *CreateGroupToChannelConversationUsecase) Execute(ctx context.Context, req *req.CreateGroupToChannelConversationRequest) (*res.FailedChannelGroupConversationResponse, error) {
	// Implement the logic to create a group conversation here
	// This may involve validating the request, checking if the users exist, creating a conversation record in the database, etc.
	dataGetChannel, err := u.conversationRepo.GetConversationByID(ctx, req.ConversationChannelID)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ConversationChannelID,
			ConversationID: "unknown",
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to get channel conversation"),
		}, err
	}
	if dataGetChannel == nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ConversationChannelID,
			ConversationID: "unknown",
			UserID:         "unknown",
			ErrorMessage:   errors.New("Channel conversation not found"),
		}, nil
	}
	entityConversation := mapper.ToEntityConversation(req.ConversationReq)
	entityConversation.Type = enum.TypeGroup
	entityConversation.RelatedChannelID = &dataGetChannel.ID
	entityConversation.ParticipantCount = len(req.User)
	err = u.conversationRepo.CreateConversation(ctx, entityConversation)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ConversationChannelID,
			ConversationID: "unknown",
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to create group conversation"),
		}, err
	}
	entity := mapper.ToEntityBulkConversationParticipant(req.User, entityConversation.ID.Hex())
	_, _, err = u.conversationParticipantRepo.CreateBulkConversationParticipants(ctx, entity)
	if err != nil {
		return &res.FailedChannelGroupConversationResponse{
			ChannelID:      req.ConversationChannelID,
			ConversationID: entityConversation.ID.Hex(),
			UserID:         "unknown",
			ErrorMessage:   errors.New("Failed to create group conversation"),
		}, nil
	}
	// For now, we will return a dummy response indicating failure
	// Optionally, you can also create conversation participants here using u.conversationParticipantRepo
	return &res.FailedChannelGroupConversationResponse{
		ChannelID:      req.ConversationChannelID,
		ConversationID: entityConversation.ID.Hex(),
		UserID:         "unknown",
		ErrorMessage:   errors.New("Create group conversation use case not implemented yet"),
	}, nil
}
