package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	ucRole usecase.IUserRoleUseCase
}

func NewRoleHandler(ucRole usecase.IUserRoleUseCase) *RoleHandler {
	return &RoleHandler{
		ucRole: ucRole,
	}
}
func (h *RoleHandler) HandlerAssignRoleToUserWithName(c *gin.Context, req *req.AssignRoleToUserWithNameReq) *response.Response {
	data, err := h.ucRole.AssignRoleToUserWithName(c.Request.Context(), c.GetString("user_id"), req.RoleName)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerAssignRoleToUserWithID(c *gin.Context, req *req.AssignRoleToUserWithIDReq) *response.Response {
	data, err := h.ucRole.AssignRoleToUserWithID(c.Request.Context(), c.GetString("user_id"), req.RoleID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
func (h *RoleHandler) HandlerUpdateRoleOfUser(c *gin.Context, req *req.UpdateRoleofUserReq) *response.Response {
	data, err := h.ucRole.UpdateRoleOfUser(c.Request.Context(), req.UserID, req.RoleID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerGetUserRoles(c *gin.Context, req *req.GetUserRolesReq) *response.Response {
	data, err := h.ucRole.GetUserRoles(c.Request.Context(), c.GetString("user_id"), req.RoleName)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerUpdateRoleDescription(c *gin.Context, req *req.UpdateRoleDescriptionReq) *response.Response {
	data, err := h.ucRole.UpdateRoleDescription(c.Request.Context(), req.RoleID, req.Description)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerCreateRole(c *gin.Context, req *req.CreateRoleReq) *response.Response {
	data, err := h.ucRole.CreateRole(c.Request.Context(), req.Role, req.Description)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerDeleteRole(c *gin.Context, req *req.DeleteRoleReq) *response.Response {
	data, err := h.ucRole.DeleteRole(c.Request.Context(), req.RoleID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerGetAllUserRoles(c *gin.Context, req *req.GetAllRoleReq) *response.Response {
	data, err := h.ucRole.GetAllUserRoles(c.Request.Context(), req.RoleID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}

func (h *RoleHandler) HandlerGetAllRoles(c *gin.Context) *response.Response {
	data, err := h.ucRole.GetAllRoles(c.Request.Context())
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
