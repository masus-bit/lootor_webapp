package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func AchievementsRouter(e *echo.Echo, jwtService *auth.JWTService, achievementsService services.AchievementsService) {
	controller := controllers.NewAchievementsController(achievementsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/achievements/items", controller.GetAllItems)
		publicGroup.GET("/achievements", controller.GetAllUsersAchievements)
		publicGroup.GET("/achievements/one", controller.GetOneAchievement)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
	}
}
