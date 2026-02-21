package IRepositoryMongoDB

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

type ICommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	CreateBulkComments(ctx context.Context, comments []*entity.Comment) (int64, []*mongodbErrors.BulkError, error)
	GetCommentByID(ctx context.Context, commentID string) (*entity.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*entity.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*entity.Comment, error)
	UpdateComment(ctx context.Context, comment *entity.Comment) error
	UpdateBulkComments(ctx context.Context, comments []*entity.Comment) (int64, []*mongodbErrors.BulkError, error)
	DeleteComment(ctx context.Context, commentID string) error
	DeleteBulkComments(ctx context.Context, commentIDs []string) (int64, []*mongodbErrors.BulkError, error)
	PaginationComments(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error)
}
