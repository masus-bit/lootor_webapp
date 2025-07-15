package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Tags struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string        `json:"name"`
	Collections []Collections `gorm:"many2many:collection_tags_collection;constraint:OnDelete:CASCADE;" json:"collections"`
}

type TagsDataResponse struct {
	Data []Tags `json:"data"`
}
