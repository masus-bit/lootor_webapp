package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func RegisterCollectionsRoutes(e *echo.Echo, jwtService *auth.JWTService, colService services.CollectionService) {
	controller := controllers.NewCollectionsController(colService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/collections/all", controller.GetByUserLogin)
		publicGroup.GET("/collections", controller.GetOneByFewParams)

	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/collections", controller.CreateCollection)
		securedGroup.DELETE("/collections", controller.DeleteCollection)
		securedGroup.POST("/collections", controller.UpdateCollection)
		securedGroup.GET("/collections/like", controller.Like)
		securedGroup.GET("/collections/tag", controller.GetByTag)
		securedGroup.GET("/collections/subscribe", controller.Subscribe)

	}
}
