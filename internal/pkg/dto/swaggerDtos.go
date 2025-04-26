package dto

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
)

type UserResponseSwagger struct {
	Login                   string   `json:"login"`
	UserName                string   `json:"userName"`
	VkId                    string   `json:"vkId"`
	TelegramId              string   `json:"telegramId"`
	Email                   string   `json:"email"`
	Created                 string   `json:"created"`
	Likes                   int      `json:"likes"`
	Dislikes                int      `json:"dislikes"`
	AvatarUrl               string   `json:"avatarUrl"`
	BackgroundUrl           string   `json:"backgroundUrl"`
	Subscribers             int      `json:"subscribers"`
	Subscriptions           []string `json:"subscriptions"`
	CollectionSubscriptions []string `json:"collectionSubscriptions"`
	CanSubscribe            bool     `json:"canSubscribe"`
	CollectionItemsCount    int      `json:"collectionItemsCount" default:"0"`
	CollectionsCount        int      `json:"collectionsCount" default:"0"`
	TotalSum                int      `json:"totalSum" default:"0"`
}

type DataUserResponseSwagger struct {
	Data UserResponseSwagger `json:"data"`
}

type LikeUserSwagger struct {
	IsLike bool `json:"isLike"`
}

type CollectionsCollectionItemsCollectionItems struct {
	CollectionID     uuid.UUID `gorm:"primaryKey;column:collectionId"`
	CollectionItemID uuid.UUID `gorm:"primaryKey;column:collectionItemId"`
}

type TagsCollectionsCollections struct {
	CollectionsID uuid.UUID `gorm:"primaryKey;column:collectionsId"`
	TagsID        uuid.UUID `gorm:"primaryKey;column:tagsId"`
}

type CollectionsSwagger struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Created          string    `json:"created"`
	IsPrivate        bool      `json:"isPrivate"`
	BannerUrl        string    `json:"bannerUrl"`
	Likes            []string  `gorm:"type:text[]" json:"likes"`
	Deleted          bool
	Transliteration  string                   `json:"transliteration"`
	SubscribersCount int64                    `json:"subscribersCount"`
	TotalPrice       int64                    `json:"totalPrice"`
	ShippingTotal    int64                    `json:"shippingTotal"`
	ShareString      string                   `json:"shareString"`
	CollectionItems  []CollectionItemsSwagger `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;" json:"collectionItems"`
	User             *UsersSwagger            `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;" json:"user"`
	Tags             []TagsSwagger            `gorm:"many2many:tags_collections_collections;constraint:OnDelete:CASCADE;" json:"tags"`
	UserLogin        string                   `gorm:"type:varchar(255);index"`
}

type UsersSwagger struct {
	Login                   string   `gorm:"primaryKey" json:"login"`
	UserName                string   `json:"userName"`
	Password                string   `gorm:"-" json:"-"`
	PasswordHash            string   `gorm:"column:password" json:"-"`
	VkId                    string   `json:"vkId"`
	TelegramId              string   `json:"telegramId"`
	Email                   string   `json:"email"`
	Created                 string   `json:"created"`
	Likes                   int      `json:"likes"`
	Dislikes                int      `json:"dislikes"`
	AvatarUrl               string   `json:"avatarUrl"`
	BackgroundUrl           string   `json:"backgroundUrl"`
	VerificationToken       string   `json:"verificationToken"`
	Subscribers             int      `json:"subscribers"`
	Subscriptions           []string `gorm:"type:text[]" json:"subscriptions"`
	CollectionSubscriptions []string `gorm:"type:text[]" json:"collectionSubscriptions"`
}

func (CollectionsCollectionItemsCollectionItems) TableName() string {
	return "collections_collection_items_collection_items"
}

func (TagsCollectionsCollections) TableName() string {
	return "tags_collections_collections"
}

type CollectionsResponseSwagger struct {
	Id                   uuid.UUID                        `json:"id"`
	Name                 string                           `json:"name"`
	Description          string                           `json:"description"`
	Created              string                           `json:"created"`
	IsPrivate            bool                             `json:"isPrivate"`
	BannerUrl            string                           `json:"bannerUrl"`
	Deleted              bool                             `json:"deleted"`
	Transliteration      string                           `json:"transliteration"`
	SubscribersCount     int64                            `json:"subscribersCount"`
	TotalPrice           float64                          `json:"totalPrice"`
	ShippingTotal        float64                          `json:"shippingTotal"`
	ShareString          string                           `json:"shareString"`
	CollectionItems      []CollectionItemsResponseSwagger `json:"collectionItems"`
	User                 UserResponseSwagger              `json:"user"`
	Tags                 []TagsSwagger                    `json:"tags"`
	CollectionItemsCount int64                            `json:"collectionItemsCount"`
	CanSubscribe         bool                             `json:"canSubscribe"`
	LikesCount           int64                            `json:"likesCount"`
	CreatedAt            time.Time                        `json:"createdAt"`
	CanLike              bool                             `json:"canLike"`
	IsOwner              bool                             `json:"isOwner"`
}

