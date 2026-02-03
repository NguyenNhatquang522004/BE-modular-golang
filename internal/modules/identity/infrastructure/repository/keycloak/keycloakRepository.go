package keycloak

import (
	"context"
	"errors"
	"log"

	"github.com/Nerzal/gocloak/v13"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

type IKeycloakRepository interface {
	// ==========================================
	// 1. NHÓM AUTHENTICATION (Đăng nhập/Đăng xuất/Token)
	// ==========================================

	// Luồng 1: Direct Grant (User nhập user/pass trên App của bạn)
	LoginWithPassword(ctx context.Context, username string, password string) (*gocloak.JWT, error)

	// Luồng 2: Auth Code Flow (Frontend redirect, Backend đổi code lấy token)
	ExchangeCodeForToken(ctx context.Context, code string, redirectURI string) (*gocloak.JWT, error)

	// Làm mới token khi AccessToken hết hạn
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)

	// Đăng xuất (Vô hiệu hóa Refresh Token)
	Logout(ctx context.Context, refreshToken string) error

	// Kiểm tra Token còn sống hay không (Introspection)
	IntrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error)

	// Giải mã Token offline để lấy claims (Sub, Email, Roles...)
	DecodeAccessToken(accessToken string) (*jwt.MapClaims, error)

	// ==========================================
	// 2. NHÓM USER MANAGEMENT (CRUD User)
	// ==========================================

	// Tạo user mới (Register)
	CreateUser(ctx context.Context, user *gocloak.User, password string) (string, error)

	// Cập nhật thông tin user (Tên, Email, Attributes...)
	UpdateUser(ctx context.Context, userID string, user *gocloak.User) error

	// Xóa user (Khi cần dọn dẹp hoặc ban vĩnh viễn)
	DeleteUser(ctx context.Context, userID string) error

	// Lấy thông tin chi tiết 1 user bằng ID (UUID)
	GetUserByID(ctx context.Context, userID string) (*gocloak.User, error)

	// Tìm user bằng Email (Dùng để check duplicate hoặc sync)
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)

	// Tìm user bằng Username
	GetUserByUsername(ctx context.Context, username string) (*gocloak.User, error)

	// ==========================================
	// 3. NHÓM PASSWORD & SECURITY
	// ==========================================

	// Admin set cứng mật khẩu mới cho user (Reset Password thủ công)
	SetPassword(ctx context.Context, userID string, password string, temporary bool) error

	// Gửi email yêu cầu user tự đổi mật khẩu (Forgot Password / Required Actions)
	SendUpdatePasswordEmail(ctx context.Context, userID string) error

	// ==========================================
	// 4. NHÓM ROLE MANAGEMENT (Phân quyền RBAC)
	// ==========================================

	// Lấy danh sách tất cả Role có trong Realm
	GetRealmRoles(ctx context.Context) ([]*gocloak.Role, error)

	// Gán quyền (Role) cho User (Ví dụ: gán làm Admin)
	AddRealmRoleToUser(ctx context.Context, userID string, roleName string) error

	// Gỡ quyền (Role) khỏi User
	DeleteRealmRoleFromUser(ctx context.Context, userID string, roleName string) error

	// Kiểm tra xem User đang có những quyền gì
	GetRealmRolesByUserID(ctx context.Context, userID string) ([]*gocloak.Role, error)
}

type KeycloakRepository struct {
	client       *gocloak.GoCloak
	ctx          context.Context
	clientId     string
	clientSecret string
	realm        string
	cfg          *configs.KeyCloakConfig
	userRepo     postgres.UserRepository
}

