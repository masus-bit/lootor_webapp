package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/feedback"
)

func FeedbackRouter(e *echo.Echo, jwtService *auth.JWTService, service feedback.FBService) {
	controller := controllers.NewFeedbackController(&service)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.POST("/feedback", controller.AddIssue)

	}
}
