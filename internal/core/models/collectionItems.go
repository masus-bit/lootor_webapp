package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type CollectionItems struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name            string
	Description     string
	Images          pq.StringArray `gorm:"type:text[]"`
	PurchaseDate    string
	PurchasePrice   float64
	Sealed          bool
	Edition         string
	CopyNumber      pq.Int64Array `gorm:"type:integer[]"`
	Rating          float64
	Likes           pq.StringArray `gorm:"type:text[]" json:"likes"`
	Deleted         bool
	ShippingCost    float64
	Transliteration string
	Collections     []Collections `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;"`
	Owner           Users         `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`
	Platform        *Platforms    `gorm:"foreignKey:PlatformID;references:ID;constraint:OnDelete:SET NULL;"`
	ItemType        *ItemTypes    `gorm:"foreignKey:ItemTypeID;references:ID;constraint:OnDelete:SET NULL;"`
	UserLogin       string        `gorm:"type:varchar(255);index"`
	PlatformID      *uuid.UUID    `gorm:"type:uuid;index"`
	ItemTypeID      *uuid.UUID    `gorm:"type:uuid;index"`
	ReportsCount    int64         `json:"-"`
	CommentsCount   int64         `json:"commentsCount"`
}
