package user

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/utils"
	"time"
)

type Service struct {
	repo       *Repository
	jwtService *auth.JWTService
}

func NewService(repo *Repository, jwtService *auth.JWTService) *Service {
	return &Service{repo: repo, jwtService: jwtService}
}

func (s *Service) GetByLogin(userLogin string) (*DataUserResponse, error) {
	user, err := s.repo.GetByLogin(userLogin)
	if err != nil {
		return nil, err
	}
	var response DataUserResponse
	e := mapstructure.Decode(user, &response)
	if e != nil {
		return nil, e
	}
	return &response, nil
}

func (s *Service) SignIn(dto *SignInRequest) (*SignInResponse, error) {
	user, err := s.repo.GetByEmailWithPassword(dto.Email)
	if err != nil {
		return nil, err
	}
	ok := auth.CheckPasswordHash(dto.Password, user.PasswordHash)
	if !ok {
		return nil, errors.New("incorrect email or password")
	}

	tokens, er := s.jwtService.GenerateTokenPair(user)

	if er != nil {
		return nil, er
	}
	return &SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil

}

func (s *Service) SignUp(dto *SignUpRequest) (*SignUpResponse, error) {
	userByMail, _ := s.repo.GetByEmail(dto.Email)
	if userByMail != nil {
		return nil, errors.New("Email должен быть уникальным")
	}
	userByLogin, _ := s.repo.GetByLogin(dto.Login)
	if userByLogin != nil {
		return nil, errors.New("Логин должен быть уникальным")
	}

	passwordHash, errorHash := auth.HashPassword(dto.Password)
	if errorHash != nil {
		return nil, errorHash
	}
	now := time.Now()
	isoTime := now.Format(time.RFC3339)

	hexString, _ := utils.GenerateRandomString(32)

	user := User{
		Login:             dto.Login,
		Email:             dto.Email,
		UserName:          dto.UserName,
		PasswordHash:      passwordHash,
		Created:           isoTime,
		VerificationToken: hexString,
	}

	err := s.repo.Create(&user)
	if err != nil {
		return nil, err
	}
	return &SignUpResponse{
		Data: "success",
	}, nil

}
