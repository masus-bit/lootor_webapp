package dto

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Notifications struct {
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserLogin   string     `json:"userLogin"`
	Type        string     `json:"type"`
	SenderLogin string     `json:"senderLogin"`
	TargetID    string     `json:"targetId"`
	IsRead      bool       `json:"isRead"`
	Date        string     `json:"date"`
	TargetUser  User       `json:"targetUser"`
	SenderUser  User       `json:"senderUser"`
	Target      TargetItem `json:"target"`
	Action      string     `json:"action"`
	Owner       User       `json:"owner"`
}
type NotificationsReadRequest struct {
	IDs []string `json:"ids"`
}

type NotificationsRequest struct {
	Login       string     `json:"login"`
	TargetID    string     `json:"targetId"`
	Type        string     `json:"type"`
	SenderLogin string     `json:"senderLogin"`
	Date        string     `json:"date"`
	TargetUser  User       `json:"targetUser"`
	SenderUser  User       `json:"senderUser"`
	Owner       User       `json:"owner"`
	OwnerLogin  string     `json:"ownerLogin"`
	TargetItem  TargetItem `json:"targetItem"`
	Action      string     `json:"action"`
}

type NotificationsRestRequest struct {
	Login       string     `json:"login"`
	TargetID    string     `json:"targetId"`
	Type        string     `json:"type"`
	SenderLogin string     `json:"senderLogin"`
	Date        string     `json:"date"`
	TargetUser  User       `json:"targetUser"`
	SenderUser  User       `json:"senderUser"`
	TargetItem  TargetItem `json:"targetItem"`
	Action      string     `json:"action"`
}
type User struct {
	Login       string `json:"login"`
	IsPremium   bool   `json:"isPremium"`
	AvatarURL   string `json:"avatarUrl"`
	ProfileName string `json:"profileName"`
}

type TargetItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Transliteration  string `json:"transliteration"`
	TargetType       string `json:"targetType"`
	TargetParentName string `json:"targetParentName"`
}
type NotificationsResponse struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserLogin   string    `json:"userLogin"`
	Type        string    `json:"type"`
	SenderLogin string    `json:"senderLogin"`
	OwnerLogin  string    `json:"ownerLogin"`
	TargetID    string    `json:"targetId"`
	IsRead      bool      `json:"isRead"`
	Date        string    `json:"date"`
	Action      string    `json:"action"`
}

type NotificationsDataResponse struct {
	Data []Notifications `json:"data"`
}
