package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

type IMessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.Message) error
	CreateBulkMessages(ctx context.Context, messages []*entity.Message) (int64, []*cassandraErrors.MessageBulkError, error)
	UpdateMessage(ctx context.Context, message *entity.Message) error
	UpdateBulkMessages(ctx context.Context, messages []*entity.Message) (int64, []*cassandraErrors.MessageBulkError, error)
	DeleteMessagesByConversationID(ctx context.Context, conversationID string) error
	DeleteBulkMessagesByConversationID(ctx context.Context, conversationIDs []string) error
	DeleteAllMessagesByMessageConversationID(ctx context.Context, conversationID string) error
	DeleteMessage(ctx context.Context, conversationID string, bucket int, messageID string) error
	DeleteBulkMessages(ctx context.Context, conversationID string, bucket int, messageIDs []string) (int64, []*cassandraErrors.MessageBulkError, error)
	GetMessagesByConversationIDs(ctx context.Context, conversationID string, bucket int, cursor string, limit int) (*dto.PaginationRes, error)
	GetMessagesByConversationID(ctx context.Context, conversationID string, bucket int) (*entity.Message, error)
	GetMessagesByMessageID(ctx context.Context, conversationID string, bucket int, messageID string) (*entity.Message, error)
}
