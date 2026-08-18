package sessions

import (
	"time"
)

type Session struct {
	ID               string     `json:"id"`
	IdentityID       string     `json:"identity_id"`
	RefreshTokenHash string     `json:"refresh_token_hash"`
	IP               string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
	ExpiresAt        time.Time  `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
}
