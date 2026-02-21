package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"github.com/gocql/gocql"
)

type IReactionsRepository interface {
	CreateReaction(ctx context.Context, reaction *entity.EntityReaction) error
	CreateBulkReactions(ctx context.Context, reactions []*entity.EntityReaction) (int64, []*cassandraErrors.ReactionBulkError, error)
	GetReactionsByTargetID(ctx context.Context, targetID string) ([]*entity.EntityReaction, error)
	GetReactionsByUserID(ctx context.Context, userID gocql.UUID) ([]*entity.EntityReaction, error)
	DeleteReaction(ctx context.Context, targetID string, userID gocql.UUID) error
	DeleteBulkReactions(ctx context.Context, targetID string, userID gocql.UUID) (int64, []*cassandraErrors.ReactionBulkError, error)
	UpdateReaction(ctx context.Context, reaction *entity.EntityReaction) error
	UpdateBulkReactions(ctx context.Context, reactions []*entity.EntityReaction) (int64, []*cassandraErrors.ReactionBulkError, error)
	PaginateReactionsByTargetID(ctx context.Context, targetID string, cursor string, limit int) (*dto.PaginationRes, error)
}
