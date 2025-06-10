package models

import (
	"time"
)

type Subscription struct {
	ID        uint   `gorm:"primaryKey"`
	UserLogin string `gorm:"index"`
	Type      string // monthly, yearly
	StartedAt time.Time
	ExpiresAt time.Time
	Status    string // "active", "canceled", "expired"
}
