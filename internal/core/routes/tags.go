package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func TagsRouter(e *echo.Echo, jwtService *auth.JWTService, tagsService services.TagsService) {
	controller := controllers.NewTagsController(tagsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/tags", controller.SearchTags)
	}
}
