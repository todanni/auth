package models

import "time"

type RefreshToken struct {
	Value     string    `json:"value"`
	UserID    int       `json:"userID"`
	Revoked   bool      `json:"revoked"`
	ExpiresAt time.Time `json:"expiresAt"`
}
