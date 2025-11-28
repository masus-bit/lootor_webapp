package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/s3"
)

func S3Router(e *echo.Echo, jwtService *auth.JWTService, s3Service s3.StorageService) {
	controller := controllers.NewFilesController(&s3Service)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/images/:key", controller.GetFile)
		publicGroup.GET("/thumbs/:key", controller.GetThumbnail)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/images/upload", controller.Upload)
		securedGroup.POST("/images/delete", controller.Delete)

	}
}
