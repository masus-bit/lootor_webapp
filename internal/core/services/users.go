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
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/achievementsclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/utils"
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type UserService struct {
	repo                 *repositories.UsersRepository
	jwtService           *auth.JWTService
	ciRepo               *repositories.CiRepository
	mailService          *mail.PostService
	eventsService        *EventsService
	notificationsService *NotificationsService
	collectionRepo       *repositories.CollectionsRepository
	postsService         *PostsService
	tagsClient           *tagsclient.GRPCTagsClient
	achClient            *achievementsclient.GRPCAchievementsClient
}

func NewUserService(
	repo *repositories.UsersRepository,
	jwtService *auth.JWTService,
	ciRepo *repositories.CiRepository,
	mailService *mail.PostService,
	eventsService *EventsService,
	notificationsService *NotificationsService,
	collectionRepo *repositories.CollectionsRepository,
	postsService *PostsService,
	tagsClient *tagsclient.GRPCTagsClient,
	achClient *achievementsclient.GRPCAchievementsClient,
) *UserService {
	_ = godotenv.Load()
	return &UserService{
		repo:                 repo,
		jwtService:           jwtService,
		ciRepo:               ciRepo,
		mailService:          mailService,
		eventsService:        eventsService,
		notificationsService: notificationsService,
		collectionRepo:       collectionRepo,
		postsService:         postsService,
		tagsClient:           tagsClient,
		achClient:            achClient,
	}
}

