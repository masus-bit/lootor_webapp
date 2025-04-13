package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type EntitiesCollectionItemCollectionItem struct {
	EntitiesId        uuid.UUID `gorm:"primaryKey;column:entitiesId"`
	CollectionItemsId uuid.UUID `gorm:"primaryKey;column:collectionItemsId"`
}

type CollectionItems struct {
	gorm.Model
	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string
	Description     string
	Images          pq.StringArray `gorm:"type:text[]"`
	PurchaseDate    string
	PurchasePrice   float64
	Sealed          bool
	Edition         string
	CopyNumber      pq.StringArray `gorm:"type:text[]"`
	Deleted         bool
	ShippingCost    float64
	Transliteration string
	Collections     []Collections `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;"`
	Owner           Users         `gorm:"constraint:OnDelete:CASCADE;"`
	Platform        Platforms     `gorm:"constraint:OnDelete:CASCADE;"`
	Entities        []Entities    `gorm:"many2many:entities_collection_item_collection_items;constraint:OnDelete:CASCADE;"`
}

func (EntitiesCollectionItemCollectionItem) TableName() string {
	return "entities_collection_item_collection_items"
}

type CollectionItemsResponse struct {
	Id            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Images        []string   `json:"images"`
	PurchaseDate  string     `json:"purchaseDate"`
	PurchasePrice float64    `json:"purchasePrice"`
	Sealed        bool       `json:"sealed"`
	Edition       string     `json:"edition"`
	CopyNumber    []string   `json:"copyNumber"`
	ShippingCost  float64    `json:"shippingCost"`
	Entities      []Entities `json:"entities"`
	Platform      Platforms  `json:"platform"`
	Collection    uuid.UUID  `json:"collection"`
}

type CollectionItemsRequestCreate struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	PurchaseDate  string   `json:"purchaseDate"`
	PurchasePrice float64  `json:"purchasePrice"`
	Sealed        bool     `json:"sealed"`
	Edition       string   `json:"edition"`
	CopyNumber    []string `json:"copyNumber"`
	ShippingCost  float64  `json:"shippingCost"`
	Entities      []string `json:"entities"`
	Platform      string   `json:"platform"`
	CollectionId  string   `json:"collectionId"`
}

type CollectionItemsDataResponse struct {
	Data CollectionItemsResponse `json:"data"`
}
