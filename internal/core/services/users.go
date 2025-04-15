package services

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/utils"
	"slices"
	"time"
)

type UserService struct {
	repo        *repositories.UsersRepository
	jwtService  *auth.JWTService
	ciRepo      *repositories.CiRepository
	mailService *mail.MailService
}

func NewUserService(repo *repositories.UsersRepository, jwtService *auth.JWTService, ciRepo *repositories.CiRepository, mailService *mail.MailService) *UserService {
	return &UserService{repo: repo, jwtService: jwtService, ciRepo: ciRepo, mailService: mailService}
}

func (s *UserService) GetByLogin(userLogin string, authUser string, isAuthenticated bool) (*models.DataUserResponse, error) {
	dbUser, err := s.repo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var authorizedUser *models.Users
	fmt.Println(authUser, 2)
	if authUser != "" && isAuthenticated {
		authorizedUser, _ = s.repo.GetUserByLogin(authUser)
	}
	var response models.UserResponse

	if userLogin == authUser {

		collectionItemsCount, _ := s.ciRepo.GetCountByUserLogin(userLogin)
		sum, _ := s.ciRepo.SumByUserLogin(userLogin)

		response.TotalSum = int(sum)
		response.CollectionItemsCount = int(collectionItemsCount)

		e := mapstructure.Decode(dbUser, &response)
		fmt.Print(dbUser, response)
		if e != nil {
			return nil, e
		}
		return &models.DataUserResponse{Data: response}, nil
	}

	var subArray []string

	if authorizedUser != nil {
		subArray = authorizedUser.Subscriptions
		canSubscribe := !slices.Contains(subArray, userLogin)
		response.CanSubscribe = canSubscribe
	}

	e := mapstructure.Decode(dbUser, &response)
	fmt.Println(dbUser, response)
	if e != nil {
		return nil, e
	}

	return &models.DataUserResponse{Data: response}, nil
}

func (s *UserService) ChangeRating(isLike bool, login string) (*dto.CommonResponse, error) {
	existsUser, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}
	if isLike {
		existsUser.Likes = existsUser.Likes + 1
	} else {
		existsUser.Dislikes = existsUser.Dislikes + 1
	}
	_, err = s.repo.UpdateUser(*existsUser, *existsUser)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) ChangePassword(password string, login string, authUserLogin string) (*dto.CommonResponse, error) {
	if authUserLogin != login {
		return nil, errors.New("ошибка смены пароля")
	}

	existsUser, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	existsUser.PasswordHash, _ = auth.HashPassword(password)
	_, err = s.repo.UpdateUser(*existsUser, *existsUser)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil

}

func (s *UserService) Subscribe(targetUserLogin string, authUserLogin string, isSubscribe bool) (*dto.CommonResponse, error) {
	subscriber, err := s.repo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}

	subscriptionTargetUser, err := s.repo.GetUserByLogin(targetUserLogin)
	if err != nil {
		return nil, err
	}

	if isSubscribe {
		subscriber.Subscriptions = append(subscriber.Subscriptions, targetUserLogin)
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers + 1
	} else {
		subscriber.Subscriptions = utils.RemoveByValue(subscriber.Subscriptions, targetUserLogin)
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers - 1
	}
	_, err = s.repo.UpdateUser(*subscriber, *subscriber)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.UpdateUser(*subscriptionTargetUser, *subscriptionTargetUser)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) SignIn(dto *models.SignInRequest) (*models.SignInResponse, error) {
	dbUser, err := s.repo.GetByEmailWithPassword(dto.Email)
	if err != nil {
		return nil, err
	}
	ok := auth.CheckPasswordHash(dto.Password, dbUser.PasswordHash)
	if !ok {
		return nil, errors.New("incorrect email or password")
	}

	tokens, er := s.jwtService.GenerateTokenPair(dbUser)

	if er != nil {
		return nil, er
	}
	return &models.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil

}

func (s *UserService) SignUp(dto *models.SignUpRequest) (*models.SignUpResponse, error) {
	userByMail, _ := s.repo.GetByEmail(dto.Email)
	if userByMail != nil {
		return nil, errors.New("Email должен быть уникальным")
	}
	userByLogin, _ := s.repo.GetUserByLogin(dto.Login)
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

	dbUser := models.Users{
		Login:             dto.Login,
		Email:             dto.Email,
		UserName:          dto.UserName,
		PasswordHash:      passwordHash,
		Created:           isoTime,
		VerificationToken: hexString,
	}

	err := s.repo.CreateUser(&dbUser)
	if err != nil {
		return nil, err
	}

	go func() {
		if s.mailService == nil {
			fmt.Println("mailService is not initialized")
			return
		}
		err := s.mailService.SendConfirmationEmail(dto.Email, hexString)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return &models.SignUpResponse{
		Data: "success",
	}, nil

}
