package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	ucSession usecase.IUserSessionService
}

func NewSessionHandler(ucSession usecase.IUserSessionService) *SessionHandler {
	return &SessionHandler{
		ucSession: ucSession,
	}
}
func (h *SessionHandler) HandlerCreateSessionLogin(c *gin.Context, req *req.UserSessionReq) *response.Response {
	data, err := h.ucSession.CreateSessionLogin(req)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
func (h *SessionHandler) HandlerGetAllUserSessions(c *gin.Context, req *req.UserSessionReq) *response.Response {

	data, err := h.ucSession.GetAllUserSessions(req)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
func (h *SessionHandler) HandlerGetUserSessionPast(c *gin.Context, req *req.UserSessionReq) *response.Response {
	data, err := h.ucSession.GetUserSessionPast(req)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus(""))
	}
	return data
}
