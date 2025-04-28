package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
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
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Created          string         `json:"created"`
	IsPrivate        bool           `json:"isPrivate"`
	BannerUrl        string         `json:"bannerUrl"`
	Likes            pq.StringArray `gorm:"type:text[]" json:"likes"`
	Deleted          bool
	Transliteration  string            `json:"transliteration"`
	SubscribersCount int64             `json:"subscribersCount"`
	TotalPrice       int64             `json:"totalPrice"`
	ShippingTotal    int64             `json:"shippingTotal"`
	ShareString      string            `json:"shareString"`
	CollectionItems  []CollectionItems `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;" json:"collectionItems"`
	User             *Users            `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;" json:"user"`
	Tags             []Tags            `gorm:"many2many:tags_collections_collections;constraint:OnDelete:CASCADE;" json:"tags"`
	UserLogin        string            `gorm:"type:varchar(255);index"`
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
	Deleted              bool                      `json:"deleted"`
	Transliteration      string                    `json:"transliteration"`
	SubscribersCount     int64                     `json:"subscribersCount"`
	TotalPrice           float64                   `json:"totalPrice"`
	ShippingTotal        float64                   `json:"shippingTotal"`
	ShareString          string                    `json:"shareString"`
	CollectionItems      []CollectionItemsResponse `json:"collectionItems"`
	User                 UserResponse              `json:"user"`
	Tags                 []Tags                    `json:"tags"`
	CollectionItemsCount int64                     `json:"collectionItemsCount"`
	CanSubscribe         bool                      `json:"canSubscribe"`
	LikesCount           int64                     `json:"likesCount"`
	CreatedAt            time.Time                 `json:"createdAt"`
	CanLike              bool                      `json:"canLike"`
	IsOwner              bool                      `json:"isOwner"`
}

type CollectionDataResponse struct {
	Data CollectionsResponse `json:"data"`
}

type AllCollectionsDataResponse struct {
	Data []CollectionsResponse `json:"data"`
}

type AllCollectionsDataByTag struct {
	Data  []CollectionsResponse `json:"data"`
	Total int64                 `json:"total"`
}

type CollectionCreateRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	IsPrivate       bool     `json:"isPrivate"`
	BannerUrl       string   `json:"bannerUrl"`
	Transliteration string   `json:"transliteration"`
	UserLogin       string   `json:"userLogin"`
	Tags            []string `json:"tags"`
}
