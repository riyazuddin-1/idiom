package loadtest

import (
	"context"
	"encoding/json"
	"sync"

	"idiom-api-services/packages/database/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db        *postgres.Postgres
	tableOnce sync.Once
	tableErr  error
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) EnsureTable(ctx context.Context) error {
	r.tableOnce.Do(func() {
		_, r.tableErr = r.db.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS loadtest_items (
				id TEXT PRIMARY KEY,
				payload JSONB NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)
		`)
	})
	return r.tableErr
}

func (r *Repository) Toggle(ctx context.Context, id string, payload json.RawMessage) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, id); err != nil {
		return "", err
	}

	var deletedID string
	err = tx.QueryRow(ctx, `
		DELETE FROM loadtest_items
		WHERE id = $1
		RETURNING id
	`, id).Scan(&deletedID)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return "", err
		}
		return "removed", nil
	}
	if err != pgx.ErrNoRows {
		return "", err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO loadtest_items (id, payload)
		VALUES ($1, $2::jsonb)
	`, id, payload); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return "added", nil
}
