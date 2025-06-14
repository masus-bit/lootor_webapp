package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/payment"
)

func PaymentRouter(e *echo.Echo, jwtService *auth.JWTService, paymentService payment.PayService) {
	controller := controllers.NewPaymentController(paymentService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.POST("/notification_payment", controller.ConfirmPayment)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/payment", controller.CreatePayment)
	}
}
