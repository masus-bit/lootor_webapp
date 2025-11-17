package models

import (
	"gorm.io/gorm"
	"time"
)

type Photos struct {
	CreatedAt string         `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id            uint64          `gorm:"type:uint;primaryKey;autoIncrement" json:"id"`
	Likes         int64           `json:"likes"`
	CollectionId  string          `json:"collectionId"`
	Path          string          `json:"path"`
	Author        SubUsers        `json:"author"`
	Collection    CollectionShort `json:"collection"`
	Tags          []ShortTags     `json:"tags"`
	CommentsCount int64           `json:"commentsCount"`
}

type PhotosDataResponse struct {
	Data  []Photos `json:"data"`
	Total int64    `json:"total"`
}

type PhotoDataResponse struct {
	Data Photos `json:"data"`
}

type PhotoCreateRequest struct {
	CollectionId string   `json:"collectionId"`
	Paths        []string `json:"paths"`
	Author       string   `json:"author"`
}

type PhotoCreate struct {
	CollectionId string   `json:"collectionId"`
	Paths        []string `json:"paths"`
}

type PhotoUpdateRequest struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type FindPhotosByCollectionsIdsRequest struct {
	CollectionIds []string `json:"collectionIds"`
}
