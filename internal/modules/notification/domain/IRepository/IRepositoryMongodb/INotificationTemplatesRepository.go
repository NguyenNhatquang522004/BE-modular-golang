package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
)

type INotificationTemplatesRepository interface {
	// GetByType lấy template thông báo theo type
	CreateTemplate(ctx context.Context, template *entity.NotificationTemplate) error
	CreateBulkTemplates(ctx context.Context, templates []*entity.NotificationTemplate) (int64, []*mongodbErrors.BulkError, error)
	GetTemplateByType(ctx context.Context, notificationType enum.NotificationType, cursor string, limit int) (*dto.PaginationRes, error)
	GetTemplateByID(ctx context.Context, id string) (*entity.NotificationTemplate, error)
	GetTemplate(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateTemplate(ctx context.Context, template *entity.NotificationTemplate) error
	UpdateBulkTemplates(ctx context.Context, templates []*entity.NotificationTemplate) (int64, []*mongodbErrors.BulkError, error)
	DeleteTemplate(ctx context.Context, id string) error
	DeleteBulkTemplates(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
}
