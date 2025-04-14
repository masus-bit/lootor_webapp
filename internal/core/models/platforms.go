package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Platforms struct {
	gorm.Model
	Id   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name string
}

type PlatformsDataResponse struct {
	Data []Platforms `json:"data"`
}
