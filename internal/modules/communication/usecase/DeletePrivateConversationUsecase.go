package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type DeletePrivateConversationUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
	pool                        IRepositoryShare.IWorkerPool
}

func NewDeletePrivateConversationUsecase() *DeletePrivateConversationUsecase {
	return &DeletePrivateConversationUsecase{}
}

func (u *DeletePrivateConversationUsecase) Execute(ctx context.Context, req *req.DeletePrivateConversationRequest) (*res.FailedPrivateConversationResponse, error) {
	// Implement the logic to delete a private conversation here
	// This may involve validating the request, checking if the conversation exists, deleting the conversation record from the database, etc.
	err := u.conversationRepo.DeleteConversation(ctx, req.ConversationID)
	if err != nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      req.ConversationID,
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Failed to delete private conversation"),
		}, err
	}
	err = u.conversationParticipantRepo.DeleteConversationParticipantByConversationID(ctx, req.ConversationID)
	if err != nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      req.ConversationID,
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Failed to delete conversation participants"),
		}, err
	}

	return &res.FailedPrivateConversationResponse{
		ConversationID:      req.ConversationID,
		UserIDSenderFirst:   "unknown",
		UserIDReceiverFirst: "unknown",
		ErrorMessage:        errors.New("Private conversation deleted successfully"),
	}, nil
}
