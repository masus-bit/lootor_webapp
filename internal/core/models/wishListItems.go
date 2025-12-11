package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type WishListItems struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserLogin string    `gorm:"type:varchar(255);index"`
	User      Users     `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`

	CollectionItemID *uuid.UUID       `gorm:"type:uuid;index"`
	CollectionItem   *CollectionItems `gorm:"foreignKey:CollectionItemID;references:ID;constraint:OnDelete:SET NULL;"`

	ItemName string
	Images   pq.StringArray `gorm:"type:text[]"`

	PurchaseLinks pq.StringArray `gorm:"type:text[]"`
	Priority      int64
	Notes         string
	ReportsCount  int64 `json:"-"`
}
