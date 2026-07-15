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
	var refreshTokenHash, err = crypto.HashJSON(refreshToken)
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, `
		SELECT *
		FROM sessions
		WHERE refresh_token_hash = $1
	`, refreshTokenHash)

	var session Session

	if err := row.Scan(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *Repository) Create(ctx context.Context, session *Session) {

}
