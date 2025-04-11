package user

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/auth"
)

func RegisterRoutes(e *echo.Echo, jwtService *auth.JWTService, userService Service) {
	controller := NewController(userService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.POST("/auth/signup", controller.SignUp)
		publicGroup.POST("/auth/signin", controller.SignIn)
		publicGroup.GET("/user", controller.GetByLogin)

	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("user/rating", controller.ChangeRating)
		securedGroup.POST("user/password", controller.ChangePass)
		securedGroup.POST("user/subscriptions", controller.Subscribe)
	}
}
