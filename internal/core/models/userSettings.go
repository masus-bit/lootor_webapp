package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type UserSettings struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	User      *Users         `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;" json:"user"`
	UserLogin string         `gorm:"type:varchar(255);index"`
	Settings  datatypes.JSON `gorm:"type:jsonb" json:"settings"`
}
