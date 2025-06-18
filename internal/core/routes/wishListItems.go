package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func WLRouter(e *echo.Echo, jwtService *auth.JWTService, wlService services.WLService) {
	controller := controllers.NewWLController(wlService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/wishlist", controller.GetWishList)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/wishlist", controller.AddWishListItem)
		securedGroup.POST("/wishlist/priority", controller.UpdatePriority)
		securedGroup.DELETE("/wishlist", controller.Delete)
		securedGroup.POST("/wishlist/update", controller.Update)
	}
}