type CollectionItemsResponseSwagger struct {
	Id            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Images        []string          `json:"images"`
	PurchaseDate  string            `json:"purchaseDate"`
	PurchasePrice float64           `json:"purchasePrice"`
	Sealed        bool              `json:"sealed"`
	Edition       string            `json:"edition"`
	CopyNumber    []string          `json:"copyNumber"`
	Rating        float64           `json:"rating"`
	ShippingCost  float64           `json:"shippingCost"`
	Entities      []EntitiesSwagger `json:"entities"`
	Platform      *models.Platforms `json:"platform"`
	ItemType      *models.ItemTypes `json:"itemType"`
	Collection    uuid.UUID         `json:"collection"`
	Owner         UsersSwagger      `json:"owner"`
	IsOwner       bool              `json:"isOwner"`
}

type CollectionDataResponseSwagger struct {
	Data CollectionsResponseSwagger `json:"data"`
}

type AllCollectionsDataResponseSwagger struct {
	Data []CollectionsResponseSwagger `json:"data"`
}

type CollectionCreateRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	IsPrivate       bool     `json:"isPrivate"`
	BannerUrl       string   `json:"bannerUrl"`
	Transliteration string   `json:"transliteration"`
	UserLogin       string   `json:"userLogin"`
	Tags            []string `json:"tags"`
}

type EntitiesSwagger struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id             uuid.UUID                `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string                   `json:"name"`
	CollectionItem []CollectionItemsSwagger `gorm:"many2many:entities_collection_item_collection_items;constraint:OnDelete:CASCADE;" json:"collectionItem"`
}

type EntitiesCreateRequestSwagger struct {
	Names []string `json:"names"`
}

type EntitiesDataResponseSwagger struct {
	Data []EntitiesSwagger `json:"data"`
}
type CollectionItemsDataResponseSwagger struct {
	Data CollectionItemsResponseSwagger `json:"data"`
}

type CollectionItemsSwagger struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name            string
	Description     string
	Images          []string `gorm:"type:text[]"`
	PurchaseDate    string
	PurchasePrice   float64
	Sealed          bool
	Edition         string
	CopyNumber      []string `gorm:"type:integer[]"`
	Rating          float64
	Deleted         bool
	ShippingCost    float64
	Transliteration string
	Collections     []CollectionsSwagger `gorm:"many2many:collections_collection_items_collection_items;constraint:OnDelete:CASCADE;"`
	Owner           UsersSwagger         `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`
	Platform        *models.Platforms    `gorm:"foreignKey:PlatformID;references:Id;constraint:OnDelete:SET NULL;"`
	ItemType        *models.ItemTypes    `gorm:"foreignKey:ItemTypeID;references:Id;constraint:OnDelete:SET NULL;"`
	Entities        []EntitiesSwagger    `gorm:"many2many:entities_collection_item_collection_items;constraint:OnDelete:CASCADE;"`
	UserLogin       string               `gorm:"type:varchar(255);index"`
	PlatformID      *uuid.UUID           `gorm:"type:uuid;index"`
	ItemTypeID      *uuid.UUID           `gorm:"type:uuid;index"`
}

type TagsSwagger struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id          uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string               `json:"name"`
	Collections []CollectionsSwagger `gorm:"many2many:collection_tags_collection;constraint:OnDelete:CASCADE;" json:"collections"`
}

type EventsSwagger struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"Id"`
	Date            string    `json:"date"`
	Action          string    `gorm:"type:varchar(100);not null" json:"action"`
	EventTargetType string    `gorm:"type:varchar(50);not null" json:"eventTargetType"` // "collection", "user" или "item"
	TargetName      string    `gorm:"type:varchar(255)" json:"targetName"`

	InitiatorLogin string       `gorm:"type:varchar(255);not null" json:"initiatorLogin"`
	Initiator      UsersSwagger `gorm:"foreignKey:InitiatorLogin;references:Login" json:"initiator"`

	TargetUserLogin    *string    `gorm:"type:varchar(255)" json:"targetUserLogin"`
	TargetCollectionID *uuid.UUID `gorm:"type:uuid" json:"targetCollectionID"`
	TargetItemID       *uuid.UUID `gorm:"type:uuid" json:"targetItemID"`

	TargetUser       *UsersSwagger           `gorm:"-" json:"targetUser"`
	TargetCollection *CollectionsSwagger     `gorm:"-" json:"targetCollection"`
	TargetItem       *CollectionItemsSwagger `gorm:"-" json:"targetItem"`
}

type EventsDataResponseSwagger struct {
	Data []EventsSwagger `json:"data"`
}

type ImagesResponse struct {
	Keys []string `json:"keys"`
}

type WishListItemsSwagger struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id        uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserLogin string       `gorm:"type:varchar(255);index"`
	User      UsersSwagger `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`

	CollectionItemID *uuid.UUID              `gorm:"type:uuid;index"`
	CollectionItem   *CollectionItemsSwagger `gorm:"foreignKey:CollectionItemID;references:Id;constraint:OnDelete:SET NULL;"`

	PurchaseLinks []string `gorm:"type:text[]"`
	Priority      int64
	Notes         string
}

type WishlistDataResponseSwagger struct {
	Data []WishlistDataResponseSwagger `json:"data"`
}
