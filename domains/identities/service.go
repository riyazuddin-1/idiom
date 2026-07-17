package identities

import (
	"context"
	"idiom-api-services/packages/crypto"
	"strings"
)

type RegisterInput struct {
	ProjectID   string
	Email       string
	Password    string
	FirstName   string
	LastName    string
	DisplayName string
	AvatarURL   string
}

func Login(ctx context.Context, repo *Repository, projectID, email, password string) (*Identity, bool, error) {
	identity, err := repo.GetIdentityByEmail(ctx, projectID, email)
	if err != nil {
		return nil, false, err
	}

	if identity.PasswordHash == nil {
		return nil, false, nil
	}

	if !identity.EmailVerified {
		return nil, false, nil
	}

	if !crypto.CheckPasswordHash(password, *identity.PasswordHash) {
		return nil, false, nil
	}

	return identity, true, nil
}

func VerifyEmail(ctx context.Context, repo *Repository, projectID, email string) error {
	return repo.VerifyIdentityEmail(ctx, projectID, strings.ToLower(email))
}

func UpdatePassword(ctx context.Context, repo *Repository, projectID, email, password string) error {
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	return repo.UpdatePassword(ctx, projectID, strings.ToLower(email), hashedPassword)
}

func Register(ctx context.Context, repo *Repository, input RegisterInput) (*Identity, error) {
	hashedPassword, err := crypto.HashPassword(input.Password)

	if err != nil {
		return nil, err
	}

	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		displayName = buildDisplayName(input.FirstName, input.LastName, input.Email)
	}

	identity := &Identity{
		ID:            crypto.GenerateID("uid_"),
		ProjectID:     input.ProjectID,
		Email:         strings.ToLower(strings.TrimSpace(input.Email)),
		EmailVerified: true,
		FirstName:     strings.TrimSpace(input.FirstName),
		LastName:      strings.TrimSpace(input.LastName),
		DisplayName:   displayName,
		AvatarURL:     strings.TrimSpace(input.AvatarURL),
		PasswordHash:  &hashedPassword,
	}

	err = repo.CreateIdentity(ctx, identity)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func buildDisplayName(firstName, lastName, email string) string {
	name := strings.TrimSpace(strings.Join([]string{
		strings.TrimSpace(firstName),
		strings.TrimSpace(lastName),
	}, " "))
	if name != "" {
		return name
	}

	if at := strings.Index(email, "@"); at > 0 {
		return email[:at]
	}

	return email
}
