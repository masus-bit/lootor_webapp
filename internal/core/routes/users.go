package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func RegisterRoutes(e *echo.Echo, jwtService *auth.JWTService, userService services.UserService) {
	controller := controllers.NewUserController(userService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.POST("/auth/signup", controller.SignUp)
		publicGroup.POST("/auth/signin", controller.SignIn)
		publicGroup.GET("/user", controller.GetByLogin)
		publicGroup.GET("/auth/verification", controller.Verification)
		publicGroup.POST("/auth/refresh", controller.RenewTokens)
		publicGroup.GET("/auth/vk/callback", controller.VkOauth)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/user/rating", controller.ChangeRating)
		securedGroup.POST("/user/password", controller.ChangePass)
		securedGroup.POST("/user/subscriptions", controller.Subscribe)
	}
}
