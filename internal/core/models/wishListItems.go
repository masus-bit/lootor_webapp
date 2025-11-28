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

type WishListDataResponse struct {
	Data        []WishListItemResponse `json:"data"`
	Total       int64                  `json:"total"`
	ProfileName string                 `json:"profileName"`
}

type WishListSingleDataResponse struct {
	Data WishListItemResponse `json:"data"`
}

type WishListCreateRequest struct {
	CollectionItemID string   `json:"collectionItemId,omitempty"`
	PurchaseLinks    []string `json:"purchaseLinks"`
	Notes            string   `json:"notes"`
	Priority         int64    `json:"priority"`
	ItemName         string   `json:"itemName"`
	Images           []string `json:"images"`
}

type WishListUpdateRequest struct {
	CollectionItemID *string  `json:"collectionItemId"`
	PurchaseLinks    []string `json:"purchaseLinks"`
	Notes            *string  `json:"notes"`
	Priority         *int64   `json:"priority"`
	ItemName         *string  `json:"itemName"`
	Images           []string `json:"images"`
}

type WishListItemResponse struct {
	ID             uuid.UUID                `json:"id"`
	User           UserResponse             `json:"user"`
	CollectionItem *CollectionItemsResponse `json:"collectionItem"`
	PurchaseLinks  []string                 `json:"purchaseLinks"`
	Priority       int64                    `json:"priority"`
	Notes          string                   `json:"notes"`
	ItemName       string                   `json:"itemName"`
	Images         []string                 `json:"images"`
}

type WishListItemUpdatePriority struct {
	Priority int64 `json:"priority"`
}
