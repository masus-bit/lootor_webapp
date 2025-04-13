package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entities struct {
	gorm.Model
	Id             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name           string
	CollectionItem []CollectionItems `gorm:"many2many:entities_collection_item_collection_items;constraint:OnDelete:CASCADE;"`
}
