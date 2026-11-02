package sessions

import (
	"context"
	"idiom-api-services/domains/identities"
	"idiom-api-services/packages/crypto"
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

type SessionToken struct {
	SessionID  string
	IdentityID string
	ProjectID  string
}

func Create(ctx context.Context, repo *Repository, identity *identities.Identity, ip string, userAgent string) (*Session, error) {
	sessionID := crypto.GenerateID("sid_")

	now := time.Now().UTC()
	session := &Session{
		ID:         sessionID,
		IdentityID: identity.ID,
		IP:         ip,
		UserAgent:  userAgent,
		ExpiresAt:  now.Add(30 * 24 * time.Hour),
	}

	if err := repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func IssueToken(ctx context.Context, repo *Repository, sessionID, identityID, projectID, verificationSecret string) (string, error) {
	return crypto.EncryptJSON(SessionToken{
		SessionID:  sessionID,
		IdentityID: identityID,
		ProjectID:  projectID,
	}, verificationSecret)
}

func Revoke(ctx context.Context, repo *Repository, sessionID string) (bool, error) {
	return repo.RevokeSession(ctx, sessionID)
}

func RevokeAll(ctx context.Context, repo *Repository, identityID string) (int64, error) {
	return repo.RevokeAllSessions(ctx, identityID)
}

func Validate(ctx context.Context, repo *Repository) {

}
