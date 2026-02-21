package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
)

type IGroupRepository interface {
	CreateGroup(ctx context.Context, group *entity.Group) error
	CreateBulkGroups(ctx context.Context, groups []*entity.Group) (int64, []*mongodbErrors.BulkError, error)
	GetGroupByID(ctx context.Context, id string) (*entity.Group, error)
	GetBulkGroupsByIDs(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error)
	GetGroupsByCreatorID(ctx context.Context, creatorID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateGroup(ctx context.Context, group *entity.Group) error
	UpdateBulkGroups(ctx context.Context, groups []*entity.Group) (int64, []*mongodbErrors.BulkError, error)
	DeleteGroup(ctx context.Context, id string) error
	DeleteBulkGroups(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
}
