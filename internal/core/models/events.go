package models

import (
	"github.com/google/uuid"
	"lootor/internal/pkg/types"
)

type Events struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date            string    `json:"date"`
	Action          string    `gorm:"type:varchar(100);not null" json:"action"`
	EventTargetType string    `gorm:"type:varchar(50);not null" json:"eventTargetType"`
	TargetName      string    `gorm:"type:varchar(255)" json:"targetName"`

	InitiatorLogin string `gorm:"type:varchar(255);not null" json:"initiatorLogin"`
	Initiator      *Users `gorm:"foreignKey:InitiatorLogin;references:Login" json:"initiator"`

	TargetUserLogin      *string    `gorm:"type:varchar(255)" json:"targetUserLogin"`
	TargetCollectionID   *uuid.UUID `gorm:"type:uuid" json:"targetCollectionId"`
	TargetItemID         *uuid.UUID `gorm:"type:uuid" json:"targetItemId"`
	TargetWishListItemID *uuid.UUID `gorm:"type:uuid" json:"targetWishListItemId"`

	TargetUser         *SubUsers              `gorm:"-" json:"targetUser"`
	TargetCollection   *types.CommonShortType `gorm:"-" json:"targetCollection"`
	TargetItem         *types.CommonShortType `gorm:"-" json:"targetItem"`
	TargetWishListItem *types.CommonShortType `gorm:"-" json:"targetWishListItem"`
}

type EventsDataResponse struct {
	Data  []Events `json:"data"`
	Total int64    `json:"total"`
}
