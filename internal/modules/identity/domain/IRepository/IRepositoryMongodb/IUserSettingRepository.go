package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserSettingRepository interface {
	GetUserSettings(ctx context.Context, userID string) (*response.Response, error)
	UpdateUserSettings(ctx context.Context, userID string, settings *entity.UserSetting) (*response.Response, error)
	DeleteUserSettings(ctx context.Context, userID string) (*response.Response, error)
	CreateUserSettingDefault(ctx context.Context, userID string) (*response.Response, error)
}
