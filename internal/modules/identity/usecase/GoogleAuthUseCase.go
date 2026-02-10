package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryKeyCloak"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleAuthUseCase struct {
	userRepo       IRepositoryPostgres.IUserRepository
	keycloakClient IRepositoryKeyCloak.IKeycloakRepository
	cfg            *configs.Config
}

func NewGoogleAuthUseCase(userRepo IRepositoryPostgres.IUserRepository, keycloakClient IRepositoryKeyCloak.IKeycloakRepository, cfg *configs.Config) *GoogleAuthUseCase {
	return &GoogleAuthUseCase{
		userRepo:       userRepo,
		keycloakClient: keycloakClient,
		cfg:            cfg,
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
	ctx := context.Background()
	// 2. Decode token để lấy User ID (sub)
	claims, err := u.keycloakClient.DecodeAccessToken(ctx, tokenResult.AccessToken)
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
func (u *GoogleAuthUseCase) LoginStandard(code string) (*response.Response, error) {
	// 1. Gọi Keycloak để đổi "Code" lấy "Token"
	// Endpoint chuẩn: POST /realms/{realm}/protocol/openid-connect/token
	tokenResult, err := u.keycloakClient.ExchangeAuthCode(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange code with keycloak: " + err.Error())
	}

	ctx := context.Background()
	claims, err := u.keycloakClient.DecodeAccessToken(ctx, tokenResult.AccessToken)
	if err != nil {
		return nil, err
	}

	mapClaims := *claims
	sub, ok := mapClaims["sub"].(string) // ID của User trong Keycloak
	email := mapClaims["email"].(string) // Email từ Google đã được Keycloak sync về
	if !ok {

		return nil, errors.New("invalid token: missing sub claim")
	}
	// 3. Đồng bộ User vào DB Postgres (Giống hệt luồng cũ)
	user, err := u.userRepo.FindByKeycloakID(sub)
	if err != nil || user == nil {
		// Tạo user mới nếu chưa có
		newUser := &entity.User{
			KeycloakID: sub,
			Email:      email,
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
	}), response.WithMessage("Login with Google (SSO) successful"), response.WithStatus("200")), nil
}

func (u *GoogleAuthUseCase) GetLoginURL(redirectURI string) (*response.Response, error) {
	return response.NewResponse(response.WithData(map[string]any{
		"login_url": u.cfg.KeyCloak.KEYCLOAK_AUTH_URL + "?client_id=" + u.cfg.KeyCloak.KEYCLOAK_CLIENT_ID + "&response_type=code&scope=openid email&redirect_uri=http://myapp.com/callback&kc_idp_hint=google",
	}), response.WithMessage("Get Google login URL successful"), response.WithStatus("200")), nil
}
