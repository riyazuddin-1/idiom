package org_members

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidMemberStatus = errors.New("invalid member status")
	ErrMemberNotFound      = errors.New("member not found")
)

func ListByOrganizationID(
	ctx context.Context,
	repo *Repository,
	organizationID string,
) ([]OrgMembers, error) {
	return repo.ListByOrganizationID(ctx, organizationID)
}

func UpdateStatus(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	identityID string,
	status string,
) (*OrgMembers, error) {
	status = strings.TrimSpace(strings.ToLower(status))

	if status != "active" && status != "inactive" {
		return nil, ErrInvalidMemberStatus
	}

	member, err := repo.UpdateStatus(
		ctx,
		organizationID,
		identityID,
		status,
	)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func Delete(
	ctx context.Context,
	repo *Repository,
	organizationID string,
	identityID string,
) error {
	return repo.Delete(ctx, organizationID, identityID)
}
