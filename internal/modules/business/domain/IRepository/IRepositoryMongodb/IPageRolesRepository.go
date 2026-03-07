package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IPageRolesRepository interface {
	// Define the methods for the PageRolesRepository interface here
	CreatePageRole(ctx context.Context, pageRole *entity.PageRole) error
	CreateBulkPageRoles(ctx context.Context, pageRoles []*entity.PageRole) (int64, []*mongodbErrors.BulkError, error)
	GetPageRoleByID(ctx context.Context, pageRoleID string) (*entity.PageRole, error)
	GetPageRolesByPageID(ctx context.Context, pageID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetPageRolesByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetPageRoleByPageIDAndUserID(ctx context.Context, pageID string, userID string) (*entity.PageRole, error)
	GetPageRolesByPageIDAndUserID(ctx context.Context, pageID string, userID string) ([]*entity.PageRole, error)
	UpdatePageRole(ctx context.Context, pageRole *entity.PageRole) error
	UpdateBulkPageRoles(ctx context.Context, pageRoles []*entity.PageRole) (int64, []*mongodbErrors.BulkError, error)
	DeletePageRole(ctx context.Context, pageRoleID string) error
	DeleteBulkPageRoles(ctx context.Context, pageRoleIDs []string) (int64, []*mongodbErrors.BulkError, error)
	DeletePageRoleByUserIDandPageID(ctx context.Context, userID string, pageID string) error
	DeleteBulkPageRolesByUserIDsAndPageID(ctx context.Context, userIDs []string, pageID string) (int64, []*mongodbErrors.BulkError, error)
	DeleteBulkPageRolesByPageIDsAndUserID(ctx context.Context, pageIDs []string, userID string) (int64, []*mongodbErrors.BulkError, error)
	DeletePageRoleByPageID(ctx context.Context, pageID string) error
}
