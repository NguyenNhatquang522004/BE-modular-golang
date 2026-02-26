package keycloak

import (
	"context"
	"errors"
	"log"

	"github.com/Nerzal/gocloak/v13"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

type KeycloakRepository struct {
	Client       *gocloak.GoCloak
	ClientId     string
	ClientSecret string
	Realm        string
	cfg          *configs.KeyCloakConfig
	userRepo     postgres.UserRepository
}

func NewKeycloakRepository(cfg *configs.KeyCloakConfig) *KeycloakRepository {
	return &KeycloakRepository{
		Client:       gocloak.NewClient(cfg.KEYCLOAK_SERVER_URL),
		ClientId:     cfg.KEYCLOAK_CLIENT_ID,
		ClientSecret: cfg.KEYCLOAK_CLIENT_SECRET,
		Realm:        cfg.KEYCLOAK_REALM,
		cfg:          cfg,
	}

}

func (k *KeycloakRepository) LoginWithPassword(ctx context.Context, username string, password string) (*gocloak.JWT, error) {
	token, err := k.Client.Login(ctx, k.ClientId, k.ClientSecret, k.Realm, username, password)
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (k *KeycloakRepository) ExchangeCodeForToken(ctx context.Context, code string, redirectURI string) (*gocloak.JWT, error) {
	// 1. Định nghĩa loại Grant Type là "authorization_code"
	grantType := "authorization_code"

	// 2. Gọi hàm GetToken của thư viện gocloak
	// Lưu ý: Phải truyền đúng RedirectURI khớp với cái Frontend đã dùng
	token, err := k.Client.GetToken(ctx, k.Realm, gocloak.TokenOptions{
		ClientID:     &k.ClientId,
		ClientSecret: &k.ClientSecret,
		GrantType:    &grantType,
		Code:         &code,
		RedirectURI:  &redirectURI,
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
func (k *KeycloakRepository) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	jwt, err := k.Client.RefreshToken(ctx, refreshToken, k.ClientId, k.ClientSecret, k.Realm)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}
	return jwt, nil

}
func (k *KeycloakRepository) Logout(ctx context.Context, refreshToken string) error {
	return k.Client.Logout(ctx, refreshToken, k.ClientId, k.ClientSecret, k.Realm)
}

func (k *KeycloakRepository) IntrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {

	return k.Client.RetrospectToken(ctx, accessToken, k.ClientId, k.ClientSecret, k.Realm)
}
func (k *KeycloakRepository) DecodeAccessToken(ctx context.Context, accessToken string) (*jwt.MapClaims, error) {
	// Gọi hàm có sẵn của thư viện gocloak
	// Hàm này sẽ parse token và verify signature (nếu cấu hình), ở đây ta chủ yếu cần parse data
	_, claims, err := k.Client.DecodeAccessToken(ctx, accessToken, k.Realm)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (k *KeycloakRepository) GetAdminToken(ctx context.Context) (*gocloak.JWT, error) {
	grantType := "client_credentials"
	return k.Client.GetToken(ctx, k.Realm, gocloak.TokenOptions{
		ClientID:     &k.ClientId,
		ClientSecret: &k.ClientSecret,
		GrantType:    &grantType,
	})
}

// Hàm tạo User mới trên Keycloak
func (k *KeycloakRepository) CreateUser(ctx context.Context, user *gocloak.User, password string) (string, error) {
	// 1. Lấy Admin Token trước
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return "", err
	}

	// 2. Tạo User
	// User này được enable luôn
	userId, err := k.Client.CreateUser(ctx, token.AccessToken, k.Realm, *user)
	if err != nil {
		return "", err
	}

	// 3. Set mật khẩu cho User vừa tạo
	err = k.Client.SetPassword(ctx, token.AccessToken, userId, k.Realm, password, false)
	if err != nil {
		// Nếu set pass lỗi thì nên xóa user vừa tạo để tránh rác (Optional)
		return "", err
	}

	return userId, nil
}

// Cập nhật thông tin user (Tên, Email, Attributes...)
func (k *KeycloakRepository) UpdateUser(ctx context.Context, userID string, user *gocloak.User) error {
	// 1. Lấy Admin Token trước
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	// 2. Cập nhật User
	err = k.Client.UpdateUser(ctx, token.AccessToken, k.Realm, *user)
	if err != nil {
		return err
	}
	return nil
}

// Xóa user (Khi cần dọn dẹp hoặc ban vĩnh viễn)
func (k *KeycloakRepository) DeleteUser(ctx context.Context, userID string) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	return k.Client.DeleteUser(ctx, token.AccessToken, k.Realm, userID)
}

// Lấy thông tin chi tiết 1 user bằng ID (UUID)
func (k *KeycloakRepository) GetUserByID(ctx context.Context, userID string) (*gocloak.User, error) {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return nil, err
	}
	return k.Client.GetUserByID(ctx, token.AccessToken, k.Realm, userID)
}

// Tìm user bằng Email (Dùng để check duplicate hoặc sync)
func (k *KeycloakRepository) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return nil, err
	}
	params := gocloak.GetUsersParams{
		Email: &email,
		Exact: gocloak.BoolP(true), // Tìm chính xác
	}
	users, err := k.Client.GetUsers(ctx, token.AccessToken, k.Realm, params)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil // Không tìm thấy user
	}
	return users[0], nil
	//
}

