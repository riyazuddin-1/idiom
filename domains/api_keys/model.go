package apikeys

import "time"

type APIKey struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	ProjectID string     `json:"project_id"`
	KeyPrefix string     `json:"key_prefix"`
	KeyHash   string     `json:"-"`
	Scopes    []string   `json:"scopes"`
	IsActive  bool       `json:"is_active"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
