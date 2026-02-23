package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	CreateBulkPosts(ctx context.Context, posts []*entity.Post) (int64, []*mongodbErrors.BulkError, error)
	GetPostByID(ctx context.Context, postID string) (*entity.Post, error)
	GetPostsBulkByIDs(ctx context.Context, postIDs []string) ([]*entity.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]*entity.Post, error)
	UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	UpdateBulkPosts(ctx context.Context, posts []*entity.Post) (int64, []*mongodbErrors.BulkError, error)
	DeletePost(ctx context.Context, postID string) error
	DeleteBulkPosts(ctx context.Context, postIDs []string) (int64, []*mongodbErrors.BulkError, error)
	PanigationPosts(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
}