// Tìm user bằng Username
func (k *KeycloakRepository) GetUserByUsername(ctx context.Context, username string) (*gocloak.User, error) {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return nil, err
	}
	params := gocloak.GetUsersParams{
		Username: &username,
		Exact:    gocloak.BoolP(true),
	}

	users, err := k.Client.GetUsers(ctx, token.AccessToken, k.Realm, params)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

// Admin set cứng mật khẩu mới cho user (Reset Password thủ công)
func (k *KeycloakRepository) SetPassword(ctx context.Context, userID string, password string, temporary bool) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	return k.Client.SetPassword(ctx, token.AccessToken, userID, k.Realm, password, temporary)
}

// Gửi email yêu cầu user tự đổi mật khẩu (Forgot Password / Required Actions)
func (k *KeycloakRepository) SendUpdatePasswordEmail(ctx context.Context, userID string) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	actions := []string{"UPDATE_PASSWORD"}

	// v13 yêu cầu truyền struct gocloak.ExecuteActionsEmail
	params := gocloak.ExecuteActionsEmail{
		UserID:  &userID,
		Actions: &actions,
	}

	return k.Client.ExecuteActionsEmail(ctx, token.AccessToken, k.Realm, params)
}

// ==========================================
// 4. NHÓM ROLE MANAGEMENT (Phân quyền RBAC)
// ==========================================
// update role
func (k *KeycloakRepository) UpdateRealmRole(ctx context.Context, rolename sharedEnums.RoleType, updatedRole gocloak.Role) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	return k.Client.UpdateRealmRole(ctx, token.AccessToken, k.Realm, rolename.String(), updatedRole)
}

// Lấy danh sách tất cả Role có trong Realm
func (k *KeycloakRepository) GetRealmRoles(ctx context.Context) ([]*gocloak.Role, error) {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return nil, err
	}
	// GetRealmRoles không nhận params filter trong v13 nên truyền nil nếu không cần
	return k.Client.GetRealmRoles(ctx, token.AccessToken, k.Realm, gocloak.GetRoleParams{})
}

