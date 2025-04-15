package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tags struct {
	gorm.Model
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string
	Collections []Collections `gorm:"many2many:collection_tags_collection;constraint:OnDelete:CASCADE;"`
}
