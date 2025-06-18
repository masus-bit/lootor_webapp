package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepository interface {
	GetUserByLogin(login string) (TokenData, error)
}
type JWTService struct {
	secretKey       []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	userRepo        UserRepository
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type Claims struct {
	Login         string `json:"login"`
	UserName      string `json:"userName"`
	VkId          string `json:"vkId"`
	TelegramId    string `json:"telegramId"`
	Email         string `json:"email"`
	Created       string `json:"created"`
	Likes         int    `json:"likes"`
	Dislikes      int    `json:"dislikes"`
	AvatarUrl     string `json:"avatarUrl"`
	BackgroundUrl string `json:"backgroundUrl"`
	Subscribers   int    `json:"subscribers"`
	Bio           string `json:"bio"`
	City          string `json:"city"`
	IsPremium     bool   `json:"isPremium"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessExp time.Duration, refreshExp time.Duration, userRepo UserRepository) *JWTService {
	return &JWTService{
		secretKey:       []byte(secret),
		accessTokenExp:  accessExp,
		refreshTokenExp: refreshExp,
		userRepo:        userRepo,
	}
}

func (s *JWTService) GenerateTokenPair(user TokenData) (*TokenPair, error) {
	// Access token
	accessToken, err := s.generateToken(user, s.accessTokenExp)
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshToken, err := s.generateToken(user, s.refreshTokenExp)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) generateToken(user TokenData, exp time.Duration) (string, error) {
	claims := Claims{
		Login:         user.GetLogin(),
		UserName:      user.GetUserName(),
		VkId:          user.GetVkId(),
		TelegramId:    user.GetTelegramId(),
		Email:         user.GetEmail(),
		Created:       user.GetCreated(),
		Likes:         user.GetLikes(),
		Dislikes:      user.GetDislikes(),
		AvatarUrl:     user.GetAvatarUrl(),
		BackgroundUrl: user.GetBackgroundUrl(),
		Subscribers:   user.GetSubscribers(),
		Bio:           user.GetBio(),
		City:          user.GetCity(),
		IsPremium:     user.GetIsPremium(),
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
	login         string
	userName      string
	vkId          string
	telegramId    string
	email         string
	created       string
	likes         int
	dislikes      int
	avatarUrl     string
	backgroundUrl string
	subscribers   int
	bio           string
	city          string
	isPremium     bool
}

func (t *tokenData) GetLogin() string         { return t.login }
func (t *tokenData) GetUserName() string      { return t.userName }
func (t *tokenData) GetVkId() string          { return t.vkId }
func (t *tokenData) GetTelegramId() string    { return t.telegramId }
func (t *tokenData) GetEmail() string         { return t.email }
func (t *tokenData) GetCreated() string       { return t.created }
func (t *tokenData) GetLikes() int            { return t.likes }
func (t *tokenData) GetDislikes() int         { return t.dislikes }
func (t *tokenData) GetAvatarUrl() string     { return t.avatarUrl }
func (t *tokenData) GetBackgroundUrl() string { return t.backgroundUrl }
func (t *tokenData) GetSubscribers() int      { return t.subscribers }
func (t *tokenData) GetBio() string           { return t.bio }
func (t *tokenData) GetCity() string          { return t.city }
func (t *tokenData) GetIsPremium() bool       { return t.isPremium }
