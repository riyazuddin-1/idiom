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
