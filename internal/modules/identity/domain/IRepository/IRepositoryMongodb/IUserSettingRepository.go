package IRepositoryMongodb

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserSettingRepository interface {
	GetUserSettings(userID string) (*response.Response, error)
	UpdateUserSettings(userID string, settings *entity.UserSetting) (*response.Response, error)
	DeleteUserSettings(userID string) (*response.Response, error)
	CreateUserSettingDefault(userID string) (*response.Response, error)
}
