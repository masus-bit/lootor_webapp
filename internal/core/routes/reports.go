package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/feedback"
)

func ReportsRouter(e *echo.Echo, jwtService *auth.JWTService, service feedback.ReportsService) {
	controller := controllers.NewReportsController(&service)

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/reports", controller.Report)
	}
}
