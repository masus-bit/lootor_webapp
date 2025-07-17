package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/payment"
	"net/http"
)

type PaymentController struct {
	payService payment.PayService
}

func NewPaymentController(payService payment.PayService) *PaymentController {
	return &PaymentController{payService: payService}
}

// CreatePayment
// @Summary Создать Платеж
// @Description создание платежа
// @Tags payment
// @Param subType query string true "тип подписки monthly/yearly"
// @Param amount query string true "сумма"
// @Success 200 {object} payment.PayResponseToClientData
// @Router /secured/payment [get]
func (c *PaymentController) CreatePayment(ctx echo.Context) error {
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	subType := ctx.QueryParam("subType")
	amount := ctx.QueryParam("amount")

	if subType == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "subType обязателен"})
	}

	response, err := c.payService.StartTransaction(authUser, subType, amount)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response)
}

func (c *PaymentController) ConfirmPayment(ctx echo.Context) error {
	var request payment.Notification

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	c.payService.EndTransaction(&request)
	return nil
}

// GetUserPayments
// @Summary Получить список платежей юзера
// @Description список платежей
// @Tags payment
// @Success 200 {object} models.PaymentsData
// @Router /secured/payment/info [get]
func (c *PaymentController) GetUserPayments(ctx echo.Context) error {
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.payService.FindAllPaymentsByLogin(authUser)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response)
}
