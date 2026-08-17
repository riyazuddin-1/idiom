package organizations

import (
	"context"
	"errors"

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

func (r *Repository) Create(
	ctx context.Context,
	organization *Organization,
	identityID string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO organizations (
			id,
			name,
			slug,
			status
		) VALUES (
			$1, $2, $3, $4
		)
	`,
		organization.ID,
		organization.Name,
		organization.Slug,
		organization.Status,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO org_members (
			organization_id,
			identity_id,
			status
		) VALUES (
			$1, $2, 'active'
		)
	`,
		organization.ID,
		identityID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(
	ctx context.Context,
	organizationID string,
	identityID string,
	name string,
	slug string,
) (*Organization, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE organizations o
		SET
			name = $3,
			slug = $4,
			updated_at = NOW()
		WHERE o.id = $1
			AND EXISTS (
				SELECT 1
				FROM org_members om
				WHERE om.organization_id = o.id
					AND om.identity_id = $2
					AND om.status = 'active'
			)
		RETURNING
			o.id,
			o.name,
			o.slug,
			o.status,
			o.created_at,
			o.updated_at
	`,
		organizationID,
		identityID,
		name,
		slug,
	)

	var organization Organization
	if err := row.Scan(
		&organization.ID,
		&organization.Name,
		&organization.Slug,
		&organization.Status,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &organization, nil
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	organizationID string,
	identityID string,
	status string,
) (*Organization, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE organizations o
		SET
			status = $3,
			updated_at = NOW()
		WHERE o.id = $1
			AND EXISTS (
				SELECT 1
				FROM org_members om
				WHERE om.organization_id = o.id
					AND om.identity_id = $2
					AND om.status = 'active'
			)
		RETURNING
			o.id,
			o.name,
			o.slug,
			o.status,
			o.created_at,
			o.updated_at
	`,
		organizationID,
		identityID,
		status,
	)

	var organization Organization
	if err := row.Scan(
		&organization.ID,
		&organization.Name,
		&organization.Slug,
		&organization.Status,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &organization, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	organizationID string,
	identityID string,
) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM organizations o
		WHERE o.id = $1
			AND EXISTS (
				SELECT 1
				FROM org_members om
				WHERE om.organization_id = o.id
					AND om.identity_id = $2
					AND om.status = 'active'
			)
	`,
		organizationID,
		identityID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("organization not found")
	}

	return nil
}
