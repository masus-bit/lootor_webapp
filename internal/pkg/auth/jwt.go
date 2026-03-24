package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/utils"
	"time"
)

type JWTService struct {
	secretKey        []byte
	accessTokenExp   time.Duration
	refreshTokenExp  time.Duration
	userRepo         repositories.UsersRepository
	tagsClient       tagsclient.GRPCTagsClient
	collectionsRepo  repositories.CollectionsRepository
	userSettingsRepo repositories.UserSettingsRepository
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type Claims struct {
	*dto.UserResponseJwt
	Email             string         `json:"email"`
	Role              string         `json:"role"`
	UserSettings      datatypes.JSON `json:"userSettings"`
	PremiumExpireDate string         `json:"premiumExpireDate"`
	Exp               int            `json:"experience"`
	jwt.RegisteredClaims
}

func NewJWTService(
	secret string,
	accessExp time.Duration,
	refreshExp time.Duration,
	userRepo repositories.UsersRepository,
	tagsClient tagsclient.GRPCTagsClient,
	collectionsRepo repositories.CollectionsRepository,
	userSettingsRepo repositories.UserSettingsRepository,
) *JWTService {
	return &JWTService{
		secretKey:        []byte(secret),
		accessTokenExp:   accessExp,
		refreshTokenExp:  refreshExp,
		userRepo:         userRepo,
		tagsClient:       tagsClient,
		collectionsRepo:  collectionsRepo,
		userSettingsRepo: userSettingsRepo,
	}
}

func (s *JWTService) GenerateTokenPair(user TokenData) (*TokenPair, error) {
	userData, _ := s.userRepo.GetUserByLogin(user.GetLogin())
	subscribersLogins, err := s.userRepo.GetForSubs(userData.SubscribersLogins)
	if err != nil {
		return nil, err
	}

	subscriptions, err := s.userRepo.GetForSubs(userData.Subscriptions)
	if err != nil {
		return nil, err
	}

	collectionsCount, err := s.collectionsRepo.GetCollectionsCount(user.GetLogin(), true)
	if err != nil {
		return nil, err
	}
	userSettings, err := s.userSettingsRepo.FindRecordByUserLogin(user.GetLogin())
	if err != nil {
		return nil, err
	}
	var premiumExpireDate string
	if !userData.IsPremium {
		premiumExpireDate = ""
	} else {
		premiumExpireDate = utils.DateZFormat(userData.PremiumUntil.String())
	}

	subTags, _ := s.tagsClient.GetTagsByIDs(
		context.Background(),
		&microservices.GetTagsByIDsRequest{Ids: userData.TagsSubscriptions},
	)
	var resultTags []dto.SubTags
	resultTags = make([]dto.SubTags, 0)
	if len(subTags.GetTags()) > 0 || subTags != nil {
		for _, tag := range subTags.GetTags() {
			resultTags = append(
				resultTags, dto.SubTags{
					ID:   tag.Id,
					Name: tag.Name,
					Slug: tag.Slug,
				},
			)
		}
	}
	// Access token
	accessToken, err := s.generateToken(
		user,
		s.accessTokenExp,
		subscribersLogins,
		subscriptions,
		resultTags,
		collectionsCount,
		userSettings.Settings,
		premiumExpireDate,
	)
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshToken, err := s.generateToken(
		user,
		s.refreshTokenExp,
		subscribersLogins,
		subscriptions,
		resultTags,
		collectionsCount,
		userSettings.Settings,
		premiumExpireDate,
	)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) generateToken(
	user TokenData,
	exp time.Duration,
	subLogins, subs []dto.SubUsers,
	tagsSubs []dto.SubTags,
	collectionsCount int64,
	userSettings datatypes.JSON,
	premiumExpire string,
) (string, error) {
	claims := Claims{
		//&dto.UserResponse{
		//	Login:             user.GetLogin(),
		//	UserName:          user.GetUserName(),
		//	VkID:              user.GetVkID(),
		//	TelegramID:        user.GetTelegramID(),
		//	Email:             user.GetEmail(),
		//	Created:           user.GetCreated(),
		//	Likes:             user.GetLikes(),
		//	Dislikes:          user.GetDislikes(),
		//	AvatarURL:         user.GetAvatarURL(),
		//	BackgroundURL:     user.GetBackgroundURL(),
		//	Subscribers:       user.GetSubscribers(),
		//	Bio:               user.GetBio(),
		//	City:              user.GetCity(),
		//	IsPremium:         user.GetIsPremium(),
		//	ProfileName:       user.GetProfileName(),
		//	SubscribersLogins: subLogins,
		//	Subscriptions:     subs,
		//	TagsSubscriptions: tagsSubs,
		//	CollectionsCount:  int(collectionsCount)
		//},
		UserResponseJwt: &dto.UserResponseJwt{
			Login:                     user.GetLogin(),
			UserName:                  user.GetUserName(),
			VkID:                      user.GetVkID(),
			TelegramID:                user.GetTelegramID(),
			Created:                   user.GetCreated(),
			Likes:                     user.GetLikes(),
			Dislikes:                  user.GetDislikes(),
			AvatarURL:                 user.GetAvatarURL(),
			BackgroundURL:             user.GetBackgroundURL(),
			Subscribers:               user.GetSubscribers(),
			Bio:                       user.GetBio(),
			City:                      user.GetCity(),
			IsPremium:                 user.GetIsPremium(),
			ProfileName:               user.GetProfileName(),
			SubscribersExtended:       subLogins,
			SubscriptionsExtended:     subs,
			TagsSubscriptionsExtended: tagsSubs,
			CollectionsCount:          int(collectionsCount),
		},
		Email:             user.GetEmail(),
		UserSettings:      userSettings,
		PremiumExpireDate: premiumExpire,
		Role:              user.GetRole(),
		Exp:               user.GetExp(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ParseToken валидирует JWT токен и возвращает userID
func (s *JWTService) ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return s.secretKey, nil
		},
	)

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", errors.New("invalid token")
}

func (s *JWTService) GetRole(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return s.secretKey, nil
		},
	)

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Role, nil
	}

	return "", errors.New("invalid token")
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
func (s *JWTService) RenewTokenPair(refreshToken string) (*TokenPair, error) {
	token, err := jwt.ParseWithClaims(
		refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return s.secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("refresh token expired")
	}

	userData, err := s.userRepo.GetUserByLogin(claims.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	return s.GenerateTokenPair(userData)
}

type tokenData struct {
	login             string
	userName          string
	vkId              string
	telegramId        string
	email             string
	created           string
	likes             int
	dislikes          int
	avatarUrl         string
	backgroundUrl     string
	subscribers       int
	bio               string
	city              string
	isPremium         bool
	profileName       string
	subscribersLogins []dto.SubUsers
	subscriptions     []dto.SubUsers
	tagsSubscriptions []dto.SubTags
	role              string
	collectionsCount  int
	experience        int
}
