package auth

type TokenData interface {
	GetLogin() string
	GetUserName() string
	GetVkId() string
	GetTelegramId() string
	GetEmail() string
	GetCreated() string
	GetLikes() int
	GetDislikes() int
	GetAvatarUrl() string
	GetBackgroundUrl() string
	GetSubscribers() int
	GetBio() string
	GetCity() string
	GetIsPremium() bool
	GetProfileName() string
	GetSubscribersLogins() []string
	GetSubscriptions() []string
}
