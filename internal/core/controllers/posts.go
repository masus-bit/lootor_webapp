package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"lootor/internal/pkg/dto"
	"net/http"
)

type PostsController struct {
	postService services.PostsService
}

func NewPostsController(postService services.PostsService) *PostsController {
	return &PostsController{postService: postService}
}

// CreatePost
// @Summary Создание post
// @Tags posts
// @Accept  json
// @Produce  json
// @Router /secured/posts [post]
func (c *PostsController) CreatePost(ctx echo.Context) error {
	var request dto.PostRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	request.Author = authUser

	context := ctx.Request().Context()
	response, err := c.postService.CreatePost(context, &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetAllPosts
// @Summary все посты с пагинацией
// @Tags posts
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Param order query string true "order"
// @Success 201 {object} dto.FeedDataResponseSwag
// @Router /public/posts/all [get]
func (c *PostsController) GetAllPosts(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	order := ctx.QueryParam("order")
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	response, err := c.postService.GetAllPosts(ctx.Request().Context(), order, limit, offset, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetPostsByUser
// @Summary получить посты юзера
// @Tags posts
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Param login query string true "login"
// @Param isDraft query string true "is draft"
// @Success 201 {object} dto.FeedDataResponseSwag
// @Router /public/posts [get]
func (c *PostsController) GetPostsByUser(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	login := ctx.QueryParam("login")
	isDraft := ctx.QueryParam("isDraft")
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	if isDraft == "true" && authInfo.UserLogin != login {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized. You can only get your own draft posts.",
		})

	}
	response, err := c.postService.GetPostsByUser(ctx.Request().Context(), login, limit, offset, authInfo.UserLogin, isDraft == "true")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetPostById
// @Summary получить пост по id
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /public/posts/{id} [get]
func (c *PostsController) GetPostById(ctx echo.Context) error {
	id := ctx.Param("id")
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	response, err := c.postService.GetPostById(ctx.Request().Context(), id, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// React
// @Summary поставить реакцию
// @Tags posts
// @Accept  json
// @Produce  json
// @Param reaction query string true "reaction"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/posts/react/{id} [get]
func (c *PostsController) React(ctx echo.Context) error {
	reaction := ctx.QueryParam("reaction")
	id := ctx.Param("id")
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	request := dto.ReactRequest{
		PostId:    id,
		Reaction:  reaction,
		UserLogin: authUser,
	}
	response, err := c.postService.React(ctx.Request().Context(), &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// DeleteReact
// @Summary убрать реакцию
// @Tags posts
// @Accept  json
// @Produce  json
// @Param reaction query string true "reaction"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/posts/react/delete/{id} [get]
func (c *PostsController) DeleteReact(ctx echo.Context) error {
	id := ctx.Param("id")
	reaction := ctx.QueryParam("reaction")
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	request := dto.ReactRequest{
		PostId:    id,
		UserLogin: authUser,
		Reaction:  reaction,
	}
	response, err := c.postService.Unreact(ctx.Request().Context(), &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Delete
// @Summary удаление posta
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/posts/{id} [delete]
func (c *PostsController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.postService.DeletePost(ctx.Request().Context(), id, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// UpdatePost
// @Summary обновить контент поста
// @Tags posts
// @Accept  json
// @Produce  json
// @Param updateRequest body dto.PostUpdateRequest true "поля"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/posts/{id} [put]
func (c *PostsController) UpdatePost(ctx echo.Context) error {
	var request dto.PostUpdateRequest
	id := ctx.Param("id")
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	request.Id = id
	response, err := c.postService.UpdatePost(ctx.Request().Context(), &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// IncrementViews
// @Summary просмотры
// @Tags posts
// @Accept  json
// @Produce  json
// @Param updateRequest body dto.IncrementRequest true "поля"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/posts/views [post]
func (c *PostsController) IncrementViews(ctx echo.Context) error {
	var request dto.IncrementRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	response, err := c.postService.IncrementViews(ctx.Request().Context(), &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
