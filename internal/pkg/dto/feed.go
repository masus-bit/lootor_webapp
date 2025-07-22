package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Feed struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date string    `json:"date"`
	Type string    `gorm:"type:varchar(255)" json:"type"`

	Content datatypes.JSON `gorm:"type:jsonb" json:"content"`
}

type FeedDataResponse struct {
	Data  []Feed `json:"data"`
	Total int64  `json:"total"`
}

type FeedResponse struct {
	Data Feed `json:"data"`
}

type FeedRequest struct {
	Date string `json:"date"`
	Type string `json:"type"`

	Content datatypes.JSON `json:"content"`
}

type FeedRequestSwag struct {
	Date string `json:"date"`
	Type string `json:"type"`

	Content map[string]interface{} `json:"content"`
}

type FeedSwag struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date string    `json:"date"`
	Type string    `gorm:"type:varchar(255)" json:"type"`

	Content map[string]interface{} `json:"content"`
}

type FeedDataResponseSwag struct {
	Data []FeedSwag `json:"data"`
}

type FeedResponseSwag struct {
	Data FeedSwag `json:"data"`
}
