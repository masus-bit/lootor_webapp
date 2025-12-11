package dto

import (
	"github.com/google/uuid"
	"lootor/internal/core/models"
)

type CollectionItemsResponse struct {
	ID                        uuid.UUID         `json:"id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	Images                    []string          `json:"images"`
	PurchaseDate              string            `json:"purchaseDate"`
	PurchasePrice             float64           `json:"purchasePrice"`
	Sealed                    bool              `json:"sealed"`
	Edition                   string            `json:"edition"`
	CopyNumber                []int64           `json:"copyNumber"`
	Rating                    float64           `json:"rating"`
	ShippingCost              float64           `json:"shippingCost"`
	Entities                  []Entities        `json:"entities"`
	Platform                  *models.Platforms `json:"platform"`
	Collection                uuid.UUID         `json:"collection"`
	Owner                     models.Users      `json:"owner"`
	ItemType                  *models.ItemTypes `json:"itemType"`
	CanLike                   bool              `json:"canLike"`
	LikesCount                int64             `json:"likesCount"`
	IsOwner                   bool              `json:"isOwner"`
	CollectionTransliteration string            `json:"collectionTransliteration"`
	CommentsCount             int64             `json:"commentsCount"`
	CollectionName            string            `json:"collectionName"`
	Tags                      []ShortTags       `json:"tags"`
}

type CollectionItemsRequestCreate struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	PurchaseDate  string   `json:"purchaseDate"`
	PurchasePrice float64  `json:"purchasePrice"`
	Sealed        bool     `json:"sealed"`
	Edition       string   `json:"edition"`
	Rating        float64  `json:"rating"`
	CopyNumber    []int64  `json:"copyNumber"`
	ShippingCost  float64  `json:"shippingCost"`
	Platform      string   `json:"platform"`
	Collection    string   `json:"collection"`
	ItemType      string   `json:"itemType"`
	Tags          []string `json:"tags"`
}

type CollectionItemsRequestUpdate struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Images        []string `json:"images"`
	PurchaseDate  *string  `json:"purchaseDate"`
	PurchasePrice *float64 `json:"purchasePrice"`
	Sealed        *bool    `json:"sealed"`
	Edition       *string  `json:"edition"`
	Rating        *float64 `json:"rating"`
	CopyNumber    []int64  `json:"copyNumber"`
	ShippingCost  *float64 `json:"shippingCost"`
	Platform      *string  `json:"platform"`
	Collection    *string  `json:"collection"`
	ItemType      *string  `json:"itemType"`
	Tags          []string `json:"tags"`
}

type CollectionItemsDataResponse struct {
	Data CollectionItemsResponse `json:"data"`
}

type CollectionItemsDataResponseWithCount struct {
	Data   []CollectionItemsResponse `json:"data"`
	Total  int64                     `json:"total"`
	Entity Entities                  `json:"entity"`
}

type CollectionItemsCopyOrMoveRequest struct {
	ID                  string   `json:"id"`
	TargetCollectionIDs []string `json:"targetCollectionIds"`
	SourceCollectionID  string   `json:"sourceCollectionId"`
}

type CollectionItemsSortedResponse struct {
	VideoGames         []CollectionItemsResponse `json:"videoGames"`
	BoardGames         []CollectionItemsResponse `json:"boardGames"`
	Comics             []CollectionItemsResponse `json:"comics"`
	GamingHardware     []CollectionItemsResponse `json:"gamingHardware"`
	CollectibleFigures []CollectionItemsResponse `json:"collectibleFigures"`
	Books              []CollectionItemsResponse `json:"books"`
	Vinyl              []CollectionItemsResponse `json:"vinyl"`
	Steelbooks         []CollectionItemsResponse `json:"steelbooks"`
	CollectibleCards   []CollectionItemsResponse `json:"collectibleCards"`
}

type CollectionItemsDataSortedResponse struct {
	Data   CollectionItemsSortedResponse `json:"data"`
	Total  int64                         `json:"total"`
	Entity Entities                      `json:"entity"`
}

type CollectionItemsDataPoor struct {
	Data []CollectionItemsResponse `json:"data"`
}
