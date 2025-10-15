package models

import (
	"gorm.io/gorm"
	"time"
)

type Tags struct {
	CreatedAt string         `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	PrimaryID   string `gorm:"type:uuid" json:"primaryId"`
	Primary     *Tags  `gorm:"foreignKey:PrimaryID" json:"primary"`
	IsPrimary   bool   `gorm:"default:true" json:"isPrimary"`
	Description string `gorm:"type:text" json:"description"`

	Slug   string `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Author string `json:"author"`

	TotalPosts           int64 `json:"totalPosts"`
	TotalCollections     int64 `json:"totalCollections"`
	TotalCollectionItems int64 `json:"totalCollectionItems"`

	TagLinks []TagLinks `gorm:"foreignKey:TagID" json:"tagLinks"`
	Synonyms []Tags     `gorm:"foreignKey:PrimaryID" json:"synonyms"`

	Entities Entities `json:"entities"`

	CollectionItemsProps *CollectionItemsProps `json:"collectionItemsProps"`
}

type CollectionItemsProps struct {
	Books              int64 `json:"books"`
	VideoGames         int64 `json:"videoGames"`
	BoardGames         int64 `json:"boardGames"`
	Comics             int64 `json:"comics"`
	GamingHardware     int64 `json:"gamingHardware"`
	Vinyl              int64 `json:"vinyl"`
	Steelbooks         int64 `json:"steelbooks"`
	CollectibleCards   int64 `json:"collectibleCards"`
	CollectibleFigures int64 `json:"collectibleFigures"`
}

type Entities struct {
	Posts           []Posts                   `json:"posts"`
	CollectionItems []CollectionItemsResponse `json:"collectionItems"`
	Collections     []CollectionsResponse     `json:"collections"`
}

type TagLinks struct {
	TagID      string    `gorm:"type:uuid;primaryKey" json:"tagId"`
	EntityType string    `gorm:"type:varchar(50);primaryKey" json:"entityType"`
	EntityID   string    `gorm:"type:uuid;primaryKey" json:"entityId"`
	CreatedAt  time.Time `gorm:"default:now()" json:"createdAt"`

	Tag Tags `gorm:"foreignKey:TagID" json:"tag"`
}

type ShortTags struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TagsDataResponse struct {
	Data []Tags `json:"data"`
}

type TagDataResponse struct {
	Data  Tags  `json:"data"`
	Total int64 `json:"total"`
}

type TagCreateRequest struct {
	Names      []string `json:"names"`
	Author     string   `json:"author"`
	EntityID   string   `json:"entityId"`
	EntityType string   `json:"entityType"`
}

type TagsSearchRequest struct {
	Name string `json:"name"`
}

type MergeTagsRequest struct {
	FromTagIDs []string `json:"fromTagIds"`
	ToTagID    string   `json:"toTagId"`
}

type AddTagToEntityRequest struct {
	TagID      string `json:"tagId"`
	EntityType string `json:"entityType"`
	EntityID   string `json:"entityId"`
}

type RemoveTagsRequest struct {
	EntityID   string   `json:"entityId"`
	TagIDs     []string `json:"tagIds"`
	EntityType string   `json:"entityType"`
}

type TagUpdateRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type TagIDsResponse struct {
	Data []string `json:"data"`
}
