package payment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/core/services"
	"lootor/internal/pkg/utils"
	"os"
	"time"
)

type PayService struct {
	userRepo     *repositories.UsersRepository
	userService  *services.UserService
	subService   *services.SubscriptionService
	paymentsRepo *repositories.PaymentsRepository
}

func NewPayService(userRepo *repositories.UsersRepository, userService *services.UserService, subService *services.SubscriptionService, paymentsRepo *repositories.PaymentsRepository) *PayService {
	return &PayService{userRepo: userRepo, userService: userService, subService: subService, paymentsRepo: paymentsRepo}
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
			ReturnUrl: "https://lootor.me/payment-success",
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

	baseUrl := os.Getenv("YOOKASSA_URL")

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

	response, err := utils.SendRequest[PayResponse](utils.RequestOptions{Method: "POST", URL: baseUrl, Headers: headers, Body: payment, QueryParams: nil, File: nil, BasicAuth: basicAuth})

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

	if data.Event != "payment.succeeded" {
		fmt.Println("expected notification: payment.succeeded")
	}

	//months := 0
	//years := 0

	subType := data.Object.Metadata.Type
	if subType == "donate" {
		return
	} else {
		//if subType == "monthly" {
		//	months = 1
		//} else {
		//	years = 1
		//}
	}

	if data.Object.Paid && data.Object.Status == "succeeded" {
		err := s.userService.ActivateTestPremium(userLogin, 3*time.Minute, subType)
		if err != nil {
			return
		}
		ctx := context.Background()
		_, err = s.subService.CreateTestSubscription(ctx, userLogin, subType, 3*time.Minute)
		if err != nil {
			return
		}
	}
	payment := models.Payments{
		Amount: data.Object.Amount.Value,
	}
	err := s.paymentsRepo.CreateRecord(&payment)
	if err != nil {
		fmt.Println(err)
	}
}
