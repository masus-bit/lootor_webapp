package models

import (
	"github.com/lib/pq"
	"time"
)

type Users struct {
	Login             string          `gorm:"primaryKey" json:"login"`
	UserName          string          `json:"userName"`
	City              string          `json:"city"`
	Bio               string          `json:"bio"`
	TmpLogin          bool            `gorm:"default:false" json:"tmpLogin"`
	Password          string          `gorm:"-" json:"-"`
	PasswordHash      string          `gorm:"column:password" json:"-"`
	VkID              string          `json:"vkId"`
	TelegramID        string          `json:"telegramId"`
	YandexID          string          `json:"yandexId"`
	Email             string          `json:"email"`
	Created           string          `json:"created"`
	Likes             int             `json:"likes"`
	Dislikes          int             `json:"dislikes"`
	AvatarURL         string          `json:"avatarUrl"`
	BackgroundURL     string          `json:"backgroundUrl"`
	VerificationToken string          `json:"verificationToken"`
	ResetToken        string          `json:"resetToken"`
	Subscribers       int             `json:"subscribers"`
	SubscribersLogins pq.StringArray  `gorm:"type:text[]" json:"subscribersLogins"`
	Subscriptions     pq.StringArray  `gorm:"type:text[]" json:"subscriptions"`
	TagsSubscriptions pq.StringArray  `gorm:"type:text[]" json:"tagsSubscriptions"`
	WishListItems     []WishListItems `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`
	IsPremium         bool            `gorm:"default:false" json:"isPremium"`
	PremiumSince      time.Time       `json:"premiumSince,omitempty"`
	PremiumUntil      time.Time       `json:"premiumUntil,omitempty"`
	PremiumType       string          `json:"premiumType,omitempty"`
	ReportsCount      int64           `json:"-"`
	DeletedAt         *time.Time      `json:"-" gorm:"column:deleted_at"`
	ProfileName       string          `json:"profileName"`

	Exp         int `json:"exp"`
	SocialScore int `json:"socialScore"`

	Role string `json:"role"`
}

func (u *Users) GetLogin() string               { return u.Login }
func (u *Users) GetUserName() string            { return u.UserName }
func (u *Users) GetVkID() string                { return u.VkID }
func (u *Users) GetTelegramID() string          { return u.TelegramID }
func (u *Users) GetEmail() string               { return u.Email }
func (u *Users) GetCreated() string             { return u.Created }
func (u *Users) GetLikes() int                  { return u.Likes }
func (u *Users) GetDislikes() int               { return u.Dislikes }
func (u *Users) GetAvatarURL() string           { return u.AvatarURL }
func (u *Users) GetBackgroundURL() string       { return u.BackgroundURL }
func (u *Users) GetSubscribers() int            { return u.Subscribers }
func (u *Users) GetBio() string                 { return u.Bio }
func (u *Users) GetCity() string                { return u.City }
func (u *Users) GetIsPremium() bool             { return u.IsPremium }
func (u *Users) GetProfileName() string         { return u.ProfileName }
func (u *Users) GetSubscribersLogins() []string { return u.SubscribersLogins }
func (u *Users) GetSubscriptions() []string {
	return u.Subscriptions
}
func (u *Users) GetTagsSubscriptions() []string {
	return u.TagsSubscriptions
}
func (u *Users) GetRole() string { return u.Role }
func (u *Users) GetExp() int     { return u.Exp }
