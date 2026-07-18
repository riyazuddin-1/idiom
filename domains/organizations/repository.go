package organizations

import (
	"context"

	"idiom-api-services/packages/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByIdentityID(ctx context.Context, identityID string) ([]Organization, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			o.id,
			o.name,
			o.slug,
			o.status,
			o.created_at,
			o.updated_at
		FROM organizations o
		INNER JOIN org_members om ON om.organization_id = o.id
		WHERE om.identity_id = $1
			AND om.status = 'active'
		ORDER BY o.created_at DESC
	`, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	organizations := make([]Organization, 0)
	for rows.Next() {
		var organization Organization
		if err := rows.Scan(
			&organization.ID,
			&organization.Name,
			&organization.Slug,
			&organization.Status,
			&organization.CreatedAt,
			&organization.UpdatedAt,
		); err != nil {
			return nil, err
		}
		organizations = append(organizations, organization)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizations, nil
}
