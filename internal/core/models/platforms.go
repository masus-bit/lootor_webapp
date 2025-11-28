package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Platforms struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string    `json:"name"`
}

type PlatformsDataResponse struct {
	Data []Platforms `json:"data"`
}
