package dto

import (
	"github.com/lib/pq"
	"lootor/internal/core/models"
)

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
	DonateTotal string `json:"donateTotal"`
}

type UserResponse struct {
	Login                 string                 `json:"login"`
	UserName              string                 `json:"userName"`
	City                  string                 `json:"city"`
	Bio                   string                 `json:"bio"`
	VkID                  string                 `json:"vkId"`
	TelegramID            string                 `json:"telegramId"`
	Created               string                 `json:"created"`
	Likes                 int                    `json:"likes"`
	Dislikes              int                    `json:"dislikes"`
	AvatarURL             string                 `json:"avatarUrl"`
	BackgroundURL         string                 `json:"backgroundUrl"`
	Subscribers           int                    `json:"subscribers"`
	Subscriptions         []string               `json:"subscriptions"`
	TagsSubscriptions     []string               `json:"tagsSubscriptions"`
	SubscribersLogins     []string               `json:"subscribersLogins"`
	SubscriptionsExtended []SubUsers             `json:"subscriptionsExtended" mapstructure:"-"`
	SubscribersExtended   []SubUsers             `json:"subscribersExtended" mapstructure:"-"`
	CanSubscribe          bool                   `json:"canSubscribe"`
	CollectionItemsCount  int                    `json:"collectionItemsCount" default:"0"`
	CollectionsCount      int                    `json:"collectionsCount" default:"0"`
	TotalSum              int                    `json:"totalSum" default:"0"`
	WishListItems         []models.WishListItems `json:"wishListItems"`
	IsPremium             bool                   `json:"isPremium"`
	ShippingTotal         int                    `json:"shippingTotal"`
	ProfileName           string                 `json:"profileName"`
	Exp                   int                    `json:"exp"`
	SocialScore           int                    `json:"socialScore"`
	TotalDonations        string                 `json:"totalDonations"`
	PostCount             int                    `json:"postCount"`
	PremiumExpireDate     string                 `json:"premiumExpireDate"`
}

type UserResponseForSingleUser struct {
	Login                string                 `json:"login"`
	UserName             string                 `json:"userName"`
	City                 string                 `json:"city"`
	Bio                  string                 `json:"bio"`
	VkID                 string                 `json:"vkId"`
	TelegramID           string                 `json:"telegramId"`
	Created              string                 `json:"created"`
	Likes                int                    `json:"likes"`
	Dislikes             int                    `json:"dislikes"`
	AvatarURL            string                 `json:"avatarUrl"`
	BackgroundURL        string                 `json:"backgroundUrl"`
	SubscribersLogins    []SubUsers             `json:"subscribersLogins" mapstructure:"-"`
	Subscribers          int                    `json:"subscribers"`
	Subscriptions        []SubUsers             `json:"subscriptions" mapstructure:"-"`
	TagsSubscriptions    []SubTags              `json:"tagsSubscriptions" mapstructure:"-"`
	CanSubscribe         bool                   `json:"canSubscribe"`
	CollectionItemsCount int                    `json:"collectionItemsCount" default:"0"`
	CollectionsCount     int                    `json:"collectionsCount" default:"0"`
	TotalSum             int                    `json:"totalSum" default:"0"`
	WishListItems        []models.WishListItems `json:"wishListItems"`
	IsPremium            bool                   `json:"isPremium"`
	ShippingTotal        int                    `json:"shippingTotal"`
	ProfileName          string                 `json:"profileName"`
	Exp                  int                    `json:"exp"`
	SocialScore          int                    `json:"socialScore"`
	TotalDonations       string                 `json:"totalDonations"`
	PostCount            int                    `json:"postCount"`
	PremiumExpireDate    string                 `json:"premiumExpireDate"`
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

type YandexSuggestRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

type YandexUserInfo struct {
	Login    string `json:"login"`
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	PSUID    string `json:"psuid"`

	DefaultEmail string   `json:"default_email,omitempty"`
	Emails       []string `json:"emails,omitempty"`

	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	RealName    string `json:"real_name,omitempty"`

	DefaultAvatarID string `json:"default_avatar_id,omitempty"`
	IsAvatarEmpty   bool   `json:"is_avatar_empty,omitempty"`
}
