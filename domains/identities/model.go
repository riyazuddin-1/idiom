package identities

import (
	"time"
)

type Identity struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	EmailVerified  bool       `json:"email_verified"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	DisplayName    string     `json:"display_name"`
	AvatarURL      string     `json:"avatar_url"`
	PasswordHash   *string    `json:"-"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
}
