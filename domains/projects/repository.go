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

func (r *Repository) GetActiveByIdentifier(ctx context.Context, identifier string) (*Project, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id,
			organization_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at
		FROM projects
		WHERE (id = $1 OR slug = $1)
			AND status = 'active'
		ORDER BY CASE WHEN id = $1 THEN 0 ELSE 1 END
		LIMIT 1
	`, identifier)

	var project Project
	if err := row.Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &project, nil
}

func (r *Repository) ActiveExists(ctx context.Context, identifier string) (bool, error) {
	project, err := r.GetActiveByIdentifier(ctx, identifier)
	if err != nil {
		return false, err
	}

	return project != nil, nil
}
