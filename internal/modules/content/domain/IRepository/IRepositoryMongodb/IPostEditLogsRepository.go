package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostEditLogsRepository interface {
	CreatePostEditLog(ctx context.Context, postEditLog *entity.PostEntityEditLog) error
	CreateBulkPostEditLog(ctx context.Context, postEditLogs []*entity.PostEntityEditLog) (int64, []*mongodbErrors.BulkError, error)
	GetByTargetID(ctx context.Context, targetID string) (*entity.PostEntityEditLog, error)
	GetBulkByTargetID(ctx context.Context, targetIDs []string) ([]*entity.PostEntityEditLog, error)
	GetByID(ctx context.Context, id string) (*entity.PostEntityEditLog, error)
	GetBulkByID(ctx context.Context, ids []string) ([]*entity.PostEntityEditLog, error)
	GetLatestVersion(ctx context.Context, targetID string) (int, error)
	GetAllVersions(ctx context.Context, targetID string) ([]*entity.PostEntityEditLog, error)
	UpdatePostEditLog(ctx context.Context, postEditLog *entity.PostEntityEditLog) error
	UpdateBulkPostEditLog(ctx context.Context, postEditLogs []*entity.PostEntityEditLog) (int64, []*mongodbErrors.BulkError, error)
	PaginationPostEditLog(ctx context.Context, targetID string, cursor string, limit int) (*dto.PaginationRes, error)
}
