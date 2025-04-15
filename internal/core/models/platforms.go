package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Platforms struct {
	gorm.Model
	Id   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string    `json:"name"`
}

type PlatformsDataResponse struct {
	Data []Platforms `json:"data"`
}
