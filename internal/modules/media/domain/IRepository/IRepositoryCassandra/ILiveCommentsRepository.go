package IRepositoryCassandra

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type ILiveCommentsRepository interface {
	// Define methods for LiveCommentsRepository here
	CreateLiveComment(ctx context.Context, comment *entity.LiveComment) error
	CreateBulkLiveComments(ctx context.Context, comments []*entity.LiveComment) (int64, []*dto.LiveCommentBulkError, error)
	GetLiveCommentsByStreamID(ctx context.Context, streamID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetLiveCommentsByUserID(ctx context.Context, streamID string, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetALLliveCommentsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateLiveComment(ctx context.Context, comment *entity.LiveComment) error
	UpdateBulkLiveComments(ctx context.Context, comments []*entity.LiveComment) (int64, []*dto.LiveCommentBulkError, error)
	DeleteLiveComment(ctx context.Context, streamID string, commentID string, createdAt time.Time) error
	DeleteBulkLiveComments(ctx context.Context, streamID string, commentIDs []string, createdAt time.Time) (int64, []*dto.LiveCommentBulkError, error)
	DeleteLiveCommentsByTimeRange(ctx context.Context, streamID string, startTime time.Time, endTime time.Time) error
	DeleteBulkLiveCommentsByTimeRange(ctx context.Context, streamIDs []string, startTime time.Time, endTime time.Time) (int64, []*dto.LiveCommentBulkError, error)
}
