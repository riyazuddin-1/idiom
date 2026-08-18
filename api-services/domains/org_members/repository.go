package org_members

import (
	"context"
	"errors"

	"idiom-api-services/packages/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByOrganizationID(
	ctx context.Context,
	organizationID string,
) ([]OrgMembers, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			organization_id,
			identity_id,
			role,
			status,
			invited_by,
			joined_at
		FROM org_members
		WHERE organization_id = $1
		ORDER BY joined_at DESC
	`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]OrgMembers, 0)

	for rows.Next() {
		var member OrgMembers

		if err := rows.Scan(
			&member.OrganizationID,
			&member.IdentityID,
			&member.Role,
			&member.Status,
			&member.InvitedBy,
			&member.JoinedAt,
		); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *Repository) UpdateStatus(
	ctx context.Context,
	organizationID string,
	identityID string,
	status string,
) (*OrgMembers, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE org_members
		SET status = $3
		WHERE organization_id = $1
			AND identity_id = $2
		RETURNING
			organization_id,
			identity_id,
			role,
			status,
			invited_by,
			joined_at
	`,
		organizationID,
		identityID,
		status,
	)

	var member OrgMembers

	if err := row.Scan(
		&member.OrganizationID,
		&member.IdentityID,
		&member.Role,
		&member.Status,
		&member.InvitedBy,
		&member.JoinedAt,
	); err != nil {
		return nil, errors.New("member not found")
	}

	return &member, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	organizationID string,
	identityID string,
) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM org_members
		WHERE organization_id = $1
			AND identity_id = $2
	`,
		organizationID,
		identityID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("member not found")
	}

	return nil
}
