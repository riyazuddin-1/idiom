package identities

import (
	"context"
	"idiom-api-services/packages/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func (r *Repository) GetIdentityByEmail(ctx context.Context, email, projectID string) (*Identity, error) {
	row := r.db.QueryRow(ctx, "SELECT email, password, project_id from identities WHERE email = $1 AND project_id = $2", email, projectID)

	var emailResult, passwordResult, projectIDResult string
	if err := row.Scan(&emailResult, &passwordResult, &projectIDResult); err != nil {
		return nil, err
	}

	return &Identity{
		Email:     emailResult,
		Password:  passwordResult,
		ProjectID: projectIDResult,
	}, nil
}

func (r *Repository) CreateIdentity(ctx context.Context, identity *Identity) error {
	_, err := r.db.Exec(ctx, "INSERT INTO identities (email, password, project_id) VALUES ($1, $2, $3)", identity.Email, identity.Password, identity.ProjectID)
	return err
}
