package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/feedback"
)

func FeedbackRouter(e *echo.Echo, jwtService *auth.JWTService, service feedback.FeedbackService) {
	controller := controllers.NewFeedbackController(&service)

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/feedback", controller.AddIssue)
	}
}
