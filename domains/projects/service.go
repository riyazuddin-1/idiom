package projects

import "context"

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
