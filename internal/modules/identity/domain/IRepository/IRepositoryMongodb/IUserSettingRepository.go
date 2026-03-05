package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserSettingRepository interface {
	GetUserSettings(ctx context.Context, userID string) (*entity.UserSetting, error)
	UpdateUserSettings(ctx context.Context, userID string, settings *entity.UserSetting) error
	DeleteUserSettings(ctx context.Context, userID string) error
	CreateUserSettingDefault(ctx context.Context, userSetting *entity.UserSetting) error
}
