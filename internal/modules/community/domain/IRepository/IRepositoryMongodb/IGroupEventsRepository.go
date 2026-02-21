package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
)

type IGroupEventsRepository interface {
	CreateGroupEvent(ctx context.Context, event *entity.GroupEvent) error
	CreateBulkGroupEvents(ctx context.Context, events []*entity.GroupEvent) (int64, []*dto.BulkError, error)
	GetGroupEventByID(ctx context.Context, id string) (*entity.GroupEvent, error)
	GetGroupEventsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateGroupEvent(ctx context.Context, event *entity.GroupEvent) error
	UpdateBulkGroupEvents(ctx context.Context, events []*entity.GroupEvent) (int64, []*dto.BulkError, error)
	UpdateGroupEventDownloadCount(ctx context.Context, id string, newCount int) error
	DeleteGroupEvent(ctx context.Context, id string) error
	DeleteBulkGroupEvents(ctx context.Context, ids []string) (int64, []*dto.BulkError, error)
	DeleteGroupEventsByGroupID(ctx context.Context, groupID string) error
	DeleteGroupEventsByCreatorIDAndGroupID(ctx context.Context, creatorID string, groupID string) error
}
