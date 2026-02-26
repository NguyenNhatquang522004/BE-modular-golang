package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/golang-jwt/jwt/v5"
)

type IUserRoleUseCase interface {
	AssignRoleToUserWithName(userID string, RoleName string) (*response.Response, error)
	AssignRoleToUserWithID(userID string, roleID string) (*response.Response, error)
	UpdateRoleOfUser(userID string, roleIDs []string) (*response.Response, error)
	GetUserRoles(userID string, roleName string) (*response.Response, error)
	UpdateRoleDescription(roleID string, description string) (*response.Response, error)
	CreateRole(role sharedEnums.RoleType, description string) (*response.Response, error)
	DeleteRole(roleID string) (*response.Response, error)
	GetAllUserRoles(RoleID string) (*response.Response, error)
	GetAllRoles() (*response.Response, error)
}
type IAdminUseCase interface {
	LoadAllUsers(Cursor string, Limit int) (*response.Response, error)
	GetUserByID(userID string) (*response.Response, error)
	UpdateUser(req *req.UpdateUserReq) (*response.Response, error)
	DeleteUser(userID string) (*response.Response, error)
	CreateUser(rolename []sharedEnums.RoleType, roleID []string, req *req.CreateUserReq) (*response.Response, error)
	AssignRoleToUser(userID string, rolename []sharedEnums.RoleType, roleID []string) (*response.Response, error)
}
type IGoogleAuthUseCase interface {
	//Token Exchange
	Login(provider string, token string) (*response.Response, error)
	GetUserInfoByToken(claim *jwt.MapClaims) (*response.Response, error)
	// Standard flow
	LoginStandard(code string) (*response.Response, error)
	GetLoginURL(redirectURI string) (*response.Response, error)
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
}
type IUserService interface {
	GetUserByID(userID string) (*response.Response, error)
	GetUserByEmail(email string) (*response.Response, error)
	PanigationUsers(Cursor string, Limit int) (*response.Response, error)
	UpdateUser(req *req.UpdateUserReq) (*response.Response, error)
	DeleteUser(userID string) (*response.Response, error)
	CreateUser(req *req.CreateUserReq) (*response.Response, error)
}
type IUserSettingService interface {
	CreateUserSetting(userID string) (*response.Response, error)
	UpdateUserSettings(userID string, settings *req.UserSettingReq) (*response.Response, error)
}

type IUserSessionService interface {
	CreateSessionLogin(session *req.UserSessionReq) (*response.Response, error)
	GetAllUserSessions(req *req.UserSessionReq) (*response.Response, error)
	GetUserSessionPast(req *req.UserSessionReq) (*response.Response, error)
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
