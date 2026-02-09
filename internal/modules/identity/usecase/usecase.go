package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type IUserRoleUseCase interface {
	AssignRoleToUserWithName(userID string, RoleName string) (*response.Response, error)
	AssignRoleToUserWithID(userID string, roleID string) (*response.Response, error)
	UpdateRoleOfUser(userID string, roleIDs []string) (*response.Response, error)
	GetUserRoles(userID string, roleName string) (*response.Response, error)
	UpdateRoleDescription(roleID string, description string) (*response.Response, error)
	CreateRole(role enum.RoleType, description string) (*response.Response, error)
	DeleteRole(roleID string) (*response.Response, error)
	GetAllUserRoles(RoleID string) (*response.Response, error)
	GetAllRoles() (*response.Response, error)
}
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
	LogOut(userID string, accessToken string) (*response.Response, error)
	ResetPassword(email string, newPassword string) (*response.Response, error)
	LoginWithGoogle(provider string, token string) (*response.Response, error)
}
type IUserSettingService interface {
	CreateUserSetting(userID string) (*response.Response, error)
	UpdateUserSettings(userID string, settings *entity.UserSetting) (*response.Response, error)
}

type IUserSessionService interface {
	CreateSessionLogin(session *entity.UserSession) (*response.Response, error)
	GetAllUserSessions(userID string) (*response.Response, error)
	GetUserSessionPast(userID string) (*response.Response, error)
}

type Usecase struct {
	Auth        IUserAuthService
	Google      IGoogleAuthUseCase
	UserSetting IUserSettingService
	Session     IUserSessionService
	Role        IUserRoleUseCase
}

func NewUsecase(
	auth IUserAuthService,
	google IGoogleAuthUseCase,
	userSetting IUserSettingService,
	session IUserSessionService,
	role IUserRoleUseCase,
) *Usecase {
	return &Usecase{
		Auth:        auth,
		Google:      google,
		UserSetting: userSetting,
		Session:     session,
		Role:        role,
	}
}
