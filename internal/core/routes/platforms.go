package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func PlatformsRouter(e *echo.Echo, jwtService *auth.JWTService, platformsService services.PlatformsService) {
	controller := controllers.NewPlatformsController(platformsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/platforms", controller.GetPlatforms)
	}
}
