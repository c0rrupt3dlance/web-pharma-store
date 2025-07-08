package models

import "time"

type RefreshToken struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	IssuedAt  time.Time `json:"Issued_at"`
}
