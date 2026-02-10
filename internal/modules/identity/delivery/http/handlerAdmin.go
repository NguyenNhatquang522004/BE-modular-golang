package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	ucAdmin usecase.AdminUseCase
}

func NewAdminHandler(ucAdmin usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		ucAdmin: ucAdmin,
	}
}

func (h *AdminHandler) HandlerLoadAllUsers(c *gin.Context, req *req.PanigationUsersReq) *response.Response {
	data, err := h.ucAdmin.LoadAllUsers(req.Cursor, req.Limit)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error load all users"), response.WithStatus("500"))
	}
	return data
}
func (h *AdminHandler) HandlerGetUserByID(c *gin.Context, req *req.GetUserByIDReq) *response.Response {
	data, err := h.ucAdmin.GetUserByID(req.UserID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error get user by id"), response.WithStatus("500"))
	}
	return data
}
func (h *AdminHandler) HandlerUpdateUser(c *gin.Context, req *req.UpdateUserReq) *response.Response {
	data, err := h.ucAdmin.UpdateUser(req)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error update user"), response.WithStatus("500"))
	}
	return data
}
func (h *AdminHandler) HandlerDeleteUser(c *gin.Context, req *req.DeleteUserReq) *response.Response {
	data, err := h.ucAdmin.DeleteUser(req.UserID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error delete user"), response.WithStatus("500"))
	}
	return data
}
func (h *AdminHandler) HandlerCreateUser(c *gin.Context, req *req.AdminCreateUserReq) *response.Response {
	data, err := h.ucAdmin.CreateUser(req.RoleName, req.RoleID, &req.CreateUserReq)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error create user"), response.WithStatus("500"))
	}
	return data
}
func (h *AdminHandler) HandlerAssignRoleToUser(c *gin.Context, req *req.IDFusionRoleReq) *response.Response {
	data, err := h.ucAdmin.AssignRoleToUser(req.UserID, req.RoleName, req.RoleID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user"), response.WithStatus("500"))
	}
	return data
}
