package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserSettingsRequest struct {
	Settings datatypes.JSON `gorm:"type:jsonb" json:"settings"`
}

type UserSettingsResponse struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserLogin string         `gorm:"type:varchar(255);index" json:"userLogin"`
	Settings  datatypes.JSON `gorm:"type:jsonb" json:"settings"`
}
