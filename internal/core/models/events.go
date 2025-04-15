package models

import (
	"github.com/google/uuid"
)

type Events struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Date            string
	Action          string `gorm:"type:varchar(100);not null"`
	EventTargetType string `gorm:"type:varchar(50);not null"` // "collection", "user" или "item"
	TargetName      string `gorm:"type:varchar(255)"`

	InitiatorLogin string `gorm:"type:varchar(255);not null"`
	Initiator      Users  `gorm:"foreignKey:InitiatorLogin;references:Login"`

	TargetUserLogin    *string    `gorm:"type:varchar(255)"`
	TargetCollectionID *uuid.UUID `gorm:"type:uuid"`
	TargetItemID       *uuid.UUID `gorm:"type:uuid"`

	TargetUser       *Users           `gorm:"-"`
	TargetCollection *Collections     `gorm:"-"`
	TargetItem       *CollectionItems `gorm:"-"`
}
