package apikeys

import "time"

type APIKeys struct {
	Name      string    `json:"name"`
	ProjectID string    `json:"project_id"`
	KeyPrefix string    `json:"key_prefix"`
	KeyHash   string    `json:"key_hash"`
	Scopes    []string  `json:"scopes"`
	IsActive  bool      `json:"is_active"`
	CreatedBy string    `json:"created_by"`
	CreatedAt string    `json:"created_at"`
	RevokedAt time.Time `json:"revoked_at"`
}
