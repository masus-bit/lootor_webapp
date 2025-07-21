package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func FeedRouter(e *echo.Echo, jwtService *auth.JWTService, feedService services.FeedService) {
	controller := controllers.NewFeedController(feedService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/feed/:id", controller.GetOne)
		publicGroup.GET("/feed", controller.GetAll)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/feed", controller.CreateFeedItem)
		securedGroup.DELETE("/feed/:id", controller.Delete)
		securedGroup.PUT("/feed/:id", controller.Update)
	}
}
