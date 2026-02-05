package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/keycloak"
	"github.com/golang-jwt/jwt/v5"
)

type GoogleAuthUseCase struct {
	userRepo       IRepositoryPostgres.IUserRepository
	keycloakClient *keycloak.KeycloakRepository
}

func NewGoogleAuthUseCase(userRepo IRepositoryPostgres.IUserRepository, keycloakClient *keycloak.KeycloakRepository) *GoogleAuthUseCase {
	return &GoogleAuthUseCase{
		userRepo:       userRepo,
		keycloakClient: keycloakClient,
	}
}
func (u *GoogleAuthUseCase) GetUserInfoByToken(claim *jwt.MapClaims) (*response.Response, error) {
	mapClaims := *claim
	email, _ := mapClaims["email"].(string)
	name, _ := mapClaims["name"].(string)
	familyName, _ := mapClaims["family_name"].(string)
	givenName, _ := mapClaims["given_name"].(string)
	picture, _ := mapClaims["picture"].(string)

	userInfo := map[string]interface{}{
		"email":   email,
		"name":    name,
		"given":   givenName,
		"family":  familyName,
		"picture": picture,
	}

	return response.NewResponse(response.WithData(userInfo), response.WithMessage("Get user info by token successful"), response.WithStatus("200")), nil
}
func (u *GoogleAuthUseCase) Login(provider string, token string) (*response.Response, error) {
	// provider: "google", "facebook"... (Phải khớp với Alias trong Keycloak)
	// token: Chuỗi ID Token mà Frontend gửi lên

	// 1. Gọi Keycloak để đổi Token
	tokenResult, err := u.keycloakClient.ExchangeExternalToken(context.Background(), provider, token)
	if err != nil {
		return nil, errors.New("failed to exchange token with keycloak: " + err.Error())
	}

	// 2. Decode token để lấy User ID (sub)
	claims, err := u.keycloakClient.DecodeAccessToken(tokenResult.AccessToken)
	if err != nil {
		return nil, err
	}
	mapClaims := *claims
	userdata, datausererr := u.GetUserInfoByToken(claims)
	if datausererr != nil {
		return nil, datausererr
	}
	sub, ok := mapClaims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token: missing sub claim")
	}
	// 3. Đồng bộ User vào DB Postgres (Giống hệt luồng Login thường)
	// Tìm xem user này có trong DB chưa
	user, err := u.userRepo.FindByKeycloakID(sub)
	if err != nil || user == nil {
		// Nếu chưa có -> Tạo mới user trong DB nội bộ (Auto Register)
		// Lưu ý: Lúc này password để trống, vì user này login bằng Google
		newUser := &entity.User{
			KeycloakID: sub,
			Email:      userdata.Data.(map[string]interface{})["email"].(string),
			IsActive:   true,
		}
		_, err = u.userRepo.CreateUser(newUser)
		if err != nil {
			return nil, err
		}
	}

	// 4. Trả về Token cho Frontend
	return response.NewResponse(response.WithData(&res.TokenResponse{
		AccessToken:  tokenResult.AccessToken,
		RefreshToken: tokenResult.RefreshToken,
		ExpiresIn:    tokenResult.ExpiresIn,
	}), response.WithMessage("Login with Google successful"), response.WithStatus("200")), nil
}
