package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"io"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/utils"
	"net/http"
	"net/url"
	"os"
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
		canSubscribe := !slices.Contains(subArray, userLogin)
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
		subscriber.Subscriptions = append(subscriber.Subscriptions, targetUserLogin)
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers + 1
	} else {
		subscriber.Subscriptions = utils.RemoveByValue(subscriber.Subscriptions, targetUserLogin)
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

	err = s.evRepo.AddEvent(authUserLogin, utils.EventActionSubscribe, utils.EventTargetUser, targetUserLogin, nil, nil, nil)
	if err != nil {
		fmt.Println(err)
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

func (s *UserService) VkOauth(dto *models.VkOauthRequest) (*models.SignInResponse, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", os.Getenv("VK_CLIENT_ID"))
	params.Add("code", dto.Code)
	params.Add("code_verifier", dto.CodeVerifier)
	params.Add("redirect_uri", os.Getenv("FRONTEND_URL"))
	params.Add("device_id", dto.DeviceId)
	params.Add("state", dto.State)

	req, err := http.NewRequest(
		"POST",
		"https://id.vk.com/oauth2/auth",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vk auth error: %s", string(body))

	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp models.VkAuthGetTokenData
	//fmt.Printf("VK OAUTH API RAW RESPONSE: %s\n", string(body))

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}
	userInfo, err := s.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	newId := fmt.Sprintf("%d", userInfo.Id)

	existsUser, _ := s.repo.GetByVkId(userInfo.FirstName + "@" + newId)

	if existsUser != nil {
		tokens, err := s.jwtService.GenerateTokenPair(existsUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &models.SignInResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
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
		PasswordHash: passwordHash,
	}

	if err := s.repo.CreateUser(&newUser); err != nil {
		return nil, fmt.Errorf("user creation error: %w", err)
	}

	user, err := s.repo.GetByVkId(newUser.VkId)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("token generation error: %w", err)
	}

	return &models.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
func (s *UserService) getUserInfo(accessToken string) (*models.VkAuthGetUserInfo, error) {
	baseUrl := "https://api.vk.com/method/users.get"

	params := url.Values{}
	params.Add("access_token", accessToken)
	params.Add("v", "5.131")
	params.Add("fields", "email,photo_200")

	resp, err := http.Get(baseUrl + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("vk api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vk api error: %s", string(body))
	}

	var response struct {
		Response []models.VkAuthGetUserInfo `json:"response"`
	}

	body, err := io.ReadAll(resp.Body)
	//fmt.Printf("VK API RAW RESPONSE: %s\n", string(body))

	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to read user info: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	if len(response.Response) == 0 {
		return nil, errors.New("empty user info response")
	}

	return &response.Response[0], nil
}

func (s *UserService) TelegramOauth(dto *models.TelegramOauthRequest) (*models.SignInResponse, error) {

	data := []string{
		fmt.Sprintf("id=%s", dto.Id),
		fmt.Sprintf("first_name=%s", dto.FirstName),
		fmt.Sprintf("last_name=%s", dto.LastName),
		fmt.Sprintf("username=%s", dto.Username),
		fmt.Sprintf("photo_url=%s", dto.PhotoUrl),
		fmt.Sprintf("auth_date=%s", dto.AuthDate),
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
	fmt.Println(hash, computedHash, "HASHI")
	if hash != computedHash {
		return nil, fmt.Errorf("Проблемы с хэшем")
	}

	newId := fmt.Sprintf("%d", dto.Id)

	existUser, err := s.repo.GetByTgId(newId)

	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	if existUser != nil {

		tokens, err := s.jwtService.GenerateTokenPair(existUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &models.SignInResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}, nil
	}

	passwordHash, err := auth.HashPassword("tg" + newId + dto.Username)
	if err != nil {
		return nil, fmt.Errorf("password hash error: %w", err)
	}

	newUser := models.Users{
		TelegramId:   newId,
		Login:        dto.Username,
		UserName:     dto.FirstName + " " + dto.LastName,
		AvatarUrl:    dto.PhotoUrl,
		Created:      time.Now().Format(time.RFC3339),
		PasswordHash: passwordHash,
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

	return &models.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
