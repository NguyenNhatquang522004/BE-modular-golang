package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ucUserAuth   usecase.IUserAuthService
	ucGoogleAuth usecase.IGoogleAuthUseCase
}

func NewAuthHandler(ucUserAuth usecase.IUserAuthService, ucGoogleAuth usecase.IGoogleAuthUseCase) *AuthHandler {
	return &AuthHandler{
		ucUserAuth:   ucUserAuth,
		ucGoogleAuth: ucGoogleAuth,
	}
}
func (h *AuthHandler) HandlerRegisterOne(c *gin.Context, req *req.RegisterOneRequest) *response.Response {
	data, err := h.ucUserAuth.RegisterOne(req.Email)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

func (h *AuthHandler) HandlerRegisterTwo(c *gin.Context, req *req.RegisterTwoRequest) *response.Response {
	data, err := h.ucUserAuth.RegisterTwo(req.Email, req.OTP)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

func (h *AuthHandler) HandlerRegisterThree(c *gin.Context, req *req.RegisterThreeRequest) *response.Response {
	data, err := h.ucUserAuth.RegisterThree(req.Email, req.Username, req.Password)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

func (h *AuthHandler) HandlerLogin(c *gin.Context, req *req.LoginRequest) *response.Response {
	data, err := h.ucUserAuth.Login(req.Email, req.Password)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	tokenData := data.Data.(*res.TokenResponse)
	if tokenData.RefreshToken == "" || tokenData.AccessToken == "" || tokenData.ExpiresIn == 0 {
		return response.NewResponse(response.WithData(nil), response.WithMessage("Internal Server Error: Failed to generate token data"), response.WithStatus("500"))
	}

	utils.SetSecureCookie(c, tokenData.RefreshToken, "RefreshToken", int(tokenData.ExpiresIn))
	utils.SetSecureCookie(c, tokenData.AccessToken, "AccessToken", int(tokenData.ExpiresIn))
	return response.NewResponse(response.WithData(tokenData), response.WithMessage("Login successful"), response.WithStatus("200"))
}

func (h *AuthHandler) HandlerReSendOTP(c *gin.Context, req *req.ReSendOTPRequest) *response.Response {

	data, err := h.ucUserAuth.ReSendOTP(req.Email)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

func (h *AuthHandler) HandlerSendLinkResetPassword(c *gin.Context, req *req.SendLinkResetPasswordRequest) *response.Response {

	data, err := h.ucUserAuth.SendLinkResetPassword(req.Email)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}
func (h *AuthHandler) HandlerResetPassword(c *gin.Context, req *req.ResetPasswordRequest) *response.Response {
	data, err := h.ucUserAuth.ResetPassword(req.Email, req.NewPassword)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}
func (h *AuthHandler) Handlerlogout(c *gin.Context, req *req.Logout) *response.Response {
	data, err := h.ucUserAuth.LogOut(req.UserID, req.AccessToken)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}

func (h *AuthHandler) HandlerLoginWithGoogle(c *gin.Context, req *req.GoogleLoginRequest) *response.Response {

	data, err := h.ucGoogleAuth.Login("google", req.Token)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}
func (h *AuthHandler) HandlerGetGoogleLoginURL(c *gin.Context, req *req.GoogleLoginURLRequest) *response.Response {
	data, err := h.ucGoogleAuth.GetLoginURL(req.RedirectURI)
	if err != nil {
		return response.NewResponse(response.WithData(nil), response.WithMessage(err.Error()), response.WithStatus("400"))
	}
	return data
}
