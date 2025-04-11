package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Events struct {
	gorm.Model
	Id                   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Action               string
	Date                 string
	EventTarget          string
	TargetName           string
	User                 Users           `gorm:"constraint:OnDelete:CASCADE;"`
	TargetUser           Users           `gorm:"constraint:OnDelete:CASCADE;"`
	TargetCollection     Collections     `gorm:"constraint:OnDelete:CASCADE;"`
	TargetCollectionItem CollectionItems `gorm:"constraint:OnDelete:CASCADE;"`
}
