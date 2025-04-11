package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntityModel struct {
	gorm.Model
	Id             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name           string
	CollectionItem []CollectionItem `gorm:"many2many:entity_model_collection_item_collection_item;constraint:OnDelete:CASCADE;"`
}

func (EntityModel) TableName() string {
	return "entity_model"
}
