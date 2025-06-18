package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type EntitiesCollectionItemCollectionItem struct {
	EntitiesId        uuid.UUID `gorm:"primaryKey;column:entitiesId"`
	CollectionItemsId uuid.UUID `gorm:"primaryKey;column:collectionItemsId"`
}

type CollectionItems struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
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
	Platform        *Platforms    `gorm:"foreignKey:PlatformID;references:Id;constraint:OnDelete:SET NULL;"`
	ItemType        *ItemTypes    `gorm:"foreignKey:ItemTypeID;references:Id;constraint:OnDelete:SET NULL;"`
	Entities        []Entities    `gorm:"many2many:entities_collection_item_collection_items;constraint:OnDelete:CASCADE;"`
	UserLogin       string        `gorm:"type:varchar(255);index"`
	PlatformID      *uuid.UUID    `gorm:"type:uuid;index"`
	ItemTypeID      *uuid.UUID    `gorm:"type:uuid;index"`
	ReportsCount    int64         `json:"-"`
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
	CopyNumber    []int64    `json:"copyNumber"`
	Rating        float64    `json:"rating"`
	ShippingCost  float64    `json:"shippingCost"`
	Entities      []Entities `json:"entities"`
	Platform      *Platforms `json:"platform"`
	Collection    uuid.UUID  `json:"collection"`
	Owner         Users      `json:"owner"`
	ItemType      *ItemTypes `json:"itemType"`
	CanLike       bool       `json:"canLike"`
	LikesCount    int64      `json:"likesCount"`
	IsOwner       bool       `json:"isOwner"`
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
	Entities      []string `json:"entities"`
	Platform      string   `json:"platform"`
	Collection    string   `json:"collection"`
	ItemType      string   `json:"itemType"`
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
	Entities      []string `json:"entities"`
	Platform      *string  `json:"platform"`
	Collection    *string  `json:"collection"`
	ItemType      *string  `json:"itemType"`
}

type CollectionItemsDataResponse struct {
	Data CollectionItemsResponse `json:"data"`
}

type CollectionItemsDataResponseWithCount struct {
	Data  []CollectionItemsResponse `json:"data"`
	Total int64                     `json:"total"`
}

type CollectionItemsCopyOrMoveRequest struct {
	Id                  string   `json:"id"`
	TargetCollectionIds []string `json:"targetCollectionIds"`
	SourceCollectionId  string   `json:"sourceCollectionId"`
}

type CollectionItemsSortedResponse struct {
	VideoGames         []CollectionItemsResponse `json:"videoGames"`
	BoardGames         []CollectionItemsResponse `json:"boardGames"`
	Comics             []CollectionItemsResponse `json:"comics"`
	GamingHardware     []CollectionItemsResponse `json:"gamingHardware"`
	CollectibleFigures []CollectionItemsResponse `json:"collectibleFigures"`
	Books              []CollectionItemsResponse `json:"books"`
	Vinyl              []CollectionItemsResponse `json:"vinyl"`
}

type CollectionItemsDataSortedResponse struct {
	Data   CollectionItemsSortedResponse `json:"data"`
	Total  int64                         `json:"total"`
	Entity Entities                      `json:"entity"`
}
