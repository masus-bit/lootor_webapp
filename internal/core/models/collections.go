package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Collections struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Created          string         `json:"created"`
	IsPrivate        bool           `json:"isPrivate"`
	BannerURL        string         `json:"bannerUrl"`
	Likes            pq.StringArray `gorm:"type:text[]" json:"likes"`
	Deleted          bool
	Transliteration  string            `json:"transliteration"`
	SubscribersCount int64             `json:"subscribersCount"`
	TotalPrice       int64             `json:"totalPrice"`
	ShippingTotal    int64             `json:"shippingTotal"`
	ShareString      string            `json:"shareString"`
	CollectionItems  []CollectionItems `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;" json:"collectionItems"`
	User             *Users            `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;" json:"user"`
	UserLogin        string            `gorm:"type:varchar(255);index"`
	ReportsCount     int64             `json:"-"`
	CommentsCount    int64             `json:"commentsCount"`
}
