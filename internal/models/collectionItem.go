package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type EntityModelCollectionItemCollectionItem struct {
	EntityModelId    uuid.UUID `gorm:"primaryKey;column:entityModelId"`
	CollectionItemId uuid.UUID `gorm:"primaryKey;column:collectionItemId"`
}

type CollectionItem struct {
	gorm.Model
	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string
	Description     string
	Images          pq.StringArray `gorm:"type:text[]"`
	PurchaseDate    string
	PurchasePrice   int
	Sealed          bool
	Edition         string
	CopyNumber      pq.StringArray `gorm:"type:text[]"`
	Deleted         bool
	ShippingCost    int
	Transliteration string
	Collections     []Collection  `gorm:"many2many:collection_collection_items;constraint:OnDelete:CASCADE;"`
	Owner           User          `gorm:"constraint:OnDelete:CASCADE;"`
	Platform        Platforms     `gorm:"constraint:OnDelete:CASCADE;"`
	Entities        []EntityModel `gorm:"many2many:entity_model_collection_item_collection_item;constraint:OnDelete:CASCADE;"`
}

func (CollectionItem) TableName() string {
	return "collection_item"
}

func (EntityModelCollectionItemCollectionItem) TableName() string {
	return "entity_model_collection_item_collection_item"
}
