package IRepositoryCassandra

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
)

type INotificationsRepository interface {
	CreateNotification(ctx context.Context, notification *entity.Notification) error
	CreateBulkNotifications(ctx context.Context, notifications []*entity.Notification) ([]*entity.Notification, []*cassandraErrors.NotificationBulkError, error)
	GetNotificationsByUserID(ctx context.Context, userID string, createdAt time.Time, notificationID string) ([]*entity.Notification, error)
	UpdateNotification(ctx context.Context, notification *entity.Notification) error
	UpdateBulkNotifications(ctx context.Context, notifications []*entity.Notification) ([]*entity.Notification, []*cassandraErrors.NotificationBulkError, error)
	DeleteNotification(ctx context.Context, userID string, createdAt time.Time, notificationID string) error
	DeleteBulkNotifications(ctx context.Context, notifications []*entity.Notification) ([]*entity.Notification, []*cassandraErrors.NotificationBulkError, error)
}
