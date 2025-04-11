package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionCollectionItemsCollectionItem struct {
	CollectionID     uuid.UUID `gorm:"primaryKey;column:collectionId"`
	CollectionItemID uuid.UUID `gorm:"primaryKey;column:collectionItemId"`
}

type CollectionTagsCollection struct {
	CollectionID uuid.UUID `gorm:"primaryKey;column:collectionId"`
	TagsID       uuid.UUID `gorm:"primaryKey;column:tagsId"`
}

type Collection struct {
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
	CollectionItems  []CollectionItem `gorm:"many2many:collection_collection_items_collection_item;constraint:OnDelete:CASCADE;"`
	User             User             `gorm:"constraint:OnDelete:CASCADE;"`
	Tags             []Tags           `gorm:"many2many:collection_tags_collection;constraint:OnDelete:CASCADE;"`
}

func (Collection) TableName() string {
	return "collection"
}

func (CollectionCollectionItemsCollectionItem) TableName() string {
	return "collection_collection_items_collection_item"
}

func (CollectionTagsCollection) TableName() string {
	return "collection_tags_collection"
}
