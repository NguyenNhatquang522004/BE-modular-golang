package auth

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
)

type TokenValidator struct {
	verifier *oidc.IDTokenVerifier
}

func NewTokenValidator(ctx context.Context, keycloakUrl, realm, clientID string) (*TokenValidator, error) {
	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("%s/realms/%s", keycloakUrl, realm))
	if err != nil {
		return nil, err
	}

	// Config verifier để check Issuer, ClientID
	config := &oidc.Config{
		ClientID: clientID,
	}
	
	return &TokenValidator{
		verifier: provider.Verifier(config),
	}, nil
}

// Hàm này trả về Claims (thông tin user) nếu token hợp lệ
func (v *TokenValidator) Validate(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	return v.verifier.Verify(ctx, rawToken)
}