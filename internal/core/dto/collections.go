package dto

import (
	"github.com/google/uuid"
	"time"
)

type CollectionsResponse struct {
	ID                   uuid.UUID                 `json:"id"`
	Name                 string                    `json:"name"`
	Description          string                    `json:"description"`
	Created              string                    `json:"created"`
	IsPrivate            bool                      `json:"isPrivate"`
	BannerURL            string                    `json:"bannerUrl"`
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
	PhotosCount          int64                     `json:"photosCount"`
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
	BannerURL       string   `json:"bannerUrl"`
	Transliteration string   `json:"transliteration"`
	UserLogin       string   `json:"userLogin"`
	Tags            []string `json:"tags"`
}

type CollectionUpdateRequest struct {
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	IsPrivate       *bool    `json:"isPrivate"`
	BannerURL       *string  `json:"bannerUrl"`
	Transliteration *string  `json:"transliteration"`
	UserLogin       *string  `json:"userLogin"`
	Tags            []string `json:"tags"`
}

type CollectionShort struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	IsPrivate       bool      `json:"isPrivate"`
	BannerURL       string    `json:"bannerUrl"`
	Transliteration string    `json:"transliteration"`
	ShareString     string    `json:"shareString"`
	UserLogin       string    `json:"userLogin"`
}
