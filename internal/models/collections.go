package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionsCollectionItemsCollectionItems struct {
	CollectionID     uuid.UUID `gorm:"primaryKey;column:collectionId"`
	CollectionItemID uuid.UUID `gorm:"primaryKey;column:collectionItemId"`
}

type TagsCollectionsCollections struct {
	CollectionsID uuid.UUID `gorm:"primaryKey;column:collectionsId"`
	TagsID        uuid.UUID `gorm:"primaryKey;column:tagsId"`
}

type Collections struct {
	gorm.Model
	Id               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name             string
	Description      string
	Created          string
	IsPrivate        bool
	BannerUrl        string
	Likes            int
	Deleted          bool
	Transliteration  string
	SubscribersCount int
	TotalPrice       int
	ShippingTotal    int
	ShareString      string
	CollectionItems  []CollectionItems `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;"`
	User             Users             `gorm:"constraint:OnDelete:CASCADE;"`
	Tags             []Tags            `gorm:"many2many:tags_collections_collections;constraint:OnDelete:CASCADE;"`
}

func (CollectionsCollectionItemsCollectionItems) TableName() string {
	return "collections_collection_items_collection_items"
}

func (TagsCollectionsCollections) TableName() string {
	return "tags_collections_collections"
}

type CollectionsResponse struct {
	Id                   uuid.UUID                 `json:"id"`
	Name                 string                    `json:"name"`
	Description          string                    `json:"description"`
	Created              string                    `json:"created"`
	IsPrivate            bool                      `json:"isPrivate"`
	BannerUrl            string                    `json:"bannerUrl"`
	Likes                int                       `json:"likes"`
	Deleted              bool                      `json:"deleted"`
	Transliteration      string                    `json:"transliteration"`
	SubscribersCount     int                       `json:"subscribersCount"`
	TotalPrice           int                       `json:"totalPrice"`
	ShippingTotal        int                       `json:"shippingTotal"`
	ShareString          string                    `json:"shareString"`
	CollectionItems      []CollectionItemsResponse `json:"collectionItems"`
	User                 UserResponse              `json:"user"`
	Tags                 []Tags                    `json:"tags"`
	CollectionItemsCount int                       `json:"collectionItemsCount"`
	CanSubscribe         bool                      `json:"canSubscribe"`
}

type CollectionCreateRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	IsPrivate       bool     `json:"isPrivate"`
	BannerUrl       string   `json:"bannerUrl"`
	Transliteration string   `json:"transliteration"`
	User            string   `json:"user"`
	Tags            []string `json:"tags"`
}
