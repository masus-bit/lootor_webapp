package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"time"
)

type JWTService struct {
	secretKey       []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	userRepo        repositories.UsersRepository
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type Claims struct {
	Login             string            `json:"login"`
	UserName          string            `json:"userName"`
	VkID              string            `json:"vkId"`
	TelegramID        string            `json:"telegramId"`
	Email             string            `json:"email"`
	Created           string            `json:"created"`
	Likes             int               `json:"likes"`
	Dislikes          int               `json:"dislikes"`
	AvatarURL         string            `json:"avatarUrl"`
	BackgroundURL     string            `json:"backgroundUrl"`
	Subscribers       int               `json:"subscribers"`
	Bio               string            `json:"bio"`
	City              string            `json:"city"`
	IsPremium         bool              `json:"isPremium"`
	ProfileName       string            `json:"profileName"`
	SubscribersLogins []models.SubUsers `json:"subscribersLogins" mapstructure:"-"`
	Subscriptions     []models.SubUsers `json:"subscriptions" mapstructure:"-"`
	Role              string            `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessExp time.Duration, refreshExp time.Duration, userRepo repositories.UsersRepository) *JWTService {
	return &JWTService{
		secretKey:       []byte(secret),
		accessTokenExp:  accessExp,
		refreshTokenExp: refreshExp,
		userRepo:        userRepo,
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
	// Access token
	accessToken, err := s.generateToken(user, s.accessTokenExp, subscribersLogins, subscriptions)
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshToken, err := s.generateToken(user, s.refreshTokenExp, subscribersLogins, subscriptions)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) generateToken(user TokenData, exp time.Duration, subLogins, subs []models.SubUsers) (string, error) {
	claims := Claims{
		Login:             user.GetLogin(),
		UserName:          user.GetUserName(),
		VkID:              user.GetVkID(),
		TelegramID:        user.GetTelegramID(),
		Email:             user.GetEmail(),
		Created:           user.GetCreated(),
		Likes:             user.GetLikes(),
		Dislikes:          user.GetDislikes(),
		AvatarURL:         user.GetAvatarURL(),
		BackgroundURL:     user.GetBackgroundURL(),
		Subscribers:       user.GetSubscribers(),
		Bio:               user.GetBio(),
		City:              user.GetCity(),
		IsPremium:         user.GetIsPremium(),
		ProfileName:       user.GetProfileName(),
		SubscribersLogins: subLogins,
		Subscriptions:     subs,
		Role:              user.GetRole(),
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", errors.New("invalid token")
}

func (s *JWTService) GetRole(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

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
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

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
	subscribersLogins []models.SubUsers
	subscriptions     []models.SubUsers
	role              string
}
