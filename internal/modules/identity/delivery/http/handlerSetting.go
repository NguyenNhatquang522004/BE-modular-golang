package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	ucSetting usecase.IUserSettingService
}

func NewSettingHandler(ucSetting usecase.IUserSettingService) *SettingHandler {
	return &SettingHandler{
		ucSetting: ucSetting,
	}
}
func (h *SettingHandler) HandlerCreateUserSetting(c *gin.Context, req *req.CreateUserSettingReq) *response.Response {
	data, err := h.ucSetting.CreateUserSetting(req.User_ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
func (h *SettingHandler) HandlerUpdateUserSettings(c *gin.Context, req *req.UserSettingReq) *response.Response {
	data, err := h.ucSetting.UpdateUserSettings(c.GetString("user_id"), req)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
