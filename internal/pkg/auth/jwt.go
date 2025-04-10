package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type JWTService struct {
	secretKey       []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
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
	Subscribers   string `json:"subscribers"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessExp time.Duration, refreshExp time.Duration) *JWTService {
	return &JWTService{
		secretKey:       []byte(secret),
		accessTokenExp:  accessExp,
		refreshTokenExp: refreshExp,
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
