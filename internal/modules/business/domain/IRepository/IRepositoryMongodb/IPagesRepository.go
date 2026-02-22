package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IPagesRepository interface {
	// Define the methods for the PagesRepository interface here
	CreatePage(ctx context.Context, page *entity.Page) error
	CreateBulkPages(ctx context.Context, pages []*entity.Page) (int64, []*mongodbErrors.BulkError, error)
	GetPageByID(ctx context.Context, pageID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdatePage(ctx context.Context, page *entity.Page) error
	UpdateBulkPages(ctx context.Context, pages []*entity.Page) (int64, []*mongodbErrors.BulkError, error)
	DeletePage(ctx context.Context, pageID string) error
	DeleteBulkPages(ctx context.Context, pageIDs []string) (int64, []*mongodbErrors.BulkError, error)
}
