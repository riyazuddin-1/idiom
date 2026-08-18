package projects

import "time"

type RedirectURL struct {
	Device string `json:"device"`
	URL    string `json:"url"`
}

type Project struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"org_id"`
	Name           string        `json:"name"`
	Slug           string        `json:"slug"`
	Description    string        `json:"description"`
	Status         string        `json:"status"`
	AuthOptions    []string      `json:"auth_options"`
	RedirectURLs   []RedirectURL `json:"redirect_urls"`
	AllowedOrigins []string      `json:"allowed_origins"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
