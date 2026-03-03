package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

type IConversationReadStateRepository interface {
	CreateConversationReadState(ctx context.Context, readState *entity.ConversationReadState) error
	CreateBulkConversationReadStates(ctx context.Context, readStates []*entity.ConversationReadState) (int64, []*cassandraErrors.ConversationReadStateBulkError, error)
	UpdateConversationReadState(ctx context.Context, readState *entity.ConversationReadState) error
	UpsertConversationReadState(ctx context.Context, readState *entity.ConversationReadState) error
	UpdateBulkConversationReadStates(ctx context.Context, readStates []*entity.ConversationReadState) (int64, []*cassandraErrors.ConversationReadStateBulkError, error)
	GetConversationReadState(ctx context.Context, conversationID string, userID string) error
	DeleteConversationReadState(ctx context.Context, conversationID string, userID string) error
	DeleteBulkConversationReadStates(ctx context.Context, conversationIDs []string, userIDs []string) (int64, []*cassandraErrors.ConversationReadStateBulkError, error)
	DeleteConversationReadStatesByConversationID(ctx context.Context, conversationID string) error
	DeleteConversationReadStatesByUserID(ctx context.Context, conversationID string, userID string) error
	DeleteBulkConversationReadStatesByManyUserID(ctx context.Context, conversationID string, userIDs []string) (int64, []*cassandraErrors.ConversationReadStateBulkError, error)
	DeleteBulkConversationReadState(ctx context.Context, conversationIDs []string) error
}
