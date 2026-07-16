package sessions

import (
	"context"
	"idiom-api-services/domains/identities"
	"idiom-api-services/packages/crypto"
	"idiom-api-services/packages/jwt"
	"time"
)

const (
	refreshTokenLength = 32
)

type SessionTokens struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
}

func Start(ctx context.Context, repo *Repository, jwtSettings *jwt.JWTSettings, identity *identities.Identity, ip string, userAgent string) (*SessionTokens, error) {
	accessToken, err := jwtSettings.CreateToken(identity.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := crypto.RandomString(refreshTokenLength)
	if err != nil {
		return nil, err
	}
	refreshTokenHash, err := crypto.HashString(refreshToken)

	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	session := &Session{
		ID:               crypto.GenerateID("sid_"),
		IdentityID:       identity.ID,
		RefreshTokenHash: refreshTokenHash,
		IP:               ip,
		UserAgent:        userAgent,
		ExpiresAt:        now.Add(30 * 24 * time.Hour),
		CreatedAt:        now,
		UpdatedAt:        now,
		RevokedAt:        nil,
	}

	if err := repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &SessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    session.ID,
	}, nil
}

func Refresh(ctx context.Context, repo *Repository, jwtSettings *jwt.JWTSettings, sessionID string) {}

func Revoke() {}

func RevokeAll() {}
