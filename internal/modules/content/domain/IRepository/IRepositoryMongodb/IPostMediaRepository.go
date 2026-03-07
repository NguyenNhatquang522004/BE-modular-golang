package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostMediaRepository interface {
	CreatePostMedia(ctx context.Context, postMedia *entity.PostMedia) error
	CreateBulkPostMedia(ctx context.Context, postMedias []*entity.PostMedia) (int64, []*mongodbErrors.BulkError, error)
	GetByPostID(ctx context.Context, postID string) (*entity.PostMedia, error)
	GetsByPostID(ctx context.Context, postID string) ([]*entity.PostMedia, error)
	GetBulkByPostIDs(ctx context.Context, postIDs []string) ([]*entity.PostMedia, error)
	UpdatePostMedia(ctx context.Context, postMedia *entity.PostMedia) error
	UpdateBulkPostMedia(ctx context.Context, postMedias []*entity.PostMedia) (int64, []*mongodbErrors.BulkError, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteBulkByIDs(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
	DeleteByPostID(ctx context.Context, postID string) error
	DeleteBulkByPostIDs(ctx context.Context, postIDs []string) (int64, []*mongodbErrors.BulkError, error)
	PanigationPostMedia(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error)
	PanigationPostsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
}
