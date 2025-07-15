package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func EntitiesRouter(e *echo.Echo, jwtService *auth.JWTService, entitiesService services.EntitiesService) {
	controller := controllers.NewEntitiesController(entitiesService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/entities/all", controller.GetAll)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/entities", controller.CreateEntities)
		securedGroup.GET("/entities", controller.SearchEntities)
	}
}
