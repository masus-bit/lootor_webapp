package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/repositories"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/elasticsearch"
)

func ReindexRouter(e *echo.Echo, jwtService *auth.JWTService, esService elasticsearch.ElasticService, userRepo repositories.UsersRepository, ciRepo repositories.CiRepository, collectionRepo repositories.CollectionsRepository, tagRepo services.TagsService) {
	controller := controllers.NewReindexController(&esService, &userRepo, &collectionRepo, &ciRepo, &tagRepo)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/recreate-elastic-search-index", controller.Reindex)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{

	}
}
