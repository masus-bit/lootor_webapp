package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func EventsRouter(e *echo.Echo, jwtService *auth.JWTService, eventsService services.EventsService) {
	controller := controllers.NewEventsController(eventsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/events/filter", controller.GetFilteredEvents)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/events", controller.GetEvents)
	}
}
