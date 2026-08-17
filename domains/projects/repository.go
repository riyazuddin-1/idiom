package projects

import (
	"context"
	"encoding/json"

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
			COALESCE(description, ''),
			status,
			COALESCE(auth_options, ARRAY[]::text[]),
			COALESCE(redirect_urls, '[]'::jsonb),
			COALESCE(allowed_origins, ARRAY[]::text[]),
			created_at,
			updated_at
		FROM projects
		WHERE (id = $1 OR slug = $1)
			AND status = 'active'
		ORDER BY CASE WHEN id = $1 THEN 0 ELSE 1 END
		LIMIT 1
	`, identifier)

	var project Project
	var redirectURLs []byte

	if err := row.Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.AuthOptions,
		&redirectURLs,
		&project.AllowedOrigins,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if len(redirectURLs) > 0 {
		if err := json.Unmarshal(redirectURLs, &project.RedirectURLs); err != nil {
			return nil, err
		}
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

func (r *Repository) ListByOrganizationID(
	ctx context.Context,
	organizationID string,
	identityID string,
) ([]Project, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			p.id,
			p.organization_id,
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.status,
			COALESCE(p.auth_options, ARRAY[]::text[]),
			COALESCE(p.redirect_urls, '[]'::jsonb),
			COALESCE(p.allowed_origins, ARRAY[]::text[]),
			p.created_at,
			p.updated_at
		FROM projects p
		INNER JOIN org_members om
			ON om.organization_id = p.organization_id
		WHERE p.organization_id = $1
			AND om.identity_id = $2
			AND om.status = 'active'
		ORDER BY p.created_at DESC
	`, organizationID, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]Project, 0)

	for rows.Next() {
		var project Project
		var redirectURLs []byte

		if err := rows.Scan(
			&project.ID,
			&project.OrganizationID,
			&project.Name,
			&project.Slug,
			&project.Description,
			&project.Status,
			&project.AuthOptions,
			&redirectURLs,
			&project.AllowedOrigins,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if len(redirectURLs) > 0 {
			if err := json.Unmarshal(redirectURLs, &project.RedirectURLs); err != nil {
				return nil, err
			}
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	organizationID string,
	projectID string,
	identityID string,
) (*Project, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			p.id,
			p.organization_id,
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.status,
			COALESCE(p.auth_options, ARRAY[]::text[]),
			COALESCE(p.redirect_urls, '[]'::jsonb),
			COALESCE(p.allowed_origins, ARRAY[]::text[]),
			p.created_at,
			p.updated_at
		FROM projects p
		INNER JOIN org_members om
			ON om.organization_id = p.organization_id
		WHERE p.id = $1
			AND p.organization_id = $2
			AND om.identity_id = $3
			AND om.status = 'active'
	`, projectID, organizationID, identityID)

	var project Project
	var redirectURLs []byte

	if err := row.Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.AuthOptions,
		&redirectURLs,
		&project.AllowedOrigins,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if len(redirectURLs) > 0 {
		if err := json.Unmarshal(redirectURLs, &project.RedirectURLs); err != nil {
			return nil, err
		}
	}

	return &project, nil
}

func (r *Repository) Create(
	ctx context.Context,
	project *Project,
	identityID string,
) error {
	redirectURLs, err := json.Marshal(project.RedirectURLs)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	var exists bool

	if err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM org_members
			WHERE organization_id = $1
				AND identity_id = $2
				AND status = 'active'
		)
	`, project.OrganizationID, identityID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return ErrProjectNotFound
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO projects (
			id,
			organization_id,
			name,
			slug,
			description,
			status,
			auth_options,
			redirect_urls,
			allowed_origins
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
	`,
		project.ID,
		project.OrganizationID,
		project.Name,
		project.Slug,
		project.Description,
		project.Status,
		project.AuthOptions,
		redirectURLs,
		project.AllowedOrigins,
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
	projectID string,
	identityID string,
	name string,
	slug string,
	description string,
	authOptions []string,
	redirectURLsData []RedirectURL,
	allowedOrigins []string,
) (*Project, error) {
	redirectURLs, err := json.Marshal(redirectURLsData)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, `
		UPDATE projects p
		SET
			name = $4,
			slug = $5,
			description = $6,
			auth_options = $7,
			redirect_urls = $8,
			allowed_origins = $9,
			updated_at = NOW()
		WHERE p.id = $1
			AND p.organization_id = $2
			AND EXISTS (
				SELECT 1
				FROM org_members om
				WHERE om.organization_id = p.organization_id
					AND om.identity_id = $3
					AND om.status = 'active'
			)
		RETURNING
			p.id,
			p.organization_id,
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.status,
			COALESCE(p.auth_options, ARRAY[]::text[]),
			COALESCE(p.redirect_urls, '[]'::jsonb),
			COALESCE(p.allowed_origins, ARRAY[]::text[]),
			p.created_at,
			p.updated_at
	`,
		projectID,
		organizationID,
		identityID,
		name,
		slug,
		description,
		authOptions,
		redirectURLs,
		allowedOrigins,
	)

	return scanProject(row)
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	organizationID string,
	projectID string,
	identityID string,
	status string,
) (*Project, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE projects p
		SET
			status = $4,
			updated_at = NOW()
		WHERE p.id = $1
			AND p.organization_id = $2
			AND EXISTS (
				SELECT 1
				FROM org_members om
				WHERE om.organization_id = p.organization_id
					AND om.identity_id = $3
					AND om.status = 'active'
			)
		RETURNING
			p.id,
			p.organization_id,
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.status,
			COALESCE(p.auth_options, ARRAY[]::text[]),
			COALESCE(p.redirect_urls, '[]'::jsonb),
			COALESCE(p.allowed_origins, ARRAY[]::text[]),
			p.created_at,
			p.updated_at
	`,
		projectID,
		organizationID,
		identityID,
		status,
	)

	return scanProject(row)
}

func (r *Repository) Delete(
	ctx context.Context,
	organizationID string,
	projectID string,
	identityID string,
) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM projects p
		WHERE p.id = $1
			AND p.organization_id = $2
			AND EXISTS (
				SELECT 1
				FROM org_members om
				WHERE om.organization_id = p.organization_id
					AND om.identity_id = $3
					AND om.status = 'active'
			)
	`,
		projectID,
		organizationID,
		identityID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func scanProject(row pgx.Row) (*Project, error) {
	var project Project
	var redirectURLs []byte

	if err := row.Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.AuthOptions,
		&redirectURLs,
		&project.AllowedOrigins,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if len(redirectURLs) > 0 {
		if err := json.Unmarshal(redirectURLs, &project.RedirectURLs); err != nil {
			return nil, err
		}
	}

	return &project, nil
}
