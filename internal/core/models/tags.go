package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tags struct {
	gorm.Model
	Id          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string        `json:"name"`
	Collections []Collections `gorm:"many2many:collection_tags_collection;constraint:OnDelete:CASCADE;" json:"collections"`
}
