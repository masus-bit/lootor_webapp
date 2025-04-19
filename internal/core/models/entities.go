package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Entities struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id             uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string            `json:"name"`
	CollectionItem []CollectionItems `gorm:"many2many:entities_collection_item_collection_items;constraint:OnDelete:CASCADE;" json:"collectionItem"`
}

type EntitiesCreateRequest struct {
	Names []string `json:"names"`
}

type EntitiesDataResponse struct {
	Data []Entities `json:"data"`
}
