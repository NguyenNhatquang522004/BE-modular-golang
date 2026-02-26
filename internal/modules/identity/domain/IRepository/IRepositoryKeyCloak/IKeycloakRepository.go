package IRepositoryKeyCloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
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
	DecodeAccessToken(ctx context.Context, accessToken string) (*jwt.MapClaims, error)
	ExchangeAuthCode(ctx context.Context, code string) (*gocloak.JWT, error)
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
	UpdateRealmRole(ctx context.Context, rolename sharedEnums.RoleType, updatedRole gocloak.Role) error
	// Lấy danh sách tất cả Role có trong Realm
	GetRealmRoles(ctx context.Context) ([]*gocloak.Role, error)

	// Gán quyền (Role) cho User (Ví dụ: gán làm Admin)
	AddRealmRoleToUser(ctx context.Context, userID string, roleName sharedEnums.RoleType) error
	UpdateListRealmRole(ctx context.Context, userID string, roleName []sharedEnums.RoleType) error
	// Gỡ quyền (Role) khỏi User
	DeleteRealmRoleFromUser(ctx context.Context, userID string, roleName sharedEnums.RoleType) error

	// Kiểm tra xem User đang có những quyền gì
	GetRealmRolesByUserID(ctx context.Context, userID string) ([]*gocloak.Role, error)
	// google
	ExchangeExternalToken(ctx context.Context, issuer string, externalToken string) (*gocloak.JWT, error)
}
