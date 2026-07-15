package sessions

import (
	"context"
	"idiom-api-services/packages/crypto"
	"time"
)

func Start(ctx context.Context, repo *Repository, identityID string, ip string, userAgent string) (string, error) {
	refresh_token, err := crypto.RandomString(32)
	if err != nil {
		return "", err
	}
	refresh_token_hash, err := crypto.HashString(refresh_token)

	if err != nil {
		return "", nil
	}

	repo.Create(ctx, &Session{
		ID:               crypto.GenerateID("session", 16),
		IdentityID:       identityID,
		RefreshTokenHash: refresh_token_hash,
		IP:               ip,
		UserAgent:        userAgent,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		RevokedAt:        nil,
	})

	return refresh_token, nil
}
