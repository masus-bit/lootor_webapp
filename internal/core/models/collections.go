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
	UserLogin        string            `gorm:"type:varchar(255);index"`
	ReportsCount     int64             `json:"-"`
	CommentsCount    int64             `json:"commentsCount"`
}

func (CollectionsCollectionItemsCollectionItems) TableName() string {
	return "collections_collection_items_collection_items"
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
	Tags                 []ShortTags               `json:"tags"`
	CollectionItemsCount int64                     `json:"collectionItemsCount"`
	LikesCount           int64                     `json:"likesCount"`
	CreatedAt            time.Time                 `json:"createdAt"`
	CanLike              bool                      `json:"canLike"`
	IsOwner              bool                      `json:"isOwner"`
	CommentsCount        int64                     `json:"commentsCount"`
}

type CollectionDataResponse struct {
	Data CollectionsResponse `json:"data"`
}

type AllCollectionsDataResponse struct {
	Data        []CollectionsResponse `json:"data"`
	Total       int64                 `json:"total"`
	ProfileName string                `json:"profileName"`
}

type AllCollectionsDataByTag struct {
	Data  []CollectionsResponse `json:"data"`
	Total int64                 `json:"total"`
	Tag   Tags                  `json:"tag"`
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

type CollectionUpdateRequest struct {
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	IsPrivate       *bool    `json:"isPrivate"`
	BannerUrl       *string  `json:"bannerUrl"`
	Transliteration *string  `json:"transliteration"`
	UserLogin       *string  `json:"userLogin"`
	Tags            []string `json:"tags"`
}
