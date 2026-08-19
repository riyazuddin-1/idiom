package projects

import (
	"context"
	"errors"
	"idiom-api-services/packages/crypto"
	"strings"
)

func IsActive(ctx context.Context, repo *Repository, projectID string) (bool, error) {
	if projectID == "" {
		return false, nil
	}

	return repo.ActiveExists(ctx, projectID)
}

func ResolveActive(ctx context.Context, repo *Repository, identifier string) (*Project, bool, error) {
	if identifier == "" {
		return nil, false, nil
	}

	project, err := repo.GetActiveByIdentifier(ctx, identifier)
	if err != nil {
		return nil, false, err
	}

	return project, project != nil, nil
}

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectNameRequired  = errors.New("project name is required")
	ErrInvalidProjectStatus = errors.New("invalid project status")
)

type CreateInput struct {
	Name           string        `json:"name"`
	Slug           string        `json:"slug"`
	Description    string        `json:"description"`
	AuthOptions    []string      `json:"auth_options"`
	RedirectURLs   []RedirectURL `json:"redirect_urls"`
	AllowedOrigins []string      `json:"allowed_origins"`
}

type UpdateInput struct {
	Name           string        `json:"name"`
	Slug           string        `json:"slug"`
	Description    string        `json:"description"`
	AuthOptions    []string      `json:"auth_options"`
	RedirectURLs   []RedirectURL `json:"redirect_urls"`
	AllowedOrigins []string      `json:"allowed_origins"`
}

func ListByOrganizationID(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	identityID string,
) ([]Project, error) {
	return repo.ListByOrganizationID(ctx, organizationID, identityID)
}

func Get(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	projectID string,
	identityID string,
) (*Project, error) {
	project, err := repo.GetByID(
		ctx,
		organizationID,
		projectID,
		identityID,
	)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project, nil
}

func Create(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	identityID string,
	input CreateInput,
) (*Project, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrProjectNameRequired
	}

	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	if slug == "" {
		slug = buildSlug(name)
	}

	project := &Project{
		ID:             crypto.GenerateID("prj_"),
		OrganizationID: organizationID,
		Name:           name,
		Slug:           slug,
		Description:    strings.TrimSpace(input.Description),
		Status:         "active",
		AuthOptions:    []string{"email"},
		RedirectURLs:   nonNilRedirectURLs(input.RedirectURLs),
		AllowedOrigins: nonNilStrings(input.AllowedOrigins),
	}

	if err := repo.Create(ctx, project, identityID); err != nil {
		return nil, err
	}

	return project, nil
}

func Update(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	projectID string,
	identityID string,
	input UpdateInput,
) (*Project, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrProjectNameRequired
	}

	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	if slug == "" {
		slug = buildSlug(name)
	}

	project, err := repo.Update(
		ctx,
		organizationID,
		projectID,
		identityID,
		name,
		slug,
		strings.TrimSpace(input.Description),
		[]string{"email"},
		nonNilRedirectURLs(input.RedirectURLs),
		nonNilStrings(input.AllowedOrigins),
	)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project, nil
}

func UpdateStatus(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	projectID string,
	identityID string,
	status string,
) (*Project, error) {
	status = strings.TrimSpace(strings.ToLower(status))

	if status != "active" && status != "inactive" {
		return nil, ErrInvalidProjectStatus
	}

	project, err := repo.UpdateStatus(
		ctx,
		organizationID,
		projectID,
		identityID,
		status,
	)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project, nil
}

func Delete(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	projectID string,
	identityID string,
) error {
	return repo.Delete(
		ctx,
		organizationID,
		projectID,
		identityID,
	)
}

func buildSlug(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	var builder strings.Builder
	lastWasDash := false

	for _, char := range name {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
			lastWasDash = false

		case char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastWasDash = false

		default:
			if !lastWasDash {
				builder.WriteRune('-')
				lastWasDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}

	return values
}

func nonNilRedirectURLs(values []RedirectURL) []RedirectURL {
	if values == nil {
		return []RedirectURL{}
	}

	return values
}
