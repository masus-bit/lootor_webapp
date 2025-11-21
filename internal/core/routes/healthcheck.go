package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
)

func HealthCheckRouter(e *echo.Echo) {
	controller := controllers.NewHealthCheckController()

	publicGroup := e.Group("/health")
	{
		publicGroup.GET("", controller.Healthcheck)
	}
}