// 16. Add Role To User (SỬA LỖI: Phải tìm Role Object trước)
func (k *KeycloakRepository) AddRealmRoleToUser(ctx context.Context, userID string, roleName sharedEnums.RoleType) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}

	// Bước 1: Tìm cái Role object dựa trên tên (string)
	role, err := k.Client.GetRealmRole(ctx, token.AccessToken, k.Realm, roleName.String())
	if err != nil {
		return err // Role không tồn tại
	}

	// Bước 2: Truyền Role object vào hàm Add
	return k.Client.AddRealmRoleToUser(ctx, token.AccessToken, k.Realm, userID, []gocloak.Role{*role})
}
func (k *KeycloakRepository) UpdateListRealmRole(ctx context.Context, userID string, roleName []sharedEnums.RoleType) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	for _, role := range roleName {
		roleObj, err := k.Client.GetRealmRole(ctx, token.AccessToken, k.Realm, role.String())
		if err != nil {
			return err // Role không tồn tại
		}
		err = k.Client.AddRealmRoleToUser(ctx, token.AccessToken, k.Realm, userID, []gocloak.Role{*roleObj})
		if err != nil {
			return err
		}
	}

	return nil
}

// 17. Delete Role From User (SỬA LỖI: Tương tự như Add)
func (k *KeycloakRepository) DeleteRealmRoleFromUser(ctx context.Context, userID string, roleName sharedEnums.RoleType) error {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return err
	}

	// Bước 1: Tìm Role object
	role, err := k.Client.GetRealmRole(ctx, token.AccessToken, k.Realm, roleName.String())
	if err != nil {
		return err
	}

	// Bước 2: Xóa
	return k.Client.DeleteRealmRoleFromUser(ctx, token.AccessToken, k.Realm, userID, []gocloak.Role{*role})
}

// Kiểm tra xem User đang có những quyền gì
func (k *KeycloakRepository) GetRealmRolesByUserID(ctx context.Context, userID string) ([]*gocloak.Role, error) {
	token, err := k.GetAdminToken(ctx)
	if err != nil {
		return nil, err
	}
	return k.Client.GetRealmRolesByUserID(ctx, token.AccessToken, k.Realm, userID)
}
func (k *KeycloakRepository) ExchangeExternalToken(ctx context.Context, issuer string, externalToken string) (*gocloak.JWT, error) {
	// 1. Grant Type
	grantType := "urn:ietf:params:oauth:grant-type:token-exchange"

	// 2. Subject Token Type (Google ID Token)
	subjectTokenType := "urn:ietf:params:oauth:token-type:id_token"

	// 3. Requested Token Type (Muốn nhận lại Access Token)
	requestedTokenType := "urn:ietf:params:oauth:token-type:access_token"

	// 4. Scope: Cực kỳ quan trọng để lấy Refresh Token và thông tin user
	// Mẹo: Nên dùng chuỗi string thay vì []string cho ổn định
	scope := "openid profile email offline_access"

	// // Gọi API
	token, err := k.Client.GetToken(ctx, k.Realm, gocloak.TokenOptions{
		ClientID:     &k.ClientId,
		ClientSecret: &k.ClientSecret,
		GrantType:    &grantType,
		// --- CÁC PHẦN BẠN BỊ THIẾU ---
		SubjectToken:       &externalToken,
		RequestedSubject:   &subjectTokenType,   // <--- QUAN TRỌNG: Phải truyền vào đây
		RequestedTokenType: &requestedTokenType, // <--- Truyền vào đây
		Scope:              &scope,              // <--- Truyền vào đây để lấy Refresh Token
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (k *KeycloakRepository) ExchangeAuthCode(ctx context.Context, code string) (*gocloak.JWT, error) {
	// 1. Định nghĩa loại Grant Type là "authorization_code"
	grantType := "authorization_code"

	// 2. Gọi hàm GetToken của thư viện gocloak
	// Lưu ý: Phải truyền đúng RedirectURI khớp với cái Frontend đã dùng
	token, err := k.Client.GetToken(ctx, k.Realm, gocloak.TokenOptions{
		ClientID:     &k.ClientId,
		ClientSecret: &k.ClientSecret,
		GrantType:    &grantType,
		Code:         &code,
		RedirectURI:  &k.cfg.REDIRECT_URI,
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
