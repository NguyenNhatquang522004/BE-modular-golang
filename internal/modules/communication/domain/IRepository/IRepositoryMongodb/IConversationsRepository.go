package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

type IConversationsRepository interface {
	CreateConversation(ctx context.Context, conversation *entity.Conversation) error
	CreateBulkConversations(ctx context.Context, conversations []entity.Conversation) (int64, []*mongodbErrors.BulkError, error)
	GetConversationByID(ctx context.Context, id string) (*entity.Conversation, error)
	GetConversationByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetConversationByCreatorID(ctx context.Context, creatorID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetConversationByOwnerID(ctx context.Context, ownerID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateConversation(ctx context.Context, conversation *entity.Conversation) error
	UpdateBulkConversations(ctx context.Context, conversations []entity.Conversation) (int64, []*mongodbErrors.BulkError, error)
	DeleteConversation(ctx context.Context, id string) error
	DeleteBulkConversations(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
	CheckConversationExists(ctx context.Context, userIDOne string, userIDTwo string) (*entity.Conversation, error)
}
