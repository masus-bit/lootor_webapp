package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func FeedRouter(e *echo.Echo, jwtService *auth.JWTService, feedService services.FeedService) {
	controller := controllers.NewFeedController(feedService)

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/feed/:id", controller.GetOne)
		securedGroup.GET("/feed", controller.GetAll)
		securedGroup.POST("/feed", controller.CreateFeedItem)
		securedGroup.DELETE("/feed/:id", controller.Delete)
		securedGroup.PUT("/feed/:id", controller.Update)
	}
}
