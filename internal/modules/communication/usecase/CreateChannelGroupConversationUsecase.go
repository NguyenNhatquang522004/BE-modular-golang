package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
)

type CreateChannelGroupConversationUsecase struct {
	conversationRepo            IRepositoryMongodb.IConversationsRepository
	conversationParticipantRepo IRepositoryMongodb.IConversationParticipantsRepository
}

func NewCreateChannelGroupConversationUsecase() *CreateChannelGroupConversationUsecase {
	return &CreateChannelGroupConversationUsecase{}
}
func (u *CreateChannelGroupConversationUsecase) Execute(ctx context.Context, req *req.CreateChannelGroupConversationRequest) (*res.FailedPrivateConversationResponse, error) {
	// Implement the logic to create a channel group conversation here
	// This may involve validating the request, creating a new conversation record in the database, adding participants to the conversation, etc.
	entityConversation := mapper.ToEntityConversation(req.ConversationReq)
	entityConversation.Type = enum.TypeChannel
	entityParticipants := mapper.ToEntityBulkConversationParticipant(req.User, entityConversation.ID.Hex())
	entityConversation.ParticipantCount = len(entityParticipants)
	err := u.conversationRepo.CreateConversation(ctx, entityConversation)
	if err != nil {
		log.Printf("Error creating conversation: %v", err)
		return &res.FailedPrivateConversationResponse{
			ConversationID:      entityConversation.ID.Hex(),
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Failed to create channel group conversation"),
		}, err
	}
	_, _, err = u.conversationParticipantRepo.CreateBulkConversationParticipants(ctx, entityParticipants)
	if err != nil {
		log.Printf("Error creating conversation participants: %v", err)
		// Optionally, you might want to roll back the conversation creation if participant creation fails
		return &res.FailedPrivateConversationResponse{
			ConversationID:      entityConversation.ID.Hex(),
			UserIDSenderFirst:   "unknown",
			UserIDReceiverFirst: "unknown",
			ErrorMessage:        errors.New("Failed to create channel group conversation participants"),
		}, err
	}
	return &res.FailedPrivateConversationResponse{
		ConversationID:      "unknown",
		UserIDSenderFirst:   "unknown",
		UserIDReceiverFirst: "unknown",
		ErrorMessage:        errors.New("Failed to create channel group conversation"),
	}, nil
}
