package keycloak

import (
	"context"
	"errors"

	"iiitn-career-portal/internal/config"

	"github.com/coreos/go-oidc"
)

type Claims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
}

func VerifyAccessToken(cfg config.Config, rawToken string) (*Claims, error) {
	ctx := context.Background()

	// 1️⃣ Discover Keycloak realm (issuer + JWKS)
	provider, err := oidc.NewProvider(
		ctx,
		cfg.BaseURL+"/realms/"+cfg.Realm,
	)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Strict verifier:
	// - verifies signature
	// - verifies expiry
	// - verifies issuer
	// - verifies audience == portal-backend
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	// 3️⃣ Verify token
	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	// 4️⃣ Extract claims
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
