package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tags struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string
	Collections []Collection `gorm:"many2many:collection_tags_collection;constraint:OnDelete:CASCADE;"`
}
