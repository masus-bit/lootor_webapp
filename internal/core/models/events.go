package models

import (
	"github.com/google/uuid"
	"lootor/internal/pkg/types"
)

type Events struct {
	Id              string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
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
	TargetPostId         string `json:"targetPostId"`

	TargetUser         *SubUsers                        `gorm:"-" json:"targetUser"`
	TargetCollection   *types.CommonShortTypeCollection `gorm:"-" json:"targetCollection"`
	TargetItem         *types.CommonShortTypeItem       `gorm:"-" json:"targetItem"`
	TargetWishListItem *types.CommonShortType           `gorm:"-" json:"targetWishListItem"`
	TargetPost         *types.CommonShortTypePost       `gorm:"-" json:"targetPost"`
}

type EventsDataResponse struct {
	Data  []Events `json:"data"`
	Total int64    `json:"total"`
}

type EventsParams struct {
	TargetUserLogin    string
	TargetCollectionID uuid.UUID
	TargetItemID       uuid.UUID
	TargetWLID         uuid.UUID
	TargetPostID       string
}

type EventsParamsStrings struct {
	TargetUserLogin    string
	TargetCollectionID string
	TargetItemID       string
	TargetWLID         string
	TargetPostID       string
}

type AddEventRequest struct {
	Action          string              `json:"action"`
	EventTargetType string              `json:"eventTargetType"`
	TargetName      string              `json:"targetName"`
	InitiatorLogin  string              `json:"initiatorLogin"`
	Params          EventsParamsStrings `json:"params"`
}

type GetEventsRequest struct {
	Subscriptions []string `json:"subscriptions"`
	Limit         string   `json:"limit"`
	Offset        string   `json:"offset"`
}

type GetFilteredEventsRequest struct {
	UserLogin        string `json:"userLogin"`
	CollectionId     string `json:"collectionId"`
	CollectionItemId string `json:"collectionItemId"`
	WishListItemId   string `json:"wishListItemId"`
	Limit            string `json:"limit"`
	Offset           string `json:"offset"`
}
