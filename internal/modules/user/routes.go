package user

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/auth"
)

func RegisterRoutes(e *echo.Echo, jwtService *auth.JWTService, userService Service) {
	controller := NewController(userService)

	publicGroup := e.Group("/public")
	{
		publicGroup.POST("/auth/signup", controller.SignUp)
		publicGroup.POST("/auth/signin", controller.SignIn)
		publicGroup.GET("/user/:login", controller.GetByLogin)

	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.EchoMiddleware())
	{
	}
}
