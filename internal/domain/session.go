package domain

import "time"

type Session struct {
	ID           string      `json:"id"`
	User         User        `json:"user"`
	Fingerprint  Fingerprint `json:"fingerprint"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    time.Time   `json:"expires_at"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type Fingerprint struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}
