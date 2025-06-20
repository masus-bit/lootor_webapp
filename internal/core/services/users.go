package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/utils"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"
)

type UserService struct {
	repo        *repositories.UsersRepository
	jwtService  *auth.JWTService
	ciRepo      *repositories.CiRepository
	mailService *mail.MailService
	evRepo      *repositories.EventsRepository
}

func NewUserService(repo *repositories.UsersRepository, jwtService *auth.JWTService, ciRepo *repositories.CiRepository, mailService *mail.MailService, evRepo *repositories.EventsRepository) *UserService {
	_ = godotenv.Load()
	return &UserService{repo: repo, jwtService: jwtService, ciRepo: ciRepo, mailService: mailService, evRepo: evRepo}
}

func (s *UserService) GetByLogin(userLogin string, authUser string, isAuthenticated bool) (*models.DataUserResponse, error) {
	dbUser, err := s.repo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var authorizedUser *models.Users
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
		var lowerCasedUsers []string
		for _, userInArray := range subArray {
			lowerCasedUsers = append(lowerCasedUsers, strings.ToLower(userInArray))
		}
		canSubscribe := !slices.Contains(lowerCasedUsers, strings.ToLower(userLogin))
		response.CanSubscribe = canSubscribe
	}

	e := mapstructure.Decode(dbUser, &response)
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
	_, err = s.repo.UpdateUser(existsUser, *existsUser)
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
	_, err = s.repo.UpdateUser(existsUser, *existsUser)
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
		subscriber.Subscriptions = append(subscriber.Subscriptions, strings.ToLower(targetUserLogin))
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers + 1
	} else {
		subscriber.Subscriptions = utils.RemoveByValue(subscriber.Subscriptions, strings.ToLower(targetUserLogin))
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers - 1
	}
	_, err = s.repo.UpdateUser(subscriber, *subscriber)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.UpdateUser(subscriptionTargetUser, *subscriptionTargetUser)
	if err != nil {
		return nil, err
	}

	err = s.evRepo.AddEvent(authUserLogin, utils.EventActionSubscribe, utils.EventTargetUser, targetUserLogin, nil, nil, nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) SignIn(dto *models.SignInRequest) (*models.SignInResponse, error) {
	dbUser, err := s.repo.GetByEmailWithPassword(dto.Email)
	if err != nil {
		return nil, errors.New("incorrect email or password")
	}
	if dbUser.VerificationToken != "" {
		return nil, errors.New("user is not verified")
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

func (s *UserService) SignUp(dto *models.SignUpRequest, ctx context.Context) (*models.SignUpResponse, error) {
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
		Bio:               dto.Bio,
		City:              dto.City,
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

func (s *UserService) Verification(token string) (*models.SignInResponse, error) {
	var tokens *auth.TokenPair
	dbUser, err := s.repo.GetByVerificationToken(token)
	if err != nil {
		return nil, err
	}
	var newUser *models.Users
	err = mapstructure.Decode(dbUser, &newUser)
	if err != nil {
		return nil, err
	}

	newUser.VerificationToken = ""
	_, err = s.repo.UpdateUser(dbUser, *newUser)
	if err != nil {
		return nil, err
	}
	tokens, err = s.jwtService.GenerateTokenPair(dbUser)
	if err != nil {
		return nil, err
	}
	return &models.SignInResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *UserService) RefreshTokens(refreshToken string) (*models.SignInResponse, error) {

	tokens, err := s.jwtService.RenewTokenPair(refreshToken)
	if err != nil {
		return nil, err
	}
	return &models.SignInResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

type VkAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	Email       string `json:"email,omitempty"`
}

func (s *UserService) VkOauth(dto *models.VkOauthRequest) (*models.SignInResponseWithTmpLogin, error) {

	bodyRequest := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     os.Getenv("VK_CLIENT_ID"),
		"code":          dto.Code,
		"code_verifier": dto.CodeVerifier,
		"redirect_uri":  os.Getenv("FRONTEND_URL"),
		"device_id":     dto.DeviceId,
		"state":         dto.State,
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	response, err := utils.SendRequest[models.VkAuthGetTokenData](utils.RequestOptions{Method: "POST", URL: "https://id.vk.com/oauth2/auth", Headers: headers, QueryParams: map[string]string{}, Body: bodyRequest, File: []byte{}, BasicAuth: nil})

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	userInfo, err := s.getUserInfo(response.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	newId := fmt.Sprintf("%d", userInfo.Id)

	existsUser, _ := s.repo.GetByVkId(newId)

	fmt.Println(fmt.Sprintf("user%+v", existsUser))
	if existsUser != nil && !existsUser.TmpLogin {

		if existsUser.VerificationToken != "" {
			return nil, errors.New("user is not verified")
		}

		tokens, err := s.jwtService.GenerateTokenPair(existsUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &models.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     false,
		}, nil
	} else if existsUser != nil && existsUser.TmpLogin {

		tokens, err := s.jwtService.GenerateTokenPair(existsUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &models.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     true,
		}, nil
	}

	passwordHash, err := auth.HashPassword("vk" + newId)
	if err != nil {
		return nil, fmt.Errorf("password hash error: %w", err)
	}

	newUser := models.Users{
		VkId:         newId,
		Login:        userInfo.FirstName + "@" + newId,
		Email:        userInfo.Email,
		UserName:     userInfo.FirstName,
		AvatarUrl:    userInfo.AvatarUrl,
		Created:      time.Now().Format(time.RFC3339),
		TmpLogin:     true,
		PasswordHash: passwordHash,
	}

	fmt.Println(1)

	err2 := s.repo.CreateUser(&newUser)
	if err2 != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	fmt.Println(2)
	user, err := s.repo.GetByVkId(newUser.VkId)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	fmt.Println(3)
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("token generation error: %w", err)
	}
	fmt.Println(4)
	return &models.SignInResponseWithTmpLogin{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TmpLogin:     true,
	}, nil
}
func (s *UserService) getUserInfo(accessToken string) (*models.VkAuthGetUserInfo, error) {
	baseUrl := "https://api.vk.com/method/users.get"

	parameters := map[string]string{
		"access_token": accessToken,
		"v":            "5.131",
		"fields":       "email,photo_200",
	}
	response, err := utils.SendRequest[struct {
		Response []models.VkAuthGetUserInfo `json:"response"`
	}](utils.RequestOptions{Method: "GET", URL: baseUrl, Headers: map[string]string{}, Body: map[string]string{}, QueryParams: parameters, File: []byte{}, BasicAuth: nil})

	if err != nil {
		return nil, err
	}

	return &response.Response[0], nil
}

func (s *UserService) TelegramOauth(dto *models.TelegramOauthRequest) (*models.SignInResponseWithTmpLogin, error) {
	if dto.LastName == "" {

	}
	data := []string{
		fmt.Sprintf("id=%d", dto.Id),
		fmt.Sprintf("username=%s", dto.Username),
		fmt.Sprintf("auth_date=%d", dto.AuthDate),
	}

	if dto.LastName != "" {
		data = append(data, fmt.Sprintf("last_name=%s", dto.LastName))
	}
	if dto.FirstName != "" {
		data = append(data, fmt.Sprintf("first_name=%s", dto.FirstName))
	}
	if dto.PhotoUrl != "" {
		data = append(data, fmt.Sprintf("photo_url=%s", dto.PhotoUrl))
	}
	hash := dto.Hash

	sort.Strings(data)
	dataCheckString := strings.Join(data, "\n")
	hSecretKey := sha256.New()
	hSecretKey.Write([]byte(os.Getenv("TELEGRAM_BOT_TOKEN")))

	h := hmac.New(sha256.New, hSecretKey.Sum(nil))

	h.Write([]byte(dataCheckString))

	hashBytes := h.Sum(nil)

	computedHash := hex.EncodeToString(hashBytes)
	if hash != computedHash {
		return nil, fmt.Errorf("Проблемы с хэшем")
	}

	newId := fmt.Sprintf("%d", dto.Id)

	existUser, _ := s.repo.GetByTgId(newId)

	if existUser != nil && !existUser.TmpLogin {
		if existUser.VerificationToken != "" {
			return nil, errors.New("user is not verified")
		}
		tokens, err := s.jwtService.GenerateTokenPair(existUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &models.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     false,
		}, nil
	} else if existUser != nil && existUser.TmpLogin {
		tokens, err := s.jwtService.GenerateTokenPair(existUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &models.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     true,
		}, nil
	}

	passwordHash, err := auth.HashPassword("tg" + newId + dto.Username)
	if err != nil {
		return nil, fmt.Errorf("password hash error: %w", err)
	}

	newUser := models.Users{
		TelegramId:   newId,
		Login:        dto.Username + newId,
		UserName:     dto.FirstName + " " + dto.LastName,
		AvatarUrl:    dto.PhotoUrl,
		Created:      time.Now().Format(time.RFC3339),
		PasswordHash: passwordHash,
		TmpLogin:     true,
	}

	err = s.repo.CreateUser(&newUser)
	if err != nil {
		return nil, fmt.Errorf("user creation error: %w", err)
	}

	user, err := s.repo.GetByTgId(newId)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("token generation error: %w", err)
	}

	return &models.SignInResponseWithTmpLogin{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TmpLogin:     true,
	}, nil
}

func (s *UserService) ResetPassword(email string) (*dto.CommonResponse, error) {
	existsUser, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	hexString, _ := utils.GenerateRandomString(32)

	existsUser.ResetToken = hexString

	_, err = s.repo.UpdateUser(existsUser, *existsUser)
	if err != nil {
		return nil, err
	}

	go func() {
		if s.mailService == nil {
			fmt.Println("mailService is not initialized")
			return
		}
		err := s.mailService.SendResetPasswordEmail(email, hexString)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) ChangeResetPassword(request *models.ChangePasswordReset) (*dto.CommonResponse, error) {
	existsUser, err := s.repo.GetByResetToken(request.Token)
	if err != nil {
		return nil, err
	}
	existsUser.PasswordHash, _ = auth.HashPassword(request.Password)
	existsUser.ResetToken = ""
	_, err = s.repo.UpdateUser(existsUser, *existsUser)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) UpdateUser(user *models.UserRequestUpdate, login string) (*models.SignInResponse, error) {
	if login == "" {
		return nil, errors.New("not authorized")
	}

	existsUser, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	userByEmail, _ := s.repo.GetByEmail(*user.Email)

	if userByEmail != nil && userByEmail.Email == *user.Email {
		return nil, errors.New("пользователь с таким email уже существует")

	}

	dst := reflect.ValueOf(existsUser).Elem()
	src := reflect.ValueOf(user).Elem()

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		if !field.IsNil() {
			fieldName := src.Type().Field(i).Name
			dstField := dst.FieldByName(fieldName)
			if dstField.IsValid() {
				dstField.Set(field.Elem())
			}
		}
	}

	updatedUser, err := s.repo.UpdateUserFull(existsUser)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtService.GenerateTokenPair(updatedUser)
	if err != nil {
		return nil, fmt.Errorf("token generation error: %w", err)
	}

	return &models.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) UpdateOnlyOnce(data *models.UserRequestUpdateFirstTime, targetUserLogin string) (*dto.CommonResponse, error) {
	if data == nil {
		return nil, errors.New("update data is nil")
	}
	emailUser, err := s.repo.GetByEmail(*data.Email)
	if err != nil {
		fmt.Errorf("failed to check user existence: %w", err)
	}

	if emailUser != nil && emailUser.Email == *data.Email {
		return nil, errors.New("email already exists")

	}
	if data.Login == nil || data.UserName == nil || data.Email == nil {
		return nil, errors.New("required fields are missing")
	}

	existsUser, err := s.repo.GetUserByLogin(targetUserLogin)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existsUser == nil {
		return nil, errors.New("user with this login not exists")
	}

	hexString, _ := utils.GenerateRandomString(32)

	existsUser.Login = *data.Login
	existsUser.UserName = *data.UserName
	existsUser.Email = *data.Email
	existsUser.AvatarUrl = *data.AvatarUrl
	existsUser.BackgroundUrl = *data.BackgroundUrl
	existsUser.City = *data.City
	existsUser.Bio = *data.Bio
	existsUser.TmpLogin = false
	existsUser.VerificationToken = hexString

	if *data.Email != "" {

		go func() {
			if s.mailService == nil {
				fmt.Println("mailService is not initialized")
				return
			}
			err := s.mailService.SendConfirmationEmail(*data.Email, hexString)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	_, err = s.repo.UpdateLogin(targetUserLogin, existsUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) ActivatePremium(userLogin string, months int, years int, premiumType string) error {
	existsUser, _ := s.repo.GetUserByLogin(userLogin)

	updatedUser := existsUser
	updatedUser.IsPremium = true
	updatedUser.PremiumSince = time.Now()
	updatedUser.PremiumUntil = time.Now().AddDate(years, months, 0)
	updatedUser.PremiumType = premiumType

	ok, err := s.repo.ActivatePremium(existsUser, updatedUser)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("error pizda")
	}
	return nil
}
func (s *UserService) ActivateTestPremium(userLogin string, duration time.Duration, premiumType string) error {
	existsUser, _ := s.repo.GetUserByLogin(userLogin)

	updatedUser := existsUser
	updatedUser.IsPremium = true
	updatedUser.PremiumSince = time.Now()
	updatedUser.PremiumUntil = time.Now().Add(duration)
	updatedUser.PremiumType = premiumType

	ok, err := s.repo.ActivatePremium(existsUser, updatedUser)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("error pizda")
	}
	return nil
}

func (s *UserService) CheckPremiumStatus(userLogin string) (bool, error) {
	return s.repo.CheckPremiumStatus(userLogin)
}

func (s *UserService) CheckExpiredSubscriptions() {
	err := s.repo.DeactivatePremium()
	if err != nil {
		fmt.Println(err)
	}
}
