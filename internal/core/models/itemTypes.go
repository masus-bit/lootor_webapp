package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ItemTypes struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name   string    `json:"name"`
	RuName string    `json:"ruName"`
	Slug   string    `json:"slug"`
}

type ItemTypesDataResponse struct {
	Data []ItemTypes `json:"data"`
}
