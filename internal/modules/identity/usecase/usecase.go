package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type IGoogleAuthUseCase interface {
	//Token Exchange
	Login(provider string, token string) (*response.Response, error)
	GetUserInfoByToken(claim *jwt.MapClaims) (*response.Response, error)
	// Standard flow
	LoginStandard(code string) (*response.Response, error)
	GetLoginURL(redirectURI string) (*response.Response, error)
	ExchangeCodeForToken(code string) (*response.Response, error)
	GetUserInfo(token *oauth2.Token) (*response.Response, error)
}
type IUserAuthService interface {
	Login(email string, password string) (*response.Response, error)
	RegisterOne(email string) (*response.Response, error)
	RegisterTwo(email string, otp string) (*response.Response, error)
	RegisterThree(email string, username string, password string) (*response.Response, error)
	ReSendOTP(email string) (*response.Response, error)
	SendLinkResetPassword(email string) (*response.Response, error)
	LogOut(userID string) (*response.Response, error)
	ResetPassword(email string, newPassword string) (*response.Response, error)
	LoginWithGoogle(provider string, token string) (*response.Response, error)
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