func NewKeycloakRepository(cfg *configs.KeyCloakConfig) *KeycloakRepository {
	return &KeycloakRepository{
		client:       gocloak.NewClient(cfg.KEYCLOAK_SERVER_URL),
		ctx:          context.Background(),
		clientId:     cfg.KEYCLOAK_CLIENT_ID,
		clientSecret: cfg.KEYCLOAK_CLIENT_SECRET,
		realm:        cfg.KEYCLOAK_REALM,
		cfg:          cfg,
	}

}
func (k *KeycloakRepository) LoginWithPassword(username string, password string) (*gocloak.JWT, error) {
	token, err := k.client.Login(k.ctx, k.clientId, k.clientSecret, k.realm, username, password)
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (k *KeycloakRepository) ExchangeCodeForToken(code string, redirectURI string) (*gocloak.JWT, error) {
	// 1. Định nghĩa loại Grant Type là "authorization_code"
	grantType := "authorization_code"

	// 2. Gọi hàm GetToken của thư viện gocloak
	// Lưu ý: Phải truyền đúng RedirectURI khớp với cái Frontend đã dùng
	token, err := k.client.GetToken(k.ctx, k.realm, gocloak.TokenOptions{
		ClientID:     &k.clientId,
		ClientSecret: &k.clientSecret,
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
	jwt, err := k.client.RefreshToken(ctx, refreshToken, k.clientId, k.clientSecret, k.realm)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}
	return jwt, nil

}
func (k *KeycloakRepository) Logout(ctx context.Context, refreshToken string) error {
	return k.client.Logout(ctx, refreshToken, k.clientId, k.clientSecret, k.realm)
}

func (k *KeycloakRepository) IntrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {

	return k.client.RetrospectToken(ctx, accessToken, k.clientId, k.clientSecret, k.realm)
}
func (k *KeycloakRepository) DecodeAccessToken(accessToken string) (*jwt.MapClaims, error) {
	// Gọi hàm có sẵn của thư viện gocloak
	// Hàm này sẽ parse token và verify signature (nếu cấu hình), ở đây ta chủ yếu cần parse data
	_, claims, err := k.client.DecodeAccessToken(k.ctx, accessToken, k.realm)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (k *KeycloakRepository) GetAdminToken() (*gocloak.JWT, error) {
	grantType := "client_credentials"
	return k.client.GetToken(k.ctx, k.realm, gocloak.TokenOptions{
		ClientID:     &k.clientId,
		ClientSecret: &k.clientSecret,
		GrantType:    &grantType,
	})
}

// Hàm tạo User mới trên Keycloak
func (k *KeycloakRepository) CreateUser(user *gocloak.User, password string) (string, error) {
	// 1. Lấy Admin Token trước
	token, err := k.GetAdminToken()
	if err != nil {
		return "", err
	}

	// 2. Tạo User
	// User này được enable luôn
	userId, err := k.client.CreateUser(k.ctx, token.AccessToken, k.realm, *user)
	if err != nil {
		return "", err
	}

	// 3. Set mật khẩu cho User vừa tạo
	err = k.client.SetPassword(k.ctx, token.AccessToken, userId, k.realm, password, false)
	if err != nil {
		// Nếu set pass lỗi thì nên xóa user vừa tạo để tránh rác (Optional)
		return "", err
	}

	return userId, nil
}

// Cập nhật thông tin user (Tên, Email, Attributes...)
func (k *KeycloakRepository) UpdateUser(ctx context.Context, userID string, user *gocloak.User) error {
	// 1. Lấy Admin Token trước
	token, err := k.GetAdminToken()
	if err != nil {
		return err
	}
	// 2. Cập nhật User
	err = k.client.UpdateUser(ctx, token.AccessToken, k.realm, *user)
	if err != nil {
		return err
	}
	return nil
}

// Xóa user (Khi cần dọn dẹp hoặc ban vĩnh viễn)
func (k *KeycloakRepository) DeleteUser(ctx context.Context, userID string) error {
	token, err := k.GetAdminToken()
	if err != nil {
		return err
	}
	return k.client.DeleteUser(ctx, token.AccessToken, k.realm, userID)
}

// Lấy thông tin chi tiết 1 user bằng ID (UUID)
func (k *KeycloakRepository) GetUserByID(ctx context.Context, userID string) (*gocloak.User, error) {
	token, err := k.GetAdminToken()
	if err != nil {
		return nil, err
	}
	return k.client.GetUserByID(ctx, token.AccessToken, k.realm, userID)
}

// Tìm user bằng Email (Dùng để check duplicate hoặc sync)
func (k *KeycloakRepository) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	token, err := k.GetAdminToken()
	if err != nil {
		return nil, err
	}
	params := gocloak.GetUsersParams{
		Email: &email,
		Exact: gocloak.BoolP(true), // Tìm chính xác
	}
	users, err := k.client.GetUsers(ctx, token.AccessToken, k.realm, params)
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
	token, err := k.GetAdminToken()
	if err != nil {
		return nil, err
	}
	params := gocloak.GetUsersParams{
		Username: &username,
		Exact:    gocloak.BoolP(true),
	}

	users, err := k.client.GetUsers(ctx, token.AccessToken, k.realm, params)
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
	token, err := k.GetAdminToken()
	if err != nil {
		return err
	}
	return k.client.SetPassword(ctx, token.AccessToken, userID, k.realm, password, temporary)
}

// Gửi email yêu cầu user tự đổi mật khẩu (Forgot Password / Required Actions)
func (k *KeycloakRepository) SendUpdatePasswordEmail(ctx context.Context, userID string) error {
	token, err := k.GetAdminToken()
	if err != nil {
		return err
	}
	actions := []string{"UPDATE_PASSWORD"}

	// v13 yêu cầu truyền struct gocloak.ExecuteActionsEmail
	params := gocloak.ExecuteActionsEmail{
		UserID:  &userID,
		Actions: &actions,
	}

	return k.client.ExecuteActionsEmail(ctx, token.AccessToken, k.realm, params)
}

// ==========================================
// 4. NHÓM ROLE MANAGEMENT (Phân quyền RBAC)
// ==========================================

// Lấy danh sách tất cả Role có trong Realm
func (k *KeycloakRepository) GetRealmRoles(ctx context.Context) ([]*gocloak.Role, error) {
	token, err := k.GetAdminToken()
	if err != nil {
		return nil, err
	}
	// GetRealmRoles không nhận params filter trong v13 nên truyền nil nếu không cần
	return k.client.GetRealmRoles(ctx, token.AccessToken, k.realm, gocloak.GetRoleParams{})
}

// 16. Add Role To User (SỬA LỖI: Phải tìm Role Object trước)
func (k *KeycloakRepository) AddRealmRoleToUser(ctx context.Context, userID string, roleName string) error {
	token, err := k.GetAdminToken()
	if err != nil {
		return err
	}

	// Bước 1: Tìm cái Role object dựa trên tên (string)
	role, err := k.client.GetRealmRole(ctx, token.AccessToken, k.realm, roleName)
	if err != nil {
		return err // Role không tồn tại
	}

	// Bước 2: Truyền Role object vào hàm Add
	return k.client.AddRealmRoleToUser(ctx, token.AccessToken, k.realm, userID, []gocloak.Role{*role})
}

// 17. Delete Role From User (SỬA LỖI: Tương tự như Add)
func (k *KeycloakRepository) DeleteRealmRoleFromUser(ctx context.Context, userID string, roleName string) error {
	token, err := k.GetAdminToken()
	if err != nil {
		return err
	}

	// Bước 1: Tìm Role object
	role, err := k.client.GetRealmRole(ctx, token.AccessToken, k.realm, roleName)
	if err != nil {
		return err
	}

	// Bước 2: Xóa
	return k.client.DeleteRealmRoleFromUser(ctx, token.AccessToken, k.realm, userID, []gocloak.Role{*role})
}

// Kiểm tra xem User đang có những quyền gì
func (k *KeycloakRepository) GetRealmRolesByUserID(ctx context.Context, userID string) ([]*gocloak.Role, error) {
	token, err := k.GetAdminToken()
	if err != nil {
		return nil, err
	}
	return k.client.GetRealmRolesByUserID(ctx, token.AccessToken, k.realm, userID)
}