func (s *UserService) GetByLogin(
	userLogin string,
	authUser string,
	isAuthenticated bool,
) (*dto.DataUserResponseForSingleUser, error) {
	dbUser, err := s.repo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var authorizedUser *models.Users
	if authUser != "" && isAuthenticated {
		authorizedUser, _ = s.repo.GetUserByLogin(authUser)
	}
	var response dto.UserResponseForSingleUser
	postsCount := s.postsService.GetCount(context.Background(), userLogin)

	if userLogin == authUser {
		totalDonations, _ := s.repo.GetDonationsTotal(userLogin)
		collectionItemsCount, _ := s.ciRepo.GetCountByUserLogin(userLogin)
		collectionsCount, _ := s.collectionRepo.GetCollectionsCount(userLogin)
		sum, _ := s.ciRepo.SumByUserLogin(userLogin)
		shippingTotal, _ := s.ciRepo.ShippingSumByUserLogin(userLogin)
		donationsString := strconv.FormatFloat(totalDonations, 'f', -1, 64)
		response.TotalSum = int(sum)
		response.CollectionItemsCount = int(collectionItemsCount)
		response.CollectionsCount = int(collectionsCount)
		response.ShippingTotal = int(shippingTotal)
		response.TotalDonations = donationsString
		response.PostCount = int(postsCount)

		e := mapstructure.Decode(dbUser, &response)

		subscribersLogins, err := s.repo.GetForSubs(dbUser.SubscribersLogins)
		if err != nil {
			return nil, err
		}
		response.SubscribersLogins = subscribersLogins
		subscriptions, err := s.repo.GetForSubs(dbUser.Subscriptions)
		if err != nil {
			return nil, err
		}

		subTags, _ := s.tagsClient.GetTagsByIDs(
			context.Background(),
			&microservices.GetTagsByIDsRequest{Ids: dbUser.TagsSubscriptions},
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
		response.TagsSubscriptions = resultTags

		response.Subscriptions = subscriptions
		if e != nil {
			return nil, e
		}
		return &dto.DataUserResponseForSingleUser{Data: response}, nil
	}

	var subArray []string
	response.CanSubscribe = true
	if authorizedUser != nil {
		subArray = authorizedUser.Subscriptions
		var lowerCasedUsers []string
		for _, userInArray := range subArray {
			lowerCasedUsers = append(lowerCasedUsers, strings.ToLower(userInArray))
		}
		canSubscribe := !slices.Contains(lowerCasedUsers, strings.ToLower(userLogin))
		response.CanSubscribe = canSubscribe
	}

	collectionItemsCount, _ := s.ciRepo.GetCountByUserLogin(userLogin)
	collectionsCount, _ := s.collectionRepo.GetCollectionsCount(userLogin)
	sum, _ := s.ciRepo.SumByUserLogin(userLogin)
	shippingTotal, _ := s.ciRepo.ShippingSumByUserLogin(userLogin)
	totalDonations, _ := s.repo.GetDonationsTotal(userLogin)
	donationsString := strconv.FormatFloat(totalDonations, 'f', -1, 64)

	response.TotalSum = int(sum)
	response.CollectionItemsCount = int(collectionItemsCount)
	response.CollectionsCount = int(collectionsCount)
	response.ShippingTotal = int(shippingTotal)
	response.TotalDonations = donationsString
	response.PostCount = int(postsCount)
	e := mapstructure.Decode(dbUser, &response)
	response.VkID = ""
	response.TelegramID = ""
	if e != nil {
		return nil, e
	}

	subscribersLogins, err := s.repo.GetForSubs(dbUser.SubscribersLogins)
	if err != nil {
		return nil, err
	}

	response.SubscribersLogins = subscribersLogins
	subscriptions, err := s.repo.GetForSubs(dbUser.Subscriptions)
	if err != nil {
		return nil, err
	}
	subTags, _ := s.tagsClient.GetTagsByIDs(
		context.Background(),
		&microservices.GetTagsByIDsRequest{Ids: dbUser.TagsSubscriptions},
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
	response.TagsSubscriptions = resultTags
	response.Subscriptions = subscriptions

	return &dto.DataUserResponseForSingleUser{Data: response}, nil
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

func (s *UserService) Subscribe(targetUserLogin string, authUserLogin string, isSubscribe bool) (
	*dto.CommonResponse,
	error,
) {
	var subscriber *models.Users
	var subscriptionTargetUser *models.Users
	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)
	go func() {
		defer wg.Done()
		subscriber, err1 = s.repo.GetUserByLogin(authUserLogin)
	}()
	go func() {
		defer wg.Done()
		subscriptionTargetUser, err2 = s.repo.GetUserByLogin(targetUserLogin)
	}()
	wg.Wait()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	if isSubscribe {
		subscriber.Subscriptions = append(subscriber.Subscriptions, strings.ToLower(targetUserLogin))
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers + 1
		subscriptionTargetUser.SubscribersLogins = append(
			subscriptionTargetUser.SubscribersLogins,
			strings.ToLower(authUserLogin),
		)
		err := s.repo.IncrementExperience(targetUserLogin, utils.UserSelfSubExp)

		subscriptionTargetUser.Exp += utils.UserSelfSubExp

		go func() {
			err = s.eventsService.AddEvent(
				authUserLogin,
				utils.EventActionSubscribe,
				utils.EventTargetUser,
				targetUserLogin,
				&dto.EventsParams{TargetUserLogin: targetUserLogin},
			)
			if err != nil {
				fmt.Println(err)
			}

			subscribers, _ := s.repo.GetTotalSubscribers(targetUserLogin)

			xp, level := utils.GetAchievementSubscribersData(subscribers + 1)

			err = utils.AddAchievement(s.achClient, utils.AchieveSubscribers, targetUserLogin, level, xp, subscribers+1)
			err = s.repo.IncrementExperience(targetUserLogin, int(xp))
		}()
		var target = dto.TargetItem{
			ID:              targetUserLogin,
			Name:            subscriptionTargetUser.ProfileName,
			Transliteration: "",
			TargetType:      "user",
		}
		go func() {
			err = s.notificationsService.SendNotification(
				context.Background(), &dto.NotificationsRequest{
					Login:       targetUserLogin,
					TargetID:    targetUserLogin,
					SenderLogin: authUserLogin,
					Type:        utils.NotificationTypeUser,
					Action:      utils.NotificationActionSubscribe,
					Date:        time.Now().Format(time.RFC3339),
					OwnerLogin:  targetUserLogin,
				}, &target,
			)

		}()
	} else {
		subscriber.Subscriptions = utils.RemoveByValue(subscriber.Subscriptions, strings.ToLower(targetUserLogin))
		subscriptionTargetUser.Subscribers = subscriptionTargetUser.Subscribers - 1
		subscriptionTargetUser.SubscribersLogins = utils.RemoveByValue(
			subscriptionTargetUser.SubscribersLogins,
			strings.ToLower(authUserLogin),
		)
		go func() {
			_ = s.notificationsService.DeleteNotification(context.Background(), targetUserLogin, authUserLogin)
			err := s.eventsService.AddEvent(
				authUserLogin,
				utils.EventActionUnsubscribe,
				utils.EventTargetUser,
				targetUserLogin,
				&dto.EventsParams{TargetUserLogin: targetUserLogin},
			)
			if err != nil {
				fmt.Println(err)
			}
		}()
		subscriptionTargetUser.Exp -= utils.UserSelfSubExp
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err1 = s.repo.UpdateUser(subscriber, *subscriber)
	}()
	go func() {
		defer wg.Done()
		_, err2 = s.repo.UpdateUser(subscriptionTargetUser, *subscriptionTargetUser)
	}()
	wg.Wait()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserService) SignIn(req *dto.SignInRequest) (*dto.SignInResponse, error) {
	dbUser, err := s.repo.GetByEmailWithPassword(req.Email)
	if err != nil {
		return nil, errors.New("incorrect email or password")
	}
	if dbUser.VerificationToken != "" {
		return nil, errors.New("user is not verified")
	}
	ok := auth.CheckPasswordHash(req.Password, dbUser.PasswordHash)
	if !ok {
		return nil, errors.New("incorrect email or password")
	}

	tokens, er := s.jwtService.GenerateTokenPair(dbUser)

	if er != nil {
		return nil, er
	}
	return &dto.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil

}

func (s *UserService) SignUp(req *dto.SignUpRequest, ctx context.Context) (*dto.SignUpResponse, error) {
	userByMail, _ := s.repo.GetByEmailForSignUp(req.Email)
	if userByMail != nil {
		return nil, errors.New("email должен быть уникальным")
	}
	userByLogin, _ := s.repo.GetUserByLoginForSignUp(req.Login)
	if userByLogin != nil {
		return nil, errors.New("логин должен быть уникальным")
	}

	passwordHash, errorHash := auth.HashPassword(req.Password)
	if errorHash != nil {
		return nil, errorHash
	}
	now := time.Now()
	isoTime := now.Format(time.RFC3339)
	hexString, _ := utils.GenerateRandomString(32)

	dbUser := models.Users{
		Login:             req.Login,
		Email:             req.Email,
		UserName:          req.UserName,
		PasswordHash:      passwordHash,
		Created:           isoTime,
		VerificationToken: hexString,
		Bio:               req.Bio,
		City:              req.City,
		ProfileName:       req.ProfileName,
		Role:              "user",
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
		err = s.mailService.SendConfirmationEmail(req.Email, hexString)
		if err != nil {
			fmt.Println(err)
		}
		_, err = s.achClient.AddZeroAchievements(
			context.Background(),
			&microservices.GetUserAchievementsRequest{UserLogin: req.Login},
		)
		err = utils.AddAchievement(s.achClient, utils.AchieveBetaTester, req.Login, 1, utils.XPBetaTester, 0)
		err = s.repo.IncrementExperience(req.Login, utils.XPBetaTester)
	}()

	return &dto.SignUpResponse{
		Data: "success",
	}, nil
}

func (s *UserService) Verification(token string) (*dto.SignInResponse, error) {
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
	return &dto.SignInResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *UserService) RefreshTokens(refreshToken string) (*dto.SignInResponse, error) {

	tokens, err := s.jwtService.RenewTokenPair(refreshToken)
	if err != nil {
		return nil, err
	}
	return &dto.SignInResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

type VkAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	Email       string `json:"email,omitempty"`
}

func (s *UserService) VkOauth(req *dto.VkOauthRequest) (*dto.SignInResponseWithTmpLogin, error) {

	bodyRequest := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     os.Getenv("VK_CLIENT_ID"),
		"code":          req.Code,
		"code_verifier": req.CodeVerifier,
		"redirect_uri":  os.Getenv("FRONTEND_URL"),
		"device_id":     req.DeviceID,
		"state":         req.State,
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	response, err := utils.SendRequest[dto.VkAuthGetTokenData](
		utils.RequestOptions{
			Method:      "POST",
			URL:         "https://id.vk.com/oauth2/auth",
			Headers:     headers,
			QueryParams: map[string]string{},
			Body:        bodyRequest,
			File:        []byte{},
			BasicAuth:   nil,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	userInfo, err := s.getUserInfo(response.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	newId := fmt.Sprintf("%d", userInfo.ID)

	existsUser, _ := s.repo.GetByVkId(newId)

	if existsUser != nil && !existsUser.TmpLogin {

		if existsUser.VerificationToken != "" {
			return nil, errors.New("user is not verified")
		}

		tokens, err := s.jwtService.GenerateTokenPair(existsUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &dto.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     false,
		}, nil
	} else if existsUser != nil && existsUser.TmpLogin {

		tokens, err := s.jwtService.GenerateTokenPair(existsUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &dto.SignInResponseWithTmpLogin{
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
		VkID:         newId,
		Login:        userInfo.FirstName + "@" + newId,
		Email:        userInfo.Email,
		UserName:     userInfo.FirstName,
		AvatarURL:    userInfo.AvatarURL,
		Created:      time.Now().Format(time.RFC3339),
		TmpLogin:     true,
		PasswordHash: passwordHash,
		ProfileName:  userInfo.FirstName,
	}

	if userInfo.AvatarURL != "" {
		err = s.repo.IncrementExperience(userInfo.FirstName+"@"+newId, utils.AvatarAddExt)
		if err != nil {
			return nil, fmt.Errorf("increment experience error: %w", err)
		}
	}

	err2 := s.repo.CreateUser(&newUser)
	if err2 != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	user, err := s.repo.GetByVkId(newUser.VkID)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("token generation error: %w", err)
	}
	return &dto.SignInResponseWithTmpLogin{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TmpLogin:     true,
	}, nil
}
func (s *UserService) getUserInfo(accessToken string) (*dto.VkAuthGetUserInfo, error) {
	baseURL := "https://api.vk.com/method/users.get"

	parameters := map[string]string{
		"access_token": accessToken,
		"v":            "5.131",
		"fields":       "email,photo_200",
	}
	response, err := utils.SendRequest[struct {
		Response []dto.VkAuthGetUserInfo `json:"response"`
	}](
		utils.RequestOptions{
			Method:      "GET",
			URL:         baseURL,
			Headers:     map[string]string{},
			Body:        map[string]string{},
			QueryParams: parameters,
			File:        []byte{},
			BasicAuth:   nil,
		},
	)

	if err != nil {
		return nil, err
	}

	return &response.Response[0], nil
}

func (s *UserService) TelegramOauth(req *dto.TelegramOauthRequest) (*dto.SignInResponseWithTmpLogin, error) {
	if req.LastName == "" {

	}
	data := []string{
		fmt.Sprintf("id=%d", req.ID),
		fmt.Sprintf("username=%s", req.Username),
		fmt.Sprintf("auth_date=%d", req.AuthDate),
	}

	if req.LastName != "" {
		data = append(data, fmt.Sprintf("last_name=%s", req.LastName))
	}
	if req.FirstName != "" {
		data = append(data, fmt.Sprintf("first_name=%s", req.FirstName))
	}
	if req.PhotoURL != "" {
		data = append(data, fmt.Sprintf("photo_url=%s", req.PhotoURL))
	}
	hash := req.Hash

	sort.Strings(data)
	dataCheckString := strings.Join(data, "\n")
	hSecretKey := sha256.New()
	hSecretKey.Write([]byte(os.Getenv("TELEGRAM_BOT_TOKEN")))

	h := hmac.New(sha256.New, hSecretKey.Sum(nil))

	h.Write([]byte(dataCheckString))

	hashBytes := h.Sum(nil)

	computedHash := hex.EncodeToString(hashBytes)
	if hash != computedHash {
		return nil, fmt.Errorf("проблемы с хэшем")
	}

	newId := fmt.Sprintf("%d", req.ID)

	existUser, _ := s.repo.GetByTgId(newId)

	if existUser != nil && !existUser.TmpLogin {
		if existUser.VerificationToken != "" {
			return nil, errors.New("user is not verified")
		}
		tokens, err := s.jwtService.GenerateTokenPair(existUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &dto.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     false,
		}, nil
	} else if existUser != nil && existUser.TmpLogin {
		tokens, err := s.jwtService.GenerateTokenPair(existUser)
		if err != nil {
			return nil, fmt.Errorf("token generation error: %w", err)
		}
		return &dto.SignInResponseWithTmpLogin{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			TmpLogin:     true,
		}, nil
	}

	passwordHash, err := auth.HashPassword("tg" + newId + req.Username)
	if err != nil {
		return nil, fmt.Errorf("password hash error: %w", err)
	}

	newUser := models.Users{
		TelegramID:   newId,
		Login:        req.Username + newId,
		UserName:     req.FirstName + " " + req.LastName,
		AvatarURL:    req.PhotoURL,
		Created:      time.Now().Format(time.RFC3339),
		PasswordHash: passwordHash,
		TmpLogin:     true,
		ProfileName:  req.Username,
	}

	if req.PhotoURL != "" {
		err = s.repo.IncrementExperience(req.Username+newId, utils.AvatarAddExt)
		if err != nil {
			return nil, fmt.Errorf("increment experience error: %w", err)
		}
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

	return &dto.SignInResponseWithTmpLogin{
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

func (s *UserService) ChangeResetPassword(request *dto.ChangePasswordReset) (*dto.SignInResponse, error) {
	existsUser, err := s.repo.GetByResetToken(request.Token)
	if err != nil {
		return nil, err
	}
	existsUser.PasswordHash, _ = auth.HashPassword(request.Password)
	existsUser.ResetToken = ""
	user, err := s.repo.UpdateUser(existsUser, *existsUser)
	if err != nil {
		return nil, err
	}
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("token generation error: %w", err)
	}
	return &dto.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) UpdateUser(user *dto.UserRequestUpdate, login string) (*dto.SignInResponse, error) {
	if login == "" {
		return nil, errors.New("not authorized")
	}

	existsUser, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}
	var userByEmail *models.Users
	if user.Email != nil {
		userByEmail, _ = s.repo.GetByEmail(*user.Email)
	}

	if userByEmail != nil && userByEmail.Email == *user.Email {
		return nil, errors.New("пользователь с таким email уже существует")

	}

	var exp int

	if existsUser.AvatarURL == "" && *user.AvatarURL != "" {
		exp += utils.AvatarAddExt
	} else if existsUser.AvatarURL != "" && *user.AvatarURL == "" {
		exp -= utils.AvatarAddExt
	}

	if existsUser.BackgroundURL == "" && *user.BackgroundURL != "" {
		exp += utils.BannerAddExp
	} else if existsUser.BackgroundURL != "" && *user.BackgroundURL == "" {
		exp -= utils.BannerAddExp
	}

	existsUser.Exp = existsUser.Exp + exp

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

	return &dto.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) UpdateOnlyOnce(data *dto.UserRequestUpdateFirstTime, targetUserLogin string) (
	*dto.CommonResponse,
	error,
) {
	userByLogin, _ := s.repo.GetUserByLogin(*data.Login)

	if userByLogin != nil && userByLogin.Login == *data.Login {
		return nil, errors.New("login already exists")
	}
	emailUser, err := s.repo.GetByEmail(*data.Email)
	if err != nil {
		_ = fmt.Errorf("failed to check user existence: %w", err)
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

	var exp int

	if existsUser.AvatarURL == "" && *data.AvatarURL != "" {
		exp += utils.AvatarAddExt
	} else if existsUser.AvatarURL != "" && *data.AvatarURL == "" {
		exp -= utils.AvatarAddExt
	}

	if existsUser.BackgroundURL == "" && *data.BackgroundURL != "" {
		exp += utils.BannerAddExp
	} else if existsUser.BackgroundURL != "" && *data.BackgroundURL == "" {
		exp -= utils.BannerAddExp
	}

	err = s.repo.IncrementExperience(targetUserLogin, exp)

	if err != nil {
		return nil, err
	}

	hexString, _ := utils.GenerateRandomString(32)

	existsUser.Login = *data.Login
	existsUser.UserName = *data.UserName
	existsUser.Email = *data.Email
	existsUser.AvatarURL = *data.AvatarURL
	existsUser.BackgroundURL = *data.BackgroundURL
	existsUser.City = *data.City
	existsUser.Bio = *data.Bio
	existsUser.TmpLogin = false
	existsUser.VerificationToken = hexString
	existsUser.ProfileName = *data.ProfileName

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

func (s *UserService) CheckDeletedUsers() {
	err := s.repo.DeleteDeletedUsersForever()
	if err != nil {
		fmt.Println(err)
	}
}

func (s *UserService) GetAll(search, limit, offset, order string) (*dto.DataUsersResponse, error) {
	users, total, donateMap, err := s.repo.GetAll(search, limit, offset, order)

	if err != nil {
		return nil, err
	}

	var result []dto.UserResponse

	for _, user := range users {
		var tempUser dto.UserResponse
		err = mapstructure.Decode(user, &tempUser)
		if err != nil {
			return nil, err
		}
		tempUser.TotalDonations = donateMap[user.Login]
		collectionItemsCount, _ := s.ciRepo.GetCountByUserLogin(user.Login)
		collectionsCount, _ := s.collectionRepo.GetCollectionsCount(user.Login)
		sum, _ := s.ciRepo.SumByUserLogin(user.Login)
		shippingTotal, _ := s.ciRepo.ShippingSumByUserLogin(user.Login)

		tempUser.TotalSum = int(sum)
		tempUser.CollectionItemsCount = int(collectionItemsCount)
		tempUser.CollectionsCount = int(collectionsCount)
		tempUser.ShippingTotal = int(shippingTotal)
		result = append(result, tempUser)
	}

	return &dto.DataUsersResponse{
		Data:  result,
		Total: total,
	}, nil
}

func (s *UserService) DeleteUser(login string) (*dto.CommonResponse, error) {

	err := s.repo.DeleteUser(login)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}
