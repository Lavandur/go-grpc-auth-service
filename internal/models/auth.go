package models

import "time"

type AuthResponse struct {
	PublicToken       string
	PublicTokenExpiry time.Time
}
