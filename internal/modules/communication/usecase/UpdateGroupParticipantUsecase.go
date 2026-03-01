package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type UpdateGroupParticipantUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewUpdateGroupParticipantUsecase() *UpdateGroupParticipantUsecase {
	return &UpdateGroupParticipantUsecase{}
}
func (u *UpdateGroupParticipantUsecase) Execute(ctx context.Context, req *req.UpdateGroupParticipantRequest) (*res.FailedParticipantResponse, error) {
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
	dataConversationParticipant, err := u.conversationParticipantRepo.GetConversationParticipant(ctx, req.GroupID, req.UserID)
	if err != nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   err,
		}, err
	}
	if dataConversationParticipant == nil {
		return &res.FailedParticipantResponse{
			ConversationID: req.GroupID,
			UserID:         req.UserID,
			ErrorMessage:   nil,
		}, nil
	}
	mapper.UpdateToEntityConversationParticipant(req.ConversationParticipantReq, dataConversationParticipant)
	err = u.conversationParticipantRepo.UpdateConversationParticipant(ctx, dataConversationParticipant)
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
