package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostExtensionRepository interface {
	CreatePostExtension(ctx context.Context, postExtension *entity.PostExtension) error
	CreateBulkPostExtension(ctx context.Context, postExtensions []*entity.PostExtension) (int64, []*mongodbErrors.BulkError, error)
	GetByPostID(ctx context.Context, postID string) (*entity.PostExtension, error)
	GetBulkByPostIDs(ctx context.Context, postIDs []string) ([]*entity.PostExtension, error)
	GetShareDataByPostID(ctx context.Context, postID string) (*entity.ShareData, error)
	GetBulkShareDataByPostIDs(ctx context.Context, postIDs []string) ([]*entity.ShareData, error)
	UpdatePostExtension(ctx context.Context, postExtension *entity.PostExtension) error
	UpdateBulkPostExtension(ctx context.Context, postExtensions []*entity.PostExtension) (int64, []*mongodbErrors.BulkError, error)
	DeleteByPostID(ctx context.Context, postID string) error
	DeleteBulkByPostIDs(ctx context.Context, postIDs []string) (int64, []*mongodbErrors.BulkError, error)
	PaginationPostExtension(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error)
}
