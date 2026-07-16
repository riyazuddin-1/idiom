package sessions

import (
	"context"
	"idiom-api-services/packages/crypto"
	"time"
)

func Start(ctx context.Context, repo *Repository, identityID string, ip string, userAgent string) (string, error) {
	refreshToken, err := crypto.RandomString(32)
	if err != nil {
		return "", err
	}
	refreshTokenHash, err := crypto.HashString(refreshToken)

	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	if err := repo.Create(ctx, &Session{
		ID:               crypto.GenerateID("sid_"),
		IdentityID:       identityID,
		RefreshTokenHash: refreshTokenHash,
		IP:               ip,
		UserAgent:        userAgent,
		ExpiresAt:        now.Add(30 * 24 * time.Hour),
		CreatedAt:        now,
		UpdatedAt:        now,
		RevokedAt:        nil,
	}); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func Refresh()

func Revoke()

func RevokeAll()
