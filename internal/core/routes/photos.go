package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func PhotosRouter(e *echo.Echo, jwtService *auth.JWTService, photosService services.PhotosService) {
	controller := controllers.NewPhotosController(photosService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/photos/:id", controller.GetById)
		publicGroup.GET("/photos/by_user", controller.GetByUserLogin)
		publicGroup.GET("/photos/by_collection", controller.FindAllByCollectionId)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/photos", controller.CreatePhotos)
		securedGroup.DELETE("/photos/:id", controller.Delete)
		securedGroup.GET("/photos/like/:id", controller.Like)
		securedGroup.PATCH("/photos", controller.Update)
	}
}
