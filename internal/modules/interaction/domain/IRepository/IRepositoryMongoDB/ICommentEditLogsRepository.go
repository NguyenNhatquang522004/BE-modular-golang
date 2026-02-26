package IRepositoryMongoDB

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

type ICommentEditLogsRepository interface {
	CreateEditLog(ctx context.Context, editlog *entity.CommentEntityEditLog) error
	CreateBulkEditLogs(ctx context.Context, editLogs []*entity.CommentEntityEditLog) (int64, []*mongodbErrors.BulkError, error)
	GetEditLogByID(ctx context.Context, editLogID string) (*entity.CommentEntityEditLog, error)
	GetBulkEditLogsByID(ctx context.Context, editLogIDs []string) ([]*entity.CommentEntityEditLog, error)
	GetEditLogsByTargetID(ctx context.Context, targetID string) ([]*entity.CommentEntityEditLog, error)
	GetLatestEditLogByTargetID(ctx context.Context, targetID string) (*entity.CommentEntityEditLog, error)
	GetVersionBulkEditLogsByTargetIDs(ctx context.Context, targetIDs []string) ([]*entity.CommentEntityEditLog, error)
	UpdateEditLog(ctx context.Context, editlog *entity.CommentEntityEditLog) error
	UpdateBulkEditLogs(ctx context.Context, editLogs []*entity.CommentEntityEditLog) (int64, []*mongodbErrors.BulkError, error)
	DeleteEditLog(ctx context.Context, editLogID string) error
	DeleteEditLogsByTargetID(ctx context.Context, targetID string) error
	DeleteBulkEditLogs(ctx context.Context, editLogsIds []string) (int64, []*mongodbErrors.BulkError, error)
	PaginationEditLogs(ctx context.Context, targetID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetVersionBulkEditLogsByTargetID(ctx context.Context, targetID string) ([]*entity.CommentEntityEditLog, error)
}
