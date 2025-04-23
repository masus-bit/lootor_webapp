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

	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserLogin string    `gorm:"type:varchar(255);index"`
	User      Users     `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`

	CollectionItemID *uuid.UUID       `gorm:"type:uuid;index"`
	CollectionItem   *CollectionItems `gorm:"foreignKey:CollectionItemID;references:Id;constraint:OnDelete:SET NULL;"`

	PurchaseLinks pq.StringArray `gorm:"type:text[]"`
	Priority      int64
	Notes         string
}

type WishListDataResponse struct {
	Data []WishListItemResponse `json:"data"`
}

type WishListCreateRequest struct {
	CollectionItemId string   `json:"collectionItemId"`
	PurchaseLinks    []string `json:"purchaseLinks"`
	Notes            string   `json:"notes"`
	Priority         int64    `json:"priority"`
}

type WishListItemResponse struct {
	Id             uuid.UUID               `json:"id"`
	User           UserResponse            `json:"user"`
	CollectionItem CollectionItemsResponse `json:"collectionItem"`
	PurchaseLinks  []string                `json:"purchaseLinks"`
	Priority       int64                   `json:"priority"`
	Notes          string                  `json:"notes"`
}
