package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type AddGroupToChannelParticipantUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewAddGroupToChannelParticipantUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository) *AddGroupToChannelParticipantUsecase {
	return &AddGroupToChannelParticipantUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
	}
}
func (u *AddGroupToChannelParticipantUsecase) Execute(ctx context.Context, req *req.AddGroupToChannelParticipantRequest) (*res.FailedParticipantResponse, error) {
	dataChannelConversation, err := u.conversationRepo.GetConversationByID(ctx, req.ChannelID)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.ChannelID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Failed to get channel conversation"),
		}, err
	}
	if dataChannelConversation == nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.ChannelID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Channel conversation not found"),
		}, nil
	}
	if dataChannelConversation.Type != sharedEnums.TypeChannel {
		return &res.FailedParticipantResponse{
			ConversationID: req.ChannelID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Conversation is not a channel conversation"),
		}, nil
	}
	dataGroupConversation, err := u.conversationRepo.GetConversationByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Failed to get group conversation"),
		}, err
	}
	if dataGroupConversation == nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Group conversation not found"),
		}, nil
	}
	if dataGroupConversation.Type != sharedEnums.TypeGroup {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Invalid group conversation type"),
		}, nil
	}
	if dataGroupConversation.RelatedGroupID.Hex() != dataChannelConversation.ID.Hex() {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Group conversation is not related to channel conversation"),
		}, nil
	}
	entityParticipant := mapper.ToEntityConversationParticipant(req.ConversationParticipantReq, dataGroupConversation.ID.Hex())
	err = u.conversationParticipantRepo.CreateConversationParticipant(ctx, entityParticipant)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Failed to add participant to group conversation"),
		}, err
	}

	return &res.FailedParticipantResponse{
		ConversationID: req.GroupID,
		UserID:         req.UserID,
		ErrorMessage:   nil,
	}, nil
}
