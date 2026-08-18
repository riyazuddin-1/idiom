package sessions

import (
	"context"
	"fmt"
	"idiom-api-services/packages/crypto"
	"idiom-api-services/packages/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBySessionID(ctx context.Context, sessionID string) (*Session, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id,
			identity_id,
			refresh_token_hash,
			ip_address,
			user_agent,
			expires_at,
			created_at,
			updated_at,
			revoked_at
		FROM sessions
		WHERE id = $1
	`, sessionID)

	var session Session

	if err := row.Scan(
		&session.ID,
		&session.IdentityID,
		&session.RefreshTokenHash,
		&session.IP,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.RevokedAt,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *Repository) Create(ctx context.Context, session *Session) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions (
			id,
			identity_id,
			refresh_token_hash,
			ip_address,
			user_agent,
			expires_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6
		)
	`,
		session.ID,
		session.IdentityID,
		session.RefreshTokenHash,
		session.IP,
		session.UserAgent,
		session.ExpiresAt,
	)
	return err
}

func (r *Repository) UpdateRefreshToken(ctx context.Context, refreshToken string) (*Session, string, string, error) {
	refreshTokenHash, err := crypto.HashString(refreshToken)
	if err != nil {
		return nil, "", "", err
	}

	newRefreshToken, err := crypto.RandomString(refreshTokenLength)
	if err != nil {
		return nil, "", "", err
	}

	newRefreshTokenHash, err := crypto.HashString(newRefreshToken)
	if err != nil {
		return nil, "", "", err
	}

	row := r.db.QueryRow(ctx, `
		UPDATE sessions
		SET refresh_token_hash = $2,
			updated_at = NOW()
		FROM identities
		WHERE sessions.identity_id = identities.id
			AND sessions.refresh_token_hash = $1
			AND sessions.revoked_at IS NULL
			AND sessions.expires_at > NOW()
		RETURNING
			sessions.id,
			sessions.identity_id,
			sessions.refresh_token_hash,
			sessions.ip_address,
			sessions.user_agent,
			sessions.expires_at,
			sessions.created_at,
			sessions.updated_at,
			sessions.revoked_at,
			identities.project_id
	`, refreshTokenHash, newRefreshTokenHash)

	var session Session
	var projectID string
	if err := row.Scan(
		&session.ID,
		&session.IdentityID,
		&session.RefreshTokenHash,
		&session.IP,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.RevokedAt,
		&projectID,
	); err != nil {
		return nil, "", "", err
	}

	return &session, newRefreshToken, projectID, nil
}

func (r *Repository) UpdateRefreshTokenBySessionID(ctx context.Context, sessionID, refreshTokenHash string) (*Session, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE sessions
		SET refresh_token_hash = $2,
			updated_at = NOW()
		WHERE id = $1
			AND revoked_at IS NULL
			AND expires_at > NOW()
		RETURNING
			id,
			identity_id,
			refresh_token_hash,
			ip_address,
			user_agent,
			expires_at,
			created_at,
			updated_at,
			revoked_at
	`, sessionID, refreshTokenHash)

	var session Session
	if err := row.Scan(
		&session.ID,
		&session.IdentityID,
		&session.RefreshTokenHash,
		&session.IP,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.RevokedAt,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *Repository) RevokeSession(ctx context.Context, sessionID string) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		`
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1
			AND revoked_at IS NULL
		`,
		sessionID,
	)

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, fmt.Errorf("Session not found")
	}

	return true, nil
}

func (r *Repository) RevokeAllSessions(ctx context.Context, identityID string) (int64, error) {
	result, err := r.db.Exec(
		ctx,
		`
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE identity_id = $1
			AND revoked_at IS NULL
		`,
		identityID,
	)

	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, fmt.Errorf("Session not found")
	}

	return result.RowsAffected(), nil
}
