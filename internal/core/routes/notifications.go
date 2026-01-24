package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func NotificationsRouter(
	e *echo.Echo,
	jwtService *auth.JWTService,
	notificationsService services.NotificationsService,
) {
	controller := controllers.NewNotificationsController(notificationsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/notifications/read", controller.ReadNotifications)
		securedGroup.GET("/notifications", controller.GetAllByLogin)
	}
}
