package payment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/core/services"
	"lootor/internal/infrastructure/achievementsclient"
	"lootor/internal/pkg/utils"
	"os"
)

type PayService struct {
	userRepo     *repositories.UsersRepository
	userService  *services.UserService
	subService   *services.SubscriptionService
	paymentsRepo *repositories.PaymentsRepository
	achClient    *achievementsclient.GRPCAchievementsClient
}

func NewPayService(
	userRepo *repositories.UsersRepository,
	userService *services.UserService,
	subService *services.SubscriptionService,
	paymentsRepo *repositories.PaymentsRepository,
	achClient *achievementsclient.GRPCAchievementsClient,
) *PayService {
	return &PayService{
		userRepo:     userRepo,
		userService:  userService,
		subService:   subService,
		paymentsRepo: paymentsRepo,
		achClient:    achClient,
	}
}

func (s *PayService) StartTransaction(login string, subType string, amount string) (*PayResponseToClientData, error) {
	existsUser, err := s.userRepo.GetUserByLogin(login)

	var description string

	if subType == "donate" {
		description = "Пожертвование"
	} else {
		description = "Подписка плана " + subType
	}

	if err != nil {
		return nil, err
	}

	item := ItemCustomer{
		Description: "subscription",
		Amount: Amount{
			Currency: "RUB",
			Value:    amount,
		},
		VatCode:        1,
		Quantity:       1,
		PaymentSubject: "service",
	}

	payment := Payment{
		Amount:  Amount{Value: amount, Currency: "RUB"},
		Capture: true,
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: "https://lootoe.space/payment-success",
		},
		Description: description,
		Receipt: Receipt{
			Customer: Customer{Email: existsUser.Email},
			Items:    []ItemCustomer{item},
		},
		Metadata: Metadata{
			UserLogin: login,
			OrderID:   uuid.New().String(),
			Type:      subType,
		},
	}

	idempotenceKey, _ := utils.GenerateRandomString(10)

	baseURL := os.Getenv("YOOKASSA_URL")

	headers := map[string]string{
		"Idempotence-Key": idempotenceKey,
		"Content-Type":    "application/json",
	}
	basicAuth := &struct {
		Username string
		Password string
	}{
		Username: os.Getenv("YOOKASSA_SHOP_ID"),
		Password: os.Getenv("YOOKASSA_SECRET_KEY"),
	}

	response, err := utils.SendRequest[PayResponse](
		utils.RequestOptions{
			Method:      "POST",
			URL:         baseURL,
			Headers:     headers,
			Body:        payment,
			QueryParams: nil,
			File:        nil,
			BasicAuth:   basicAuth,
		},
	)

	if err != nil {
		return nil, err
	}
	var result PayResponseToClient
	err = mapstructure.Decode(response, &result)
	if err != nil {
		return nil, err
	}

	return &PayResponseToClientData{Data: result}, nil
}

func (s *PayService) EndTransaction(data *Notification) {
	userLogin := data.Object.Metadata.UserLogin
	mode := os.Getenv("MODE")

	if data.Event != "payment.succeeded" {
		fmt.Println("expected notification: payment.succeeded")
	}

	months := 0
	years := 0

	subType := data.Object.Metadata.Type
	if subType == "donate" {
		payment := models.Payments{
			Amount:        data.Object.Amount.Value,
			InnerOrderID:  data.Object.Metadata.OrderID,
			TransactionID: data.Object.ID,
			UserLogin:     userLogin,
			PaymentType:   subType,
		}
		err := s.paymentsRepo.CreateRecord(&payment)
		if err != nil {
			fmt.Println(err)
		}

		return
	}
	if subType == "monthly" {
		months = 1
	} else {
		years = 1
	}

	if data.Object.Paid && data.Object.Status == "succeeded" {
		err := s.userService.ActivatePremium(userLogin, months, years, subType)
		if err != nil {
			return
		}
		ctx := context.Background()
		_, err = s.subService.CreateSubscription(ctx, userLogin, subType, months, years)
		if err != nil {
			return
		}
	}
	payment := models.Payments{
		Amount:        data.Object.Amount.Value,
		InnerOrderID:  data.Object.Metadata.OrderID,
		TransactionID: data.Object.ID,
		UserLogin:     userLogin,
		PaymentType:   subType,
	}
	err := s.paymentsRepo.CreateRecord(&payment)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		if mode == "beta" {
			err = utils.AddAchievement(
				s.achClient,
				utils.AchieveBetaTesterDonate,
				userLogin,
				1,
				utils.XPBetaTesterDonate,
				0,
			)
			err = s.userRepo.IncrementExperience(userLogin, utils.XPBetaTesterDonate)
		}
	}()

	if subType == "monthly" {
		err = s.userRepo.IncrementExperience(userLogin, utils.MonthlyDonateExp+utils.MonthlyDonateExp)
		if err != nil {
			return
		}
		go func() {
			err = utils.AddAchievement(s.achClient, utils.AchieveDonate, userLogin, 1, utils.XPDonateLevel1, 0)
		}()

	} else if subType == "yearly" {
		err = s.userRepo.IncrementExperience(userLogin, utils.YearlyDonateExp+utils.YearlyDonateExp)
		if err != nil {
			return
		}
		go func() {
			err = utils.AddAchievement(s.achClient, utils.AchieveDonate, userLogin, 2, utils.XPDonateLevel2, 0)
		}()
	}
}

func (s *PayService) FindAllPaymentsByLogin(login string) (*models.PaymentsData, error) {
	payments, err := s.paymentsRepo.FindAllPaymentsByLogin(login)

	if err != nil {
		return nil, err
	}
	return &models.PaymentsData{Data: payments}, nil
}

func (s *PayService) GetInfo() (*dto.PaymentsInfo, error) {
	yearTotalDonates, payments, err := s.paymentsRepo.FindYearPayments()
	if err != nil {
		return nil, err
	}
	var userLogins []string
	var paymentsResult []dto.PaymentDto
	for _, payment := range payments {
		userLogins = append(userLogins, payment.UserLogin)
	}
	subUsers, err := s.userRepo.GetForSubsMap(userLogins)
	if err != nil {
		return nil, err
	}

	for _, payment := range payments {
		paymentsResult = append(
			paymentsResult, dto.PaymentDto{
				CreatedAt:     payment.CreatedAt,
				ID:            payment.ID.String(),
				Amount:        payment.Amount,
				TransactionID: payment.TransactionID,
				InnerOrderID:  payment.InnerOrderID,
				User:          subUsers[payment.UserLogin],
				PaymentType:   payment.PaymentType,
			},
		)
	}
	data := struct {
		YearDonationsSum string           `json:"yearDonationsSum"`
		LastPayments     []dto.PaymentDto `json:"lastPayments"`
	}{
		YearDonationsSum: yearTotalDonates,
		LastPayments:     paymentsResult,
	}

	return &dto.PaymentsInfo{Data: data}, nil
}
