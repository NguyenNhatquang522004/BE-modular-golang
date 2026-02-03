package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
)

type IUserSettingService interface {
	CreateDefaultSettings(userID string) (*response.Response, error)
	GetUserSettings(userID string) (*response.Response, error)
	UpdateUserSettings(userID string, settings map[string]interface{}) (*response.Response, error)
}

type IUserSessionService interface {
	CreateSession(userID string, deviceInfo string) (*response.Response, error)
	ValidateSession(sessionToken string) (*response.Response, error)
	RevokeSession(sessionToken string) (*response.Response, error)
}

type Usecase struct {
}
