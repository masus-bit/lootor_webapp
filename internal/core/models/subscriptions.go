package models

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserLogin string    `gorm:"index"`
	Type      string    // monthly, yearly
	StartedAt time.Time
	ExpiresAt time.Time
	Status    string // "active", "canceled", "expired"
}
