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

func (r *Repository) GetIdentityByEmail(ctx context.Context, projectID, email string) (*Identity, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, project_id, email, password_hash, email_verified
		FROM identities
		WHERE project_id = $1
			AND email = $2
	`, projectID, email)

	var identity Identity
	var passwordHash string
	if err := row.Scan(&identity.ID, &identity.ProjectID, &identity.Email, &passwordHash, &identity.EmailVerified); err != nil {
		return nil, err
	}

	identity.PasswordHash = &passwordHash
	return &identity, nil
}

func (r *Repository) VerifyIdentityEmail(ctx context.Context, projectID, email string) error {
	result, err := r.db.Exec(ctx, `
		UPDATE identities
		SET email_verified = true,
			updated_at = NOW()
		WHERE project_id = $1
			AND email = $2
	`, projectID, email)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("identity not found")
	}
	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, projectID, email, passwordHash string) error {
	result, err := r.db.Exec(ctx, `
		UPDATE identities
		SET password_hash = $3,
			updated_at = NOW()
		WHERE project_id = $1
			AND email = $2
	`, projectID, email, passwordHash)
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
			password_hash
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9
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
	)
	return err
}
