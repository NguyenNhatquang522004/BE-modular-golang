package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type RemoveGroupParticipantUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewRemoveGroupParticipantUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository) *RemoveGroupParticipantUsecase {
	return &RemoveGroupParticipantUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
	}
}
func (u *RemoveGroupParticipantUsecase) Execute(ctx context.Context, req *req.RemoveGroupParticipantRequest) (*res.FailedParticipantResponse, error) {
	dataConversation, err := u.conversationRepo.GetConversationByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   err,
		}, err
	}
	if dataConversation == nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   nil,
		}, nil
	}
	dataConversation.ParticipantCount -= 1
	err = u.conversationRepo.UpdateConversation(ctx, dataConversation)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   err,
		}, err
	}
	err = u.conversationParticipantRepo.DeleteConversationParticipant(ctx, req.GroupID, req.UserID)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   err,
		}, err
	}
	return &res.FailedParticipantResponse{
		ConversationID: req.GroupID,
		UserID:         req.UserID,
		ErrorMessage:   nil,
	}, nil
}
