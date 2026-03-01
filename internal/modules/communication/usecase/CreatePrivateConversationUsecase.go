package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
)

type CreatePrivateConversationUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
	pool                        IRepositoryShare.IWorkerPool
}

func NewCreatePrivateConversationUsecase(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository) *CreatePrivateConversationUsecase {
	return &CreatePrivateConversationUsecase{
		conversationRepo:            conversationRepo,
		conversationParticipantRepo: conversationParticipantRepo,
	}
}

func (u *CreatePrivateConversationUsecase) Execute(ctx context.Context, req *req.CreatePrivateConversationRequest) (*res.FailedPrivateConversationResponse, error) {
	// Implement the logic to create a private conversation here
	// This may involve validating the request, checking if the users exist, creating a conversation record in the database, etc.
	entityConversation := mapper.ToEntityConversation(req.ConversationReq)
	entityConversation.Type = enum.TypePrivate
	err := u.conversationRepo.CreateConversation(ctx, entityConversation)
	if err != nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      entityConversation.ID.Hex(),
			UserIDSenderFirst:   req.UserSenderFirst.UserID,
			UserIDReceiverFirst: req.UserReceiverFirst.UserID,
			ErrorMessage:        errors.New("Failed to create private conversation"),
		}, err
	}
	entitySender1 := mapper.ToEntityConversationParticipant(req.UserSenderFirst, entityConversation.ID.Hex())
	entitySender1.ConversationID = entityConversation.ID
	entitySender2 := mapper.ToEntityConversationParticipant(req.UserReceiverFirst, entityConversation.ID.Hex())
	entitySender2.ConversationID = entityConversation.ID
	countsucess, detailError, err := u.conversationParticipantRepo.CreateBulkConversationParticipants(ctx, []*entity.ConversationParticipant{entitySender1, entitySender2})
	if err != nil {
		return &res.FailedPrivateConversationResponse{
			ConversationID:      "dummy_conversation_id",
			UserIDSenderFirst:   req.UserSenderFirst.UserID,
			UserIDReceiverFirst: req.UserReceiverFirst.UserID,
			ErrorMessage:        errors.New("Failed to create private conversation"),
		}, nil
	}
	log.Printf("Successfully created %d conversation participants, with %d errors", countsucess, len(detailError))
	for _, errDetail := range detailError {
		log.Printf("Error creating conversation participant: %v", errDetail)
	}
	// For now, we will return a dummy response indicating failure
	return &res.FailedPrivateConversationResponse{
		ConversationID:      "dummy_conversation_id",
		UserIDSenderFirst:   req.UserSenderFirst.UserID,
		UserIDReceiverFirst: req.UserReceiverFirst.UserID,
		ErrorMessage:        errors.New("Failed to create private conversation"),
	}, nil
}
