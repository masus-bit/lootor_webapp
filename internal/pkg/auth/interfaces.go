package auth

type TokenData interface {
	GetLogin() string
	GetUserName() string
	GetVkID() string
	GetTelegramID() string
	GetEmail() string
	GetCreated() string
	GetLikes() int
	GetDislikes() int
	GetAvatarURL() string
	GetBackgroundURL() string
	GetSubscribers() int
	GetBio() string
	GetCity() string
	GetIsPremium() bool
	GetProfileName() string
	GetSubscribersLogins() []string
	GetSubscriptions() []string
	GetRole() string
	GetExp() int
}
