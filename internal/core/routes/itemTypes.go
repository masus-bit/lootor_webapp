package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func ItemTypesRouter(e *echo.Echo, jwtService *auth.JWTService, itemTypesService services.ItemTypesService) {
	controller := controllers.NewItemTypesController(itemTypesService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/item_types", controller.GetItemTypes)

	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
	}
}
