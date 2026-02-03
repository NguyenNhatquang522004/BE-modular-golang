package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"golang.org/x/oauth2"
)

type IGoogleAuthUseCase interface {
	Login(email string, password string) (*response.Response, error)
	// 1. Tạo đường dẫn để chuyển hướng người dùng đến trang đăng nhập Google
	GetLoginURL(state string) (*response.Response, error)

	// 2. Đổi Authorization Code lấy Access Token (xử lý khi Google gọi lại callback)
	ExchangeCodeForToken(code string) (*response.Response, error)

	// 3. Sử dụng Token để lấy thông tin người dùng (Email, Name, Avatar...)
	GetUserInfo(token *oauth2.Token) (*response.Response, error)
}

type GoogleAuthUseCase struct {
	oauthConfig *oauth2.Config
}
