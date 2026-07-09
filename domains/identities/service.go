package identities

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"idiom-api-services/packages/crypto"
	"strings"
	"time"
)

const (
	ProviderEmail  = "email"
	StatusActive   = "active"
	identityIDSize = 16
)

type RegisterInput struct {
	Email       string
	Password    string
	ProjectID   string
	FirstName   string
	LastName    string
	DisplayName string
	AvatarURL   string
}

func Login(ctx context.Context, email, password, projectID string) (bool, error) {
	repo := &Repository{}

	identity, err := repo.GetIdentityByEmail(ctx, email, projectID)
	if err != nil {
		return false, err
	}

	if identity.PasswordHash == nil {
		return false, nil
	}

	if !identity.EmailVerified {
		return false, nil
	}

	return crypto.CheckPasswordHash(password, *identity.PasswordHash), nil
}

func VerifyEmail(ctx context.Context, email, projectID string) error {
	repo := &Repository{}
	return repo.VerifyIdentityEmail(ctx, strings.ToLower(strings.TrimSpace(email)), projectID)
}

func Register(ctx context.Context, input RegisterInput) (*Identity, error) {
	repo := &Repository{}

	hashedPassword, err := crypto.HashPassword(input.Password)

	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		displayName = buildDisplayName(input.FirstName, input.LastName, input.Email)
	}

	identity := &Identity{
		ID:             newIdentityID(),
		ProjectID:      input.ProjectID,
		Email:          strings.ToLower(strings.TrimSpace(input.Email)),
		EmailVerified:  false,
		FirstName:      strings.TrimSpace(input.FirstName),
		LastName:       strings.TrimSpace(input.LastName),
		DisplayName:    displayName,
		AvatarURL:      strings.TrimSpace(input.AvatarURL),
		PasswordHash:   &hashedPassword,
		Provider:       ProviderEmail,
		ProviderUserID: "",
		Status:         StatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastLoginAt:    nil,
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

func newIdentityID() string {
	bytes := make([]byte, identityIDSize)
	if _, err := rand.Read(bytes); err != nil {
		return "idn_" + hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}

	return "idn_" + hex.EncodeToString(bytes)
}
