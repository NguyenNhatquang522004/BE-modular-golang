package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ucUser usecase.IUserService
}

func NewUserHandler(ucUser usecase.IUserService) *UserHandler {
	return &UserHandler{
		ucUser: ucUser,
	}
}

func (h *UserHandler) HandlerGetUserByID(c *gin.Context, req *req.GetUserByIDReq) *response.Response {
	data, err := h.ucUser.GetUserByID(req.UserID)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data

}
func (h *UserHandler) HandlerGetUserByEmail(c *gin.Context, req *req.GetUserByEmailReq) *response.Response {
	data, err := h.ucUser.GetUserByEmail(req.Email)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

func (h *UserHandler) HandlerPanigationUsers(c *gin.Context, req *req.PanigationUsersReq) *response.Response {
	data, err := h.ucUser.PanigationUsers(req.Cursor, req.Limit)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

