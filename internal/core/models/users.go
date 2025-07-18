package models

import (
	"github.com/lib/pq"
	"time"
)

type Users struct {
	Login                   string          `gorm:"primaryKey" json:"login"`
	UserName                string          `json:"userName"`
	City                    string          `json:"city"`
	Bio                     string          `json:"bio"`
	TmpLogin                bool            `gorm:"default:false" json:"tmpLogin"`
	Password                string          `gorm:"-" json:"-"`
	PasswordHash            string          `gorm:"column:password" json:"-"`
	VkId                    string          `json:"vkId"`
	TelegramId              string          `json:"telegramId"`
	Email                   string          `json:"email"`
	Created                 string          `json:"created"`
	Likes                   int             `json:"likes"`
	Dislikes                int             `json:"dislikes"`
	AvatarUrl               string          `json:"avatarUrl"`
	BackgroundUrl           string          `json:"backgroundUrl"`
	VerificationToken       string          `json:"verificationToken"`
	ResetToken              string          `json:"resetToken"`
	Subscribers             int             `json:"subscribers"`
	Subscriptions           pq.StringArray  `gorm:"type:text[]" json:"subscriptions"`
	CollectionSubscriptions pq.StringArray  `gorm:"type:text[]" json:"collectionSubscriptions"`
	WishListItems           []WishListItems `gorm:"foreignKey:UserLogin;references:Login;constraint:OnDelete:CASCADE;"`
	IsPremium               bool            `gorm:"default:false" json:"isPremium"`
	PremiumSince            time.Time       `json:"premiumSince,omitempty"`
	PremiumUntil            time.Time       `json:"premiumUntil,omitempty"`
	PremiumType             string          `json:"premiumType,omitempty"`
	ReportsCount            int64           `json:"-"`
	CreatorRating           int             `json:"creatorRating"`
	DeletedAt               *time.Time      `json:"-" gorm:"column:deleted_at"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Login    string `json:"login" validate:"required"`
	UserName string `json:"userName" validate:"required"`
	City     string `json:"city"`
	Bio      string `json:"bio"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangeRatingRequest struct {
	IsLike bool `json:"isLike"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

type ChangeLoginRequest struct {
	Login string `json:"login" validate:"required,min=3"`
}

type ChangePasswordReset struct {
	Password string `json:"password" validate:"required,min=8"`
	Token    string `json:"token" validate:"required"`
}

type UserResponse struct {
	Login                   string          `json:"login"`
	UserName                string          `json:"userName"`
	City                    string          `json:"city"`
	Bio                     string          `json:"bio"`
	VkId                    string          `json:"vkId"`
	TelegramId              string          `json:"telegramId"`
	Email                   string          `json:"email"`
	Created                 string          `json:"created"`
	Likes                   int             `json:"likes"`
	Dislikes                int             `json:"dislikes"`
	AvatarUrl               string          `json:"avatarUrl"`
	BackgroundUrl           string          `json:"backgroundUrl"`
	Subscribers             int             `json:"subscribers"`
	Subscriptions           pq.StringArray  `json:"subscriptions"`
	CollectionSubscriptions pq.StringArray  `json:"collectionSubscriptions"`
	CanSubscribe            bool            `json:"canSubscribe"`
	CollectionItemsCount    int             `json:"collectionItemsCount" default:"0"`
	CollectionsCount        int             `json:"collectionsCount" default:"0"`
	TotalSum                int             `json:"totalSum" default:"0"`
	WishListItems           []WishListItems `json:"wishListItems"`
	IsPremium               bool            `json:"isPremium"`
	CreatorRating           int             `json:"creatorRating"`
	ShippingTotal           int             `json:"shippingTotal"`
}

type UserRequestUpdate struct {
	UserName      *string `json:"userName,omitempty"`
	Email         *string `json:"email,omitempty"`
	AvatarUrl     *string `json:"avatarUrl,omitempty"`
	BackgroundUrl *string `json:"backgroundUrl,omitempty"`
	City          *string `json:"city,omitempty"`
	Bio           *string `json:"bio,omitempty"`
}

type UserRequestUpdateFirstTime struct {
	Login         *string `json:"login,omitempty"`
	UserName      *string `json:"userName,omitempty"`
	Email         *string `json:"email,omitempty"`
	AvatarUrl     *string `json:"avatarUrl,omitempty"`
	BackgroundUrl *string `json:"backgroundUrl,omitempty"`
	City          *string `json:"city,omitempty"`
	Bio           *string `json:"bio,omitempty"`
}

type UserResponseForCollection struct {
	Login         string         `json:"login"`
	UserName      string         `json:"userName"`
	City          string         `json:"city"`
	Bio           string         `json:"bio"`
	Email         string         `json:"email"`
	AvatarUrl     string         `json:"avatarUrl"`
	Subscriptions pq.StringArray `json:"subscriptions"`
}

type DataUserResponse struct {
	Data UserResponse `json:"data"`
}

type DataUsersResponse struct {
	Data []UserResponse `json:"data"`
}

type SignInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignInResponseWithTmpLogin struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TmpLogin     bool   `json:"tmpLogin"`
}

type RenewTokensRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type SignUpResponse struct {
	Data string `json:"data"`
}

type VkOauthRequest struct {
	Code          string `json:"code" validate:"required"`
	CodeVerifier  string `json:"codeVerifier" validate:"required"`
	DeviceId      string `json:"deviceId" validate:"required"`
	State         string `json:"state" validate:"required"`
	CodeChallenge string `json:"codeChallenge" validate:"required"`
}

type VkAuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	RedirectUri  string `json:"redirect_uri"`
	DeviceId     string `json:"device_id"`
	State        string `json:"state"`
}

type VkAuthGetTokenData struct {
	AccessToken string `json:"access_token"`
}

type VkAuthGetUserInfo struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	AvatarUrl string `json:"photo_200"`
}

type VkAuthGetUserInfoResponse struct {
	Response []VkAuthGetUserInfo `json:"response"`
}

type TelegramOauthRequest struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoUrl  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func (u *Users) GetLogin() string         { return u.Login }
func (u *Users) GetUserName() string      { return u.UserName }
func (u *Users) GetVkId() string          { return u.VkId }
func (u *Users) GetTelegramId() string    { return u.TelegramId }
func (u *Users) GetEmail() string         { return u.Email }
func (u *Users) GetCreated() string       { return u.Created }
func (u *Users) GetLikes() int            { return u.Likes }
func (u *Users) GetDislikes() int         { return u.Dislikes }
func (u *Users) GetAvatarUrl() string     { return u.AvatarUrl }
func (u *Users) GetBackgroundUrl() string { return u.BackgroundUrl }
func (u *Users) GetSubscribers() int      { return u.Subscribers }
func (u *Users) GetBio() string           { return u.Bio }
func (u *Users) GetCity() string          { return u.City }
func (u *Users) GetIsPremium() bool       { return u.IsPremium }
