package model

import "time"

type RefreshToken struct {
	ExpiresAt time.Time
	UserID    string
}
