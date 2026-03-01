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

type AddChannelParticipantUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewAddChannelParticipantUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository) *AddChannelParticipantUsecase {
	return &AddChannelParticipantUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
	}
}
func (u *AddChannelParticipantUsecase) Execute(ctx context.Context, req *req.AddChannelParticipantRequest) (*res.FailedParticipantResponse, error) {
	dataConversation, err := u.conversationRepo.GetConversationByID(ctx, req.ConversationID)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.ConversationID,
			UserID:         req.UserID,
			ErrorMessage:   err,
		}, err
	}
	if dataConversation == nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.ConversationID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Conversation not found"),
		}, nil
	}
	if dataConversation.Type != enum.TypeChannel {
		return &res.FailedParticipantResponse{
			ConversationID: req.ConversationID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Conversation is not a channel conversation"),
		}, nil
	}
	entityParticipant := mapper.ToEntityConversationParticipant(req.ConversationParticipantReq, dataConversation.ID.Hex())
	err = u.conversationParticipantRepo.CreateConversationParticipant(ctx, entityParticipant)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.ConversationID,
			UserID:         req.UserID,
			ErrorMessage:   errors.New("Failed to add participant to channel conversation"),
		}, err
	}
	return &res.FailedParticipantResponse{
		ConversationID: req.ConversationID,
		UserID:         req.UserID,
		ErrorMessage:   nil,
	}, nil
}
