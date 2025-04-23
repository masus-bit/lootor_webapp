package models

import (
	"github.com/google/uuid"
)

type Events struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"Id"`
	Date            string    `json:"date"`
	Action          string    `gorm:"type:varchar(100);not null" json:"action"`
	EventTargetType string    `gorm:"type:varchar(50);not null" json:"eventTargetType"`
	TargetName      string    `gorm:"type:varchar(255)" json:"targetName"`

	InitiatorLogin string `gorm:"type:varchar(255);not null" json:"initiatorLogin"`
	Initiator      Users  `gorm:"foreignKey:InitiatorLogin;references:Login" json:"initiator"`

	TargetUserLogin      *string    `gorm:"type:varchar(255)" json:"targetUserLogin"`
	TargetCollectionID   *uuid.UUID `gorm:"type:uuid" json:"targetCollectionID"`
	TargetItemID         *uuid.UUID `gorm:"type:uuid" json:"targetItemID"`
	TargetWishListItemID *uuid.UUID `gorm:"type:uuid" json:"targetWishListItemID"`

	TargetUser         *Users           `gorm:"-" json:"targetUser"`
	TargetCollection   *Collections     `gorm:"-" json:"targetCollection"`
	TargetItem         *CollectionItems `gorm:"-" json:"targetItem"`
	TargetWishListItem *WishListItems   `gorm:"-" json:"targetWishListItem"`
}

type EventsDataResponse struct {
	Data []Events `json:"data"`
}
