package projects

import (
	"context"

	"idiom-api-services/packages/database/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ActiveExists(ctx context.Context, projectID string) (bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id
		FROM projects
		WHERE id = $1
			AND status = 'active'
	`, projectID)

	var id string
	if err := row.Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
