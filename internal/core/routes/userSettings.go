package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func UserSettingsRouter(e *echo.Echo, jwtService *auth.JWTService, userSettingsService services.UserSettingsService) {
	controller := controllers.NewUserSettingsController(userSettingsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/settings", controller.CreateOrUpdateSettings)
	}
}
