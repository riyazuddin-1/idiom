package apikeys

import (
	"context"
	"errors"

	"idiom-api-services/packages/database/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByProjectID(
	ctx context.Context,
	projectID string,
	identityID string,
) ([]APIKey, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			ak.id,
			ak.name,
			ak.project_id,
			ak.key_prefix,
			ak.key_hash,
			COALESCE(ak.scopes, ARRAY[]::text[]),
			ak.is_active,
			ak.created_by,
			ak.created_at,
			ak.revoked_at
		FROM api_keys ak
		INNER JOIN projects p
			ON p.id = ak.project_id
		INNER JOIN org_members om
			ON om.organization_id = p.organization_id
		WHERE ak.project_id = $1
			AND om.identity_id = $2
			AND om.status = 'active'
		ORDER BY ak.created_at DESC
	`, projectID, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apiKeys := make([]APIKey, 0)

	for rows.Next() {
		var apiKey APIKey

		if err := rows.Scan(
			&apiKey.ID,
			&apiKey.Name,
			&apiKey.ProjectID,
			&apiKey.KeyPrefix,
			&apiKey.KeyHash,
			&apiKey.Scopes,
			&apiKey.IsActive,
			&apiKey.CreatedBy,
			&apiKey.CreatedAt,
			&apiKey.RevokedAt,
		); err != nil {
			return nil, err
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (r *Repository) Create(
	ctx context.Context,
	apiKey *APIKey,
	identityID string,
) error {
	result, err := r.db.Exec(ctx, `
		INSERT INTO api_keys (
			id,
			name,
			project_id,
			key_prefix,
			key_hash,
			scopes,
			is_active,
			created_by
		)
		SELECT
			$1,
			$2,
			p.id,
			$3,
			$4,
			$5,
			true,
			$6
		FROM projects p
		INNER JOIN org_members om
			ON om.organization_id = p.organization_id
		WHERE p.id = $7
			AND om.identity_id = $6
			AND om.status = 'active'
	`,
		apiKey.ID,
		apiKey.Name,
		apiKey.KeyPrefix,
		apiKey.KeyHash,
		apiKey.Scopes,
		identityID,
		apiKey.ProjectID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Rotate(
	ctx context.Context,
	apiKeyID string,
	projectID string,
	identityID string,
	keyPrefix string,
	keyHash string,
) (*APIKey, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE api_keys ak
		SET
			key_prefix = $4,
			key_hash = $5
		WHERE ak.id = $1
			AND ak.project_id = $2
			AND ak.is_active = true
			AND EXISTS (
				SELECT 1
				FROM projects p
				INNER JOIN org_members om
					ON om.organization_id = p.organization_id
				WHERE p.id = ak.project_id
					AND om.identity_id = $3
					AND om.status = 'active'
			)
		RETURNING
			ak.id,
			ak.name,
			ak.project_id,
			ak.key_prefix,
			ak.key_hash,
			COALESCE(ak.scopes, ARRAY[]::text[]),
			ak.is_active,
			ak.created_by,
			ak.created_at,
			ak.revoked_at
	`,
		apiKeyID,
		projectID,
		identityID,
		keyPrefix,
		keyHash,
	)

	var apiKey APIKey

	if err := row.Scan(
		&apiKey.ID,
		&apiKey.Name,
		&apiKey.ProjectID,
		&apiKey.KeyPrefix,
		&apiKey.KeyHash,
		&apiKey.Scopes,
		&apiKey.IsActive,
		&apiKey.CreatedBy,
		&apiKey.CreatedAt,
		&apiKey.RevokedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &apiKey, nil
}

func (r *Repository) Revoke(
	ctx context.Context,
	apiKeyID string,
	projectID string,
	identityID string,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE api_keys ak
		SET
			is_active = false,
			revoked_at = NOW()
		WHERE ak.id = $1
			AND ak.project_id = $2
			AND ak.is_active = true
			AND EXISTS (
				SELECT 1
				FROM projects p
				INNER JOIN org_members om
					ON om.organization_id = p.organization_id
				WHERE p.id = ak.project_id
					AND om.identity_id = $3
					AND om.status = 'active'
			)
	`,
		apiKeyID,
		projectID,
		identityID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) GetActiveByPrefix(
	ctx context.Context,
	keyPrefix string,
) (*APIKey, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id,
			name,
			project_id,
			key_prefix,
			key_hash,
			COALESCE(scopes, ARRAY[]::text[]),
			is_active,
			created_by,
			created_at,
			revoked_at
		FROM api_keys
		WHERE key_prefix = $1
			AND is_active = true
		LIMIT 1
	`, keyPrefix)

	var apiKey APIKey

	if err := row.Scan(
		&apiKey.ID,
		&apiKey.Name,
		&apiKey.ProjectID,
		&apiKey.KeyPrefix,
		&apiKey.KeyHash,
		&apiKey.Scopes,
		&apiKey.IsActive,
		&apiKey.CreatedBy,
		&apiKey.CreatedAt,
		&apiKey.RevokedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &apiKey, nil
}
