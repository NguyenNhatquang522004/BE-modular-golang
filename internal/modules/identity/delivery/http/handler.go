package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/gin-gonic/gin"
)

type IHandlerAuth interface {
	HandlerRegisterOne(c *gin.Context, req *req.RegisterOneRequest) *response.Response
	HandlerRegisterTwo(c *gin.Context, req *req.RegisterTwoRequest) *response.Response
	HandlerRegisterThree(c *gin.Context, req *req.RegisterThreeRequest) *response.Response
	HandlerLogin(c *gin.Context, req *req.LoginRequest) *response.Response
	HandlerReSendOTP(c *gin.Context, req *req.ReSendOTPRequest) *response.Response
	HandlerSendLinkResetPassword(c *gin.Context, req *req.SendLinkResetPasswordRequest) *response.Response
	HandlerResetPassword(c *gin.Context, req *req.ResetPasswordRequest) *response.Response
	Handlerlogout(c *gin.Context, req *req.Logout) *response.Response
	HandlerLoginWithGoogle(c *gin.Context, req *req.GoogleLoginRequest) *response.Response
	HandlerGetGoogleLoginURL(c *gin.Context, req *req.GoogleLoginURLRequest) *response.Response
}
type IHandlerUser interface {
	HandlerGetUserByID(c *gin.Context, req *req.GetUserByIDReq) *response.Response
	HandlerGetUserByEmail(c *gin.Context, req *req.GetUserByEmailReq) *response.Response
	HandlerPanigationUsers(c *gin.Context, req *req.PanigationUsersReq) *response.Response
}
type IHandlerRole interface {
	HandlerAssignRoleToUserWithName(c *gin.Context, req *req.AssignRoleToUserWithNameReq) *response.Response
	HandlerAssignRoleToUserWithID(c *gin.Context, req *req.AssignRoleToUserWithIDReq) *response.Response
	HandlerUpdateRoleOfUser(c *gin.Context, req *req.UpdateRoleofUserReq) *response.Response
	HandlerGetUserRoles(c *gin.Context, req *req.GetUserRolesReq) *response.Response
	HandlerUpdateRoleDescription(c *gin.Context, req *req.UpdateRoleDescriptionReq) *response.Response
	HandlerCreateRole(c *gin.Context, req *req.CreateRoleReq) *response.Response
	HandlerDeleteRole(c *gin.Context, req *req.DeleteRoleReq) *response.Response
	HandlerGetAllUserRoles(c *gin.Context, req *req.GetAllRoleReq) *response.Response
	HandlerGetAllRoles(c *gin.Context) *response.Response
}
type IHandlerUserSetting interface {
	HandlerCreateUserSetting(c *gin.Context, req *req.CreateUserSettingReq) *response.Response
	HandlerUpdateUserSettings(c *gin.Context, req *req.UserSettingReq) *response.Response
}
type IHandlerUserSession interface {
	HandlerCreateSessionLogin(c *gin.Context, req *req.UserSessionReq) *response.Response
	HandlerGetAllUserSessions(c *gin.Context, req *req.UserSessionReq) *response.Response
	HandlerGetUserSessionPast(c *gin.Context, req *req.UserSessionReq) *response.Response
}
type IHandleAdmin interface {
	HandlerLoadAllUsers(c *gin.Context, req *req.PanigationUsersReq) *response.Response
	HandlerGetUserByID(c *gin.Context, req *req.GetUserByIDReq) *response.Response
	HandlerUpdateUser(c *gin.Context, req *req.UpdateUserReq) *response.Response
	HandlerDeleteUser(c *gin.Context, req *req.DeleteUserReq) *response.Response
	HandlerCreateUser(c *gin.Context, req *req.AdminCreateUserReq) *response.Response
	HandlerAssignRoleToUser(c *gin.Context, req *req.IDFusionRoleReq) *response.Response
}

type Handler struct {
	Auth        IHandlerAuth
	User        IHandlerUser
	Role        IHandlerRole
	UserSetting IHandlerUserSetting
	UserSession IHandlerUserSession
	Admin       IHandleAdmin
}

func NewHandler(
	auth IHandlerAuth,
	user IHandlerUser,
	role IHandlerRole,
	userSetting IHandlerUserSetting,
	userSession IHandlerUserSession,

	admin IHandleAdmin,
) *Handler {
	return &Handler{
		Auth:        auth,
		User:        user,
		Role:        role,
		UserSetting: userSetting,
		UserSession: userSession,
		Admin:       admin,
	}
}
