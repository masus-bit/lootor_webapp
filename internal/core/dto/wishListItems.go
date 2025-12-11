package dto

import (
	"github.com/google/uuid"
)

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
