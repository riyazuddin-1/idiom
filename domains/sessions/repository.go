package sessions

import (
	"context"
	"idiom-api-services/packages/crypto"
	"idiom-api-services/packages/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func (r *Repository) GetByRefreshToken(ctx context.Context, refreshToken string) (*Session, error) {
	refreshTokenHash, err := crypto.HashString(refreshToken)
	if err != nil {
		return nil, err
	}
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
		WHERE refresh_token_hash = $1
	`, refreshTokenHash)

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
			expires_at,
			created_at,
			updated_at,
			revoked_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
	`,
		session.ID,
		session.IdentityID,
		session.RefreshTokenHash,
		session.IP,
		session.UserAgent,
		session.ExpiresAt,
		session.CreatedAt,
		session.UpdatedAt,
		session.RevokedAt,
	)
	return err
}
