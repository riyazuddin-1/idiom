package organizations

import (
	"context"
	"errors"
	"strings"

	"idiom-api-services/packages/crypto"
)

type CreateInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ListByIdentityID(ctx context.Context, repo *Repository, identityID string) ([]Organization, error) {
	return repo.ListByIdentityID(ctx, identityID)
}

func Create(
	ctx context.Context,
	repo *Repository,
	identityID string,
	input CreateInput,
) (*Organization, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("organization name is required")
	}

	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	if slug == "" {
		slug = buildSlug(name)
	}

	organization := &Organization{
		ID:     crypto.GenerateID("org_"),
		Name:   name,
		Slug:   slug,
		Status: "active",
	}

	if err := repo.Create(ctx, organization, identityID); err != nil {
		return nil, err
	}

	return organization, nil
}

func Update(
	ctx context.Context,
	repo *Repository,
	identityID string,
	organizationID string,
	input UpdateInput,
) (*Organization, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(strings.ToLower(input.Slug))

	if name == "" {
		return nil, errors.New("organization name is required")
	}

	if slug == "" {
		slug = buildSlug(name)
	}

	organization, err := repo.Update(
		ctx,
		organizationID,
		identityID,
		name,
		slug,
	)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func UpdateStatus(
	ctx context.Context,
	repo *Repository,
	identityID string,
	organizationID string,
	status string,
) (*Organization, error) {
	status = strings.TrimSpace(strings.ToLower(status))

	if status != "active" && status != "inactive" {
		return nil, errors.New("invalid organization status")
	}

	organization, err := repo.UpdateStatus(
		ctx,
		organizationID,
		identityID,
		status,
	)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func Delete(
	ctx context.Context,
	repo *Repository,
	identityID string,
	organizationID string,
) error {
	return repo.Delete(ctx, organizationID, identityID)
}

func buildSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))

	var builder strings.Builder
	lastWasDash := false

	for _, char := range slug {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
			lastWasDash = false
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastWasDash = false
		case !lastWasDash:
			builder.WriteRune('-')
			lastWasDash = true
		}
	}

	return strings.Trim(builder.String(), "-")
}
