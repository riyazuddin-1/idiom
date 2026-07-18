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
	IdentityID   string
	ProjectID    string
}

func Start(ctx context.Context, repo *Repository, jwtSettings *jwt.JWTSettings, identity *identities.Identity, ip string, userAgent string) (*SessionTokens, error) {
	sessionID := crypto.GenerateID("sid_")
	accessToken, err := jwtSettings.CreateToken(jwt.CustomClaims{
		"sid": sessionID,
		"sub": identity.ID,
		"pid": identity.ProjectID,
	})
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
		ID:               sessionID,
		IdentityID:       identity.ID,
		RefreshTokenHash: refreshTokenHash,
		IP:               ip,
		UserAgent:        userAgent,
		ExpiresAt:        now.Add(30 * 24 * time.Hour),
	}

	if err := repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &SessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    session.ID,
		IdentityID:   identity.ID,
		ProjectID:    identity.ProjectID,
	}, nil
}

func Refresh(ctx context.Context, repo *Repository, jwtSettings *jwt.JWTSettings, refreshToken string) (*SessionTokens, error) {
	session, newRefreshToken, projectID, err := repo.UpdateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := jwtSettings.CreateToken(jwt.CustomClaims{
		"sid": session.ID,
		"sub": session.IdentityID,
		"pid": projectID,
	})
	if err != nil {
		return nil, err
	}

	return &SessionTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		SessionID:    session.ID,
		IdentityID:   session.IdentityID,
		ProjectID:    projectID,
	}, nil
}

func Revoke(ctx context.Context, repo *Repository, sessionID string) (bool, error) {
	return repo.RevokeSession(ctx, sessionID)
}

func RevokeAll(ctx context.Context, repo *Repository, identityID string) (int64, error) {
	return repo.RevokeAllSessions(ctx, identityID)
}

func Validate(ctx context.Context, repo *Repository) {

}
