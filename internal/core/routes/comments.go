package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func CommentsRouter(e *echo.Echo, jwtService *auth.JWTService, commentsService services.CommentsService) {
	controller := controllers.NewCommentsController(commentsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/comments", controller.GetAllComments)
		publicGroup.GET("/comments/answers", controller.LoadAnswers)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/comments", controller.CreateComment)
		securedGroup.DELETE("/comments/:id", controller.Delete)
		securedGroup.GET("/comments/like/:id", controller.Like)
		securedGroup.GET("/comments/dislike/:id", controller.Dislike)
	}
}
