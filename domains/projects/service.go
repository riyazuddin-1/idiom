package projects

import "context"

func IsActive(ctx context.Context, repo *Repository, projectID string) (bool, error) {
	if projectID == "" {
		return false, nil
	}

	return repo.ActiveExists(ctx, projectID)
}
