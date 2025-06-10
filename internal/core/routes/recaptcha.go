package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
)

func RecaptchaRouter(e *echo.Echo, jwtService *auth.JWTService, service auth.RecaptchaService) {
	controller := controllers.NewRecaptchaController(&service)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/recaptcha_verify", controller.CheckRecaptcha)
	}
}
