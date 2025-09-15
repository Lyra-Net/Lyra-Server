package models

import "time"

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    uint
	Token     string
	DeviceID  string
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
}
