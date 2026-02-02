package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"golang.org/x/oauth2"
)

type IUserAuthService interface {
	Login(email string, password string) (*response.Response, error)
	Register(username string, password string, email string) (*response.Response, error)
	SendOTP(email string) (*response.Response, error)
	VerifyOTP(email string, otp string) (*response.Response, error)
	ResetPassword(email string, newPassword string) (*response.Response, error)
	ChangeIsactiveStatus(userID string, isActive bool) (*response.Response, error)
}
type IGoogleAuthService interface {
	// 1. Tạo đường dẫn để chuyển hướng người dùng đến trang đăng nhập Google
	GetLoginURL(state string) (*response.Response, error)

	// 2. Đổi Authorization Code lấy Access Token (xử lý khi Google gọi lại callback)
	ExchangeCodeForToken(code string) (*response.Response, error)

	// 3. Sử dụng Token để lấy thông tin người dùng (Email, Name, Avatar...)
	GetUserInfo(token *oauth2.Token) (*response.Response, error)
}

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
