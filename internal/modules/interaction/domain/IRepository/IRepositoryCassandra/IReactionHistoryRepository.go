package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"github.com/gocql/gocql"
)

type IReactionHistoryRepository interface {
	CreateReactionHistory(ctx context.Context, reaction *entity.UserReactionHistory) error
	CreateBulkReactionHistory(ctx context.Context, reactions []*entity.UserReactionHistory) (int64, []*cassandraErrors.ReactionBulkError, error)
	UpdateReactionHistory(ctx context.Context, reaction *entity.UserReactionHistory) error
	GetReactionHistoryByUserID(ctx context.Context, userID gocql.UUID) ([]*entity.UserReactionHistory, error)
	GetReactionHistoryByUserIDAndTargetID(ctx context.Context, userID gocql.UUID, targetID gocql.UUID) (*entity.UserReactionHistory, error)
	PanigationReactionHistoryByUserID(ctx context.Context, userID gocql.UUID, cursor string, limit int) (*dto.PaginationRes, error)
}
