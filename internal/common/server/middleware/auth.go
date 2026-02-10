package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type KeycloakClaims struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	EmailVerified     bool   `json:"email_verified"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]interface{} `json:"resource_access"` // Client roles
}

// AuthMiddleware struct giữ state của provider để không phải init lại mỗi request
type AuthMiddleware struct {
	Verifier *oidc.IDTokenVerifier
}

// NewAuthMiddleware khởi tạo kết nối đến Keycloak để lấy JWKS (Public Keys)
func NewAuthMiddleware(ctx context.Context, issuerURL string, clientID string) (*AuthMiddleware, error) {
	// 1. Auto-discovery: Golang gọi Keycloak để lấy config
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	// 2. Cấu hình Verifier
	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true, // Bật cái này nếu validate Access Token thay vì ID Token
		// Access Token của Keycloak thường có 'aud' là 'account', không phải clientID
	}

	// Mẹo Best Practice: Với Access Token, ta thường skip check ClientID
	// và chỉ quan tâm Token đó có được ký bởi đúng Issuer hay không.

	return &AuthMiddleware{
		Verifier: provider.Verifier(oidcConfig),
	}, nil
}

// Handler là middleware function
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Lấy token từ Header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 2. Verify Token (Kiểm tra chữ ký, hạn dùng, issuer)
		idToken, err := am.Verifier.Verify(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		// 3. Extract Claims (Lấy thông tin user: sub, email, roles...)
		var claims KeycloakClaims

		// Thư viện tự động map JSON vào Struct (Gọn gàng!)
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// 4. Inject vào Context để Controller sử dụng
		ctx := context.WithValue(r.Context(), "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_roles", claims.RealmAccess.Roles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
