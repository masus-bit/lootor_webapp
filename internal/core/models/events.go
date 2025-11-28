package models

import (
	"github.com/google/uuid"
)

type Events struct {
	ID              string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date            string `json:"date"`
	Action          string `gorm:"type:varchar(100);not null" json:"action"`
	EventTargetType string `gorm:"type:varchar(50);not null" json:"eventTargetType"`
	TargetName      string `gorm:"type:varchar(255)" json:"targetName"`

	InitiatorLogin string    `gorm:"type:varchar(255);not null" json:"initiatorLogin"`
	Initiator      *SubUsers `gorm:"foreignKey:InitiatorLogin;references:Login" json:"initiator"`

	TargetUserLogin      string `gorm:"type:varchar(255)" json:"targetUserLogin"`
	TargetCollectionID   string `gorm:"type:uuid" json:"targetCollectionId"`
	TargetItemID         string `gorm:"type:uuid" json:"targetItemId"`
	TargetWishListItemID string `gorm:"type:uuid" json:"targetWishListItemId"`
	TargetPostID         string `json:"targetPostId"`
	TargetTagID          string `json:"targetTagId"`
	TargetPhotoID        string `json:"targetPhotoId"`

	TargetUser         *SubUsers                `gorm:"-" json:"targetUser"`
	TargetCollection   *CollectionsResponse     `gorm:"-" json:"targetCollection"`
	TargetItem         *CollectionItemsResponse `gorm:"-" json:"targetItem"`
	TargetWishListItem *WishListItemResponse    `gorm:"-" json:"targetWishListItem"`
	TargetPost         *Posts                   `gorm:"-" json:"targetPost"`
	TargetTag          *ShortTags               `gorm:"-" json:"targetTag"`
	TargetPhoto        *Photos                  `gorm:"-" json:"targetPhoto"`

	TagRelatedEntityType string `json:"tagRelatedEntityType"`
}

type EventsDataResponse struct {
	Data  []Events `json:"data"`
	Total int64    `json:"total"`
}

type EventsParams struct {
	TargetUserLogin      string
	TargetCollectionID   uuid.UUID
	TargetItemID         uuid.UUID
	TargetWLID           uuid.UUID
	TargetPostID         string
	TargetTagID          uuid.UUID
	TagRelatedEntityType string
	TargetPhotoID        string
}

type EventsParamsStrings struct {
	TargetUserLogin      string
	TargetCollectionID   string
	TargetItemID         string
	TargetWLID           string
	TargetPostID         string
	TargetTagID          string
	TagRelatedEntityType string
	TargetPhotoID        string
}

type AddEventRequest struct {
	Action          string              `json:"action"`
	EventTargetType string              `json:"eventTargetType"`
	TargetName      string              `json:"targetName"`
	InitiatorLogin  string              `json:"initiatorLogin"`
	Params          EventsParamsStrings `json:"params"`
}

type GetEventsRequest struct {
	Subscriptions    []string `json:"subscriptions"`
	Limit            string   `json:"limit"`
	Offset           string   `json:"offset"`
	Actions          []string `json:"actions"`
	EventTargetTypes []string `json:"eventTargetTypes"`
}
type GetEventsRequestForAll struct {
	Limit            int64    `json:"limit"`
	Offset           int64    `json:"offset"`
	Actions          []string `json:"actions"`
	EventTargetTypes []string `json:"eventTargetTypes"`
}

type GetFilteredEventsRequest struct {
	UserLogin         string   `json:"userLogin"`
	CollectionID      string   `json:"collectionId"`
	CollectionItemID  string   `json:"collectionItemId"`
	WishListItemID    string   `json:"wishListItemId"`
	Limit             string   `json:"limit"`
	Offset            string   `json:"offset"`
	CollectionItemIDs []string `json:"collectionItemIds"`
	TagID             string   `json:"tagId"`
	PhotoID           string   `json:"photoId"`
}
