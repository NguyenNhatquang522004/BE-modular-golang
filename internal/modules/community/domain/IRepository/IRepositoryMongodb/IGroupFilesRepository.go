package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
)

type IGroupFilesRepository interface {
	CreateGroupFile(ctx context.Context, file *entity.GroupFile) error
	CreateBulkGroupFiles(ctx context.Context, files []entity.GroupFile) (int64, []*dto.BulkError, error)
	GetGroupFileByID(ctx context.Context, id string) (*entity.GroupFile, error)
	GetGroupFilesByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetGroupFilesByGroupIDAndUploaderID(ctx context.Context, groupID string, UploaderID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateGroupFileDownloadCount(ctx context.Context, id string, newCount int) error
	DeleteGroupFile(ctx context.Context, id string) error
	deleteBulkGroupFiles(ctx context.Context, ids []string) (int64, []*dto.BulkError, error)
	DeleteGroupFilesByGroupID(ctx context.Context, groupID string) error
	DeleteBulkGroupFilesByGroupID(ctx context.Context, groupIDs []string) (int64, []*dto.BulkError, error)
	DeleteGroupFileByIDAndGroupID(ctx context.Context, id string, groupID string) error
	DeleteBulkGroupFileByIDAndGroupID(ctx context.Context, ids []string, groupID string) (int64, []*dto.BulkError, error)
	DeleteGroupFilesByGroupIDAndUploaderID(ctx context.Context, groupID string, uploaderID string) error
	DeleteBulkGroupFilesByGroupIDAndUploaderID(ctx context.Context, groupID string, uploaderIDs []string) (int64, []*dto.BulkError, error)
}
