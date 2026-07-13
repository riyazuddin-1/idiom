package identities

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

func (r *Repository) GetIdentityByEmail(ctx context.Context, email, projectID string) (*Identity, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, project_id, email_verified
		FROM identities
		WHERE email = $1 AND project_id = $2
	`, email, projectID)

	var identity Identity
	var passwordHash string
	if err := row.Scan(&identity.ID, &identity.Email, &passwordHash, &identity.ProjectID, &identity.EmailVerified); err != nil {
		return nil, err
	}

	identity.PasswordHash = &passwordHash
	return &identity, nil
}

func (r *Repository) VerifyIdentityEmail(ctx context.Context, email, projectID string) error {
	result, err := r.db.Exec(ctx, `
		UPDATE identities
		SET email_verified = true,
			updated_at = NOW()
		WHERE email = $1 AND project_id = $2
	`, email, projectID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("identity not found")
	}
	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, email, projectID, passwordHash string) error {
	result, err := r.db.Exec(ctx, `
		UPDATE identities
		SET password_hash = $3,
			updated_at = NOW()
		WHERE email = $1 AND project_id = $2
	`, email, projectID, passwordHash)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("identity not found")
	}
	return nil
}

func (r *Repository) CreateIdentity(ctx context.Context, identity *Identity) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO identities (
			id,
			project_id,
			email,
			email_verified,
			first_name,
			last_name,
			display_name,
			avatar_url,
			password_hash,
			provider,
			provider_user_id,
			status,
			created_at,
			updated_at,
			last_login_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15
		)
	`,
		identity.ID,
		identity.ProjectID,
		identity.Email,
		identity.EmailVerified,
		identity.FirstName,
		identity.LastName,
		identity.DisplayName,
		identity.AvatarURL,
		*identity.PasswordHash,
		identity.Provider,
		identity.ProviderUserID,
		identity.Status,
		identity.CreatedAt,
		identity.UpdatedAt,
		identity.LastLoginAt,
	)
	return err
}
