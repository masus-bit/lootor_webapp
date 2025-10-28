package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func TagsRouter(e *echo.Echo, jwtService *auth.JWTService, tagsService services.TagsService) {
	controller := controllers.NewTagsController(tagsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/tags/:id", controller.FindEntitiesByTag)
		publicGroup.GET("/tags/all", controller.GetAllTags)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.GET("/tags", controller.SearchTags)
		securedGroup.POST("/tags", controller.CreateTag)
		securedGroup.POST("/tags/adm/entity", controller.AddTagToEntity)
		securedGroup.POST("/tags/adm/entity/remove", controller.RemoveTagsFromEntity)
		securedGroup.POST("/tags/adm/merge", controller.MergeTags)
		securedGroup.POST("/tags/adm/merge_series", controller.MergeSeries)
		securedGroup.POST("/tags/adm/update/:id", controller.UpdateTag)
		securedGroup.GET("/tags/subscribe/:id", controller.Subscribe)
	}
}
