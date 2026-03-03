package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IConversationParticipantsRepository interface {
	CreateConversationParticipant(ctx context.Context, participant *entity.ConversationParticipant) error
	CreateBulkConversationParticipants(ctx context.Context, participants []*entity.ConversationParticipant) (int64, []*mongodbErrors.BulkError, error)
	GetConversationParticipant(ctx context.Context, conversationID string, userID string) (*entity.ConversationParticipant, error)
	GetConversationParticipantsByConversationID(ctx context.Context, conversationID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetConversationParticipantsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateConversationParticipant(ctx context.Context, participant *entity.ConversationParticipant) error
	UpdateBulkConversationParticipants(ctx context.Context, participants []entity.ConversationParticipant) (int64, []*mongodbErrors.BulkError, error)
	DeleteConversationParticipantByConversationID(ctx context.Context, conversationID string) error
	DeleteBulkConversationParticipantsByConversationIDs(ctx context.Context, conversationIDs []primitive.ObjectID) error
	DeleteConversationParticipant(ctx context.Context, conversationID string, userID string) error
	DeleteBulkConversationParticipants(ctx context.Context, conversationIDs []string, userIDs []string) (int64, []*mongodbErrors.BulkError, error)
}
