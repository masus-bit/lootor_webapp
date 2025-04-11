package user

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/models"
	"lootor/internal/modules/collectionItem"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"slices"
	"time"
)

type Service struct {
	repo       *Repository
	jwtService *auth.JWTService
	ciRepo     *collectionItem.Repository
}

func NewService(repo *Repository, jwtService *auth.JWTService, ciRepo *collectionItem.Repository) *Service {
	return &Service{repo: repo, jwtService: jwtService, ciRepo: ciRepo}
}

func (s *Service) GetByLogin(userLogin string, authUser string, isAuthenticated bool) (*models.DataUserResponse, error) {
	user, err := s.repo.GetByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var authorizedUser *models.User

	if authUser != "" && isAuthenticated {
		authorizedUser, _ = s.repo.GetByLogin(authUser)
	}
	var response models.DataUserResponse
	if userLogin == authUser {
		//TODO itemsCount и тд

		collectionItemsCount, _ := s.ciRepo.GetCountByUserLogin(userLogin)
		sum, _ := s.ciRepo.SumByUserLogin(userLogin)

		response.Data.TotalSum = int(sum)
		response.Data.CollectionItemsCount = int(collectionItemsCount)

		e := mapstructure.Decode(user, &response)
		if e != nil {
			return nil, e
		}
		return &response, nil
	}
	var subArray []string

	if authorizedUser != nil {
		subArray = authorizedUser.Subscriptions
		canSubscribe := !slices.Contains(subArray, userLogin)
		response.Data.CanSubscribe = canSubscribe
	}

	e := mapstructure.Decode(user, &response)

	if e != nil {
		return nil, e
	}

	return &response, nil
}

func (s *Service) ChangeRating(isLike bool, login string) (*dto.CommonResponse, error) {
	existsUser, err := s.repo.GetByLogin(login)
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

func (s *Service) ChangePassword(password string, login string, authUserLogin string) (*dto.CommonResponse, error) {
	if authUserLogin != login {
		return nil, errors.New("ошибка смены пароля")
	}

	existsUser, err := s.repo.GetByLogin(login)
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

func (s *Service) Subscribe(targetUserLogin string, authUserLogin string, isSubscribe bool) (*dto.CommonResponse, error) {
	subscriber, err := s.repo.GetByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}

	subscriptionTargetUser, err := s.repo.GetByLogin(targetUserLogin)
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

func (s *Service) SignIn(dto *models.SignInRequest) (*models.SignInResponse, error) {
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
	return &models.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil

}

func (s *Service) SignUp(dto *models.SignUpRequest) (*models.SignUpResponse, error) {
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

	user := models.User{
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
	return &models.SignUpResponse{
		Data: "success",
	}, nil

}
