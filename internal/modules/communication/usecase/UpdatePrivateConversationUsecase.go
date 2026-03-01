package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type UpdatePrivateConversationUsecase struct {
	conversationRepo IRepositoryMongodb.IConversationsRepository
}

func NewUpdatePrivateConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository) *UpdatePrivateConversationUsecase {
	return &UpdatePrivateConversationUsecase{
		conversationRepo: conversationRepo,
	}
}

func (u *UpdatePrivateConversationUsecase) Execute(ctx context.Context, req *req.UpdatePrivateConversationRequest) (*res.FailedPrivateConversationResponse, error) {

	data, err := u.conversationRepo.GetConversationByID(ctx, req.ID)
	if err != nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      req.ID,
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Failed to retrieve conversation for update"),
		}, err
	}
	if data == nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      req.ID,
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Conversation not found for update"),
		}, nil
	}
	mapper.UpdateToEntityConversation(req.ConversationReq, data)
	err = u.conversationRepo.UpdateConversation(ctx, data)
	if err != nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      req.ID,
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Failed to update private conversation"),
		}, err
	}
	return &res.FailedPrivateConversationResponse{
		ConversationID:      req.ID,
		UserIDSenderFirst:   "unknown",
		UserIDReceiverFirst: "unknown",
		ErrorMessage:        errors.New("Private conversation updated successfully"),
	}, nil
}
