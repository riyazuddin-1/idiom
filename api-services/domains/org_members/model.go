package org_members

type OrgMembers struct {
	OrganizationID string `json:"org_id"`
	IdentityID     string `json:"identity_id"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	InvitedBy      string `json:"invited_by"`
	JoinedAt       string `json:"joined_at"`
}
