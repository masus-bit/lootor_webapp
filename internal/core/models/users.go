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

type SignUpRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	Login       string `json:"login" validate:"required"`
	UserName    string `json:"userName" validate:"required"`
	City        string `json:"city"`
	Bio         string `json:"bio"`
	ProfileName string `json:"profileName"`
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

type SubUsers struct {
	Login       string `json:"login"`
	AvatarURL   string `json:"avatarUrl" gorm:"column:avatar_url"`
	ProfileName string `json:"profileName" gorm:"column:profile_name"`
	IsPremium   bool   `json:"isPremium" gorm:"column:is_premium"`
}

type UserResponse struct {
	Login                 string          `json:"login"`
	UserName              string          `json:"userName"`
	City                  string          `json:"city"`
	Bio                   string          `json:"bio"`
	VkID                  string          `json:"vkId"`
	TelegramID            string          `json:"telegramId"`
	Created               string          `json:"created"`
	Likes                 int             `json:"likes"`
	Dislikes              int             `json:"dislikes"`
	AvatarURL             string          `json:"avatarUrl"`
	BackgroundURL         string          `json:"backgroundUrl"`
	Subscribers           int             `json:"subscribers"`
	Subscriptions         []string        `json:"subscriptions"`
	TagsSubscriptions     []string        `json:"tagsSubscriptions"`
	SubscribersLogins     []string        `json:"subscribersLogins"`
	SubscriptionsExtended []SubUsers      `json:"subscriptionsExtended" mapstructure:"-"`
	SubscribersExtended   []SubUsers      `json:"subscribersExtended" mapstructure:"-"`
	CanSubscribe          bool            `json:"canSubscribe"`
	CollectionItemsCount  int             `json:"collectionItemsCount" default:"0"`
	CollectionsCount      int             `json:"collectionsCount" default:"0"`
	TotalSum              int             `json:"totalSum" default:"0"`
	WishListItems         []WishListItems `json:"wishListItems"`
	IsPremium             bool            `json:"isPremium"`
	ShippingTotal         int             `json:"shippingTotal"`
	ProfileName           string          `json:"profileName"`
	Exp                   int             `json:"exp"`
	SocialScore           int             `json:"socialScore"`
	TotalDonations        string          `json:"totalDonations"`
	PostCount             int             `json:"postCount"`
}

type UserResponseForSingleUser struct {
	Login                string          `json:"login"`
	UserName             string          `json:"userName"`
	City                 string          `json:"city"`
	Bio                  string          `json:"bio"`
	VkID                 string          `json:"vkId"`
	TelegramID           string          `json:"telegramId"`
	Created              string          `json:"created"`
	Likes                int             `json:"likes"`
	Dislikes             int             `json:"dislikes"`
	AvatarURL            string          `json:"avatarUrl"`
	BackgroundURL        string          `json:"backgroundUrl"`
	SubscribersLogins    []SubUsers      `json:"subscribersLogins" mapstructure:"-"`
	Subscribers          int             `json:"subscribers"`
	Subscriptions        []SubUsers      `json:"subscriptions" mapstructure:"-"`
	TagsSubscriptions    []SubTags       `json:"tagsSubscriptions" mapstructure:"-"`
	CanSubscribe         bool            `json:"canSubscribe"`
	CollectionItemsCount int             `json:"collectionItemsCount" default:"0"`
	CollectionsCount     int             `json:"collectionsCount" default:"0"`
	TotalSum             int             `json:"totalSum" default:"0"`
	WishListItems        []WishListItems `json:"wishListItems"`
	IsPremium            bool            `json:"isPremium"`
	ShippingTotal        int             `json:"shippingTotal"`
	ProfileName          string          `json:"profileName"`
	Exp                  int             `json:"exp"`
	SocialScore          int             `json:"socialScore"`
	TotalDonations       string          `json:"totalDonations"`
	PostCount            int             `json:"postCount"`
}

type UserRequestUpdate struct {
	UserName      *string `json:"userName,omitempty"`
	Email         *string `json:"email,omitempty"`
	AvatarURL     *string `json:"avatarUrl,omitempty"`
	BackgroundURL *string `json:"backgroundUrl,omitempty"`
	City          *string `json:"city,omitempty"`
	Bio           *string `json:"bio,omitempty"`
	ProfileName   *string `json:"profileName"`
}

type UserRequestUpdateFirstTime struct {
	Login         *string `json:"login,omitempty"`
	UserName      *string `json:"userName,omitempty"`
	Email         *string `json:"email,omitempty"`
	AvatarURL     *string `json:"avatarUrl,omitempty"`
	BackgroundURL *string `json:"backgroundUrl,omitempty"`
	City          *string `json:"city,omitempty"`
	Bio           *string `json:"bio,omitempty"`
	ProfileName   *string `json:"profileName"`
}

type UserResponseForCollection struct {
	Login         string         `json:"login"`
	UserName      string         `json:"userName"`
	City          string         `json:"city"`
	Bio           string         `json:"bio"`
	AvatarURL     string         `json:"avatarUrl"`
	Subscriptions pq.StringArray `json:"subscriptions"`
	ProfileName   string         `json:"profileName"`
}

type DataUserResponse struct {
	Data UserResponse `json:"data"`
}

type DataUserResponseForSingleUser struct {
	Data UserResponseForSingleUser `json:"data"`
}

type DataUsersResponse struct {
	Data  []UserResponse `json:"data"`
	Total int64          `json:"total"`
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
	DeviceID      string `json:"deviceId" validate:"required"`
	State         string `json:"state" validate:"required"`
	CodeChallenge string `json:"codeChallenge" validate:"required"`
}

type VkAuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	RedirectURI  string `json:"redirect_uri"`
	DeviceID     string `json:"device_id"`
	State        string `json:"state"`
}

type VkAuthGetTokenData struct {
	AccessToken string `json:"access_token"`
}

type VkAuthGetUserInfo struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	AvatarURL string `json:"photo_200"`
}

type VkAuthGetUserInfoResponse struct {
	Response []VkAuthGetUserInfo `json:"response"`
}

type TelegramOauthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

type SubTags struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
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
func (u *Users) GetRole() string { return u.Role }
