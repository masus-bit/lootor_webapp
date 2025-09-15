package routes

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/controllers"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
)

func PostsRouter(e *echo.Echo, jwtService *auth.JWTService, postsService services.PostsService) {
	controller := controllers.NewPostsController(postsService)

	publicGroup := e.Group("/public")
	publicGroup.Use(jwtService.AuthInfoMiddleware())
	{
		publicGroup.GET("/posts/:id", controller.GetPostById)
		publicGroup.GET("/posts/translit/:id", controller.GetPostByTranslit)
		publicGroup.GET("/posts/all", controller.GetAllPosts)
		publicGroup.GET("/posts", controller.GetPostsByUser)
		publicGroup.POST("/posts/views", controller.IncrementViews)
	}

	securedGroup := e.Group("/secured")
	securedGroup.Use(jwtService.RequireAuthMiddleware())
	{
		securedGroup.POST("/posts", controller.CreatePost)
		securedGroup.DELETE("/posts/:id", controller.Delete)
		securedGroup.GET("/posts/react/:id", controller.React)
		securedGroup.GET("/posts/react/delete/:id", controller.DeleteReact)
		securedGroup.PUT("/posts/:id", controller.UpdatePost)
	}
}
