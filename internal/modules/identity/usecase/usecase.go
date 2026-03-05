package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/golang-jwt/jwt/v5"
)

type IUserRoleUseCase interface {
	AssignRoleToUserWithName(ctx context.Context, userID string, RoleName string) (*response.Response, error)
	AssignRoleToUserWithID(ctx context.Context, userID string, roleID string) (*response.Response, error)
	UpdateRoleOfUser(ctx context.Context, userID string, roleIDs []string) (*response.Response, error)
	GetUserRoles(ctx context.Context, userID string, roleName string) (*response.Response, error)
	UpdateRoleDescription(ctx context.Context, roleID string, description string) (*response.Response, error)
	CreateRole(ctx context.Context, role sharedEnums.RoleType, description string) (*response.Response, error)
	DeleteRole(ctx context.Context, roleID string) (*response.Response, error)
	GetAllUserRoles(ctx context.Context, RoleID string) (*response.Response, error)
	GetAllRoles(ctx context.Context) (*response.Response, error)
}
type IAdminUseCase interface {
	LoadAllUsers(ctx context.Context, Cursor string, Limit int) (*response.Response, error)
	GetUserByID(ctx context.Context, userID string) (*response.Response, error)
	UpdateUser(ctx context.Context, req *req.UpdateUserReq) (*response.Response, error)
	DeleteUser(ctx context.Context, userID string) (*response.Response, error)
	CreateUser(ctx context.Context, rolename []sharedEnums.RoleType, roleID []string, req *req.CreateUserReq) (*response.Response, error)
	AssignRoleToUser(ctx context.Context, userID string, rolename []sharedEnums.RoleType, roleID []string) (*response.Response, error)
}
type IGoogleAuthUseCase interface {
	//Token Exchange
	Login(ctx context.Context, provider string, token string) (*response.Response, error)
	GetUserInfoByToken(ctx context.Context, claim *jwt.MapClaims) (*response.Response, error)
	// Standard flow
	LoginStandard(ctx context.Context, code string) (*response.Response, error)
	GetLoginURL(ctx context.Context, redirectURI string) (*response.Response, error)
}
type IUserAuthService interface {
	Login(ctx context.Context, email string, password string) (*response.Response, error)
	RegisterOne(ctx context.Context, email string) (*response.Response, error)
	RegisterTwo(ctx context.Context, email string, otp string) (*response.Response, error)
	RegisterThree(ctx context.Context, email string, username string, password string) (*response.Response, error)
	ReSendOTP(ctx context.Context, email string) (*response.Response, error)
	SendLinkResetPassword(ctx context.Context, email string) (*response.Response, error)
	LogOut(ctx context.Context, userID string, accessToken string) (*response.Response, error)
	ResetPassword(ctx context.Context, email string, newPassword string) (*response.Response, error)
}
type IUserService interface {
	GetUserByID(ctx context.Context, userID string) (*response.Response, error)
	GetUserByEmail(ctx context.Context, email string) (*response.Response, error)
	PanigationUsers(ctx context.Context, Cursor string, Limit int) (*response.Response, error)
	UpdateUser(ctx context.Context, req *req.UpdateUserReq) (*response.Response, error)
	DeleteUser(ctx context.Context, userID string) (*response.Response, error)
	CreateUser(ctx context.Context, req *req.CreateUserReq) (*response.Response, error)
}
type IUserSettingService interface {
	CreateUserSetting(ctx context.Context, userID string) (*response.Response, error)
	UpdateUserSettings(ctx context.Context, userID string, settings *req.UserSettingReq) (*response.Response, error)
}

type IUserSessionService interface {
	CreateSessionLogin(ctx context.Context, session *req.UserSessionReq) (*response.Response, error)
	GetAllUserSessions(ctx context.Context, req *req.UserSessionReq) (*response.Response, error)
	GetUserSessionPast(ctx context.Context, req *req.UserSessionReq) (*response.Response, error)
}

type Usecase struct {
	Auth        IUserAuthService
	Google      IGoogleAuthUseCase
	UserSetting IUserSettingService
	Session     IUserSessionService
	Role        IUserRoleUseCase
	User        IUserService
	Admin       IAdminUseCase
}

func NewUsecase(
	auth IUserAuthService,
	google IGoogleAuthUseCase,
	userSetting IUserSettingService,
	session IUserSessionService,
	role IUserRoleUseCase,
	user IUserService,
	admin IAdminUseCase,

) *Usecase {
	return &Usecase{
		Auth:        auth,
		Google:      google,
		UserSetting: userSetting,
		Session:     session,
		Role:        role,
		User:        user,
		Admin:       admin,
	}
}
