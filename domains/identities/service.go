package identities

import (
	"context"
	"idiom-api-services/packages/crypto"
)

func Login(ctx context.Context, email, password, projectID string) (bool, error) {
	repo := &Repository{}

	identity, err := repo.GetIdentityByEmail(ctx, email, projectID)
	if err != nil {
		return false, err
	}

	return crypto.CheckPasswordHash(password, identity.Password), nil
}

func Register(ctx context.Context, email, password, projectID string) (bool, error) {
	repo := &Repository{}

	hashedPassword, err := crypto.HashPassword(password)

	if err != nil {
		return false, err
	}

	identity := &Identity{
		Email:     email,
		Password:  hashedPassword,
		ProjectID: projectID,
	}

	err = repo.CreateIdentity(ctx, identity)
	if err != nil {
		return false, err
	}

	return true, nil
}
