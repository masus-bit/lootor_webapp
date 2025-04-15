package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/elasticsearch"
)

func SearchRouter(e *echo.Echo, jwtService *auth.JWTService, searchService elasticsearch.ElasticService) {
	controller := controllers.NewSearchController(&searchService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/search", controller.Search)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{

	}
}
