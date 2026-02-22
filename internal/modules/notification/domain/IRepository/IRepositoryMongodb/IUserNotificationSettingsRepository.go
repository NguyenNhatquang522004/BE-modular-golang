package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
)

type IUserNotificationSettingsRepository interface {
	// GetByUserID lấy cài đặt thông báo của user theo user_id
	CreateUserNotificationSettings(ctx context.Context, settings *entity.UserNotificationSetting) error
	CreateBulkUserNotificationSettings(ctx context.Context, settings []*entity.UserNotificationSetting) (int64, []*mongodbErrors.BulkError, error)
	GetUserNotificationSettingsByUserID(ctx context.Context, userID string) (*entity.UserNotificationSetting, error)
	UpdateUserNotificationSettings(ctx context.Context, settings *entity.UserNotificationSetting) error
	UpdateBulkUserNotificationSettings(ctx context.Context, settings []*entity.UserNotificationSetting) (int64, []*mongodbErrors.BulkError, error)
	DeleteUserNotificationSettings(ctx context.Context, userID string) error
	DeleteBulkUserNotificationSettings(ctx context.Context, userIDs []string) (int64, []*mongodbErrors.BulkError, error)
}
