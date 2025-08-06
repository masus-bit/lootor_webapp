package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"lootor/internal/pkg/dto"
	"net/http"
)

type CommentsController struct {
	commentService services.CommentsService
}

func NewCommentsController(commentService services.CommentsService) *CommentsController {
	return &CommentsController{commentService: commentService}
}

// CreateComment
// @Summary Создание коммента
// @Tags comments
// @Accept  json
// @Produce  json
// @Param createRequest body dto.CommentsRequestSwag true "поля создания"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/comments [post]
func (c *CommentsController) CreateComment(ctx echo.Context) error {
	var request dto.CommentsRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	context := ctx.Request().Context()
	response, err := c.commentService.CreateComment(context, &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetAllComments
// @Summary все комменты к сущности
// @Tags comments
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Param targetId query string true "target id"
// @Success 201 {object} dto.FeedDataResponseSwag
// @Router /public/comments [get]
func (c *CommentsController) GetAllComments(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	targetId := ctx.QueryParam("targetId")
	response, err := c.commentService.GetAllComments(ctx.Request().Context(), targetId, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Like
// @Summary лайк коммента
// @Tags comments
// @Accept  json
// @Produce  json
// @Param targetId path string true "target id"
// @Param login query string true "login"
// @Param isLike query bool true "isLike"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/comments/like/{id} [get]
func (c *CommentsController) Like(ctx echo.Context) error {
	id := ctx.Param("id")
	login := ctx.QueryParam("login")
	isLike := ctx.QueryParam("isLike")
	response, err := c.commentService.LikeComment(ctx.Request().Context(), id, login, isLike == "true")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Delete
// @Summary удаление коммента
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/comments/{id} [delete]
func (c *CommentsController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	response, err := c.commentService.DeleteComment(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
