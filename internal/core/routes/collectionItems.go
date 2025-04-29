package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func RegisterCollectionItemsRoutes(e *echo.Echo, jwtService *auth.JWTService, ciService services.CiService) {
	controller := controllers.NewCIController(ciService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/collection_item", controller.GetCollectionItem)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/collection_item", controller.CreateCollectionItem)
		securedGroup.POST("/collection_item/update", controller.UpdateCollectionItem)
		securedGroup.DELETE("/collection_item", controller.DeleteCollectionItem)
		securedGroup.POST("/collection_item/copy", controller.CopyOrMoveCollectionItem)
		securedGroup.GET("/collection_item/like", controller.Like)
		securedGroup.GET("/collection_items/entity", controller.GetByEntity)
	}
}
