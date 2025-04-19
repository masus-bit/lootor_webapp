package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"net/http"
)

type TagsController struct {
	tagsService services.TagsService
}

func NewTagsController(tagsService services.TagsService) *TagsController {
	return &TagsController{tagsService: tagsService}
}

// SearchTags
// @Summary Поиск по тегам
// @Tags tags
// @Accept  json
// @Produce  json
// @Param name query string true "name"
// @Success 201 {object} dto.TagsSwagger
// @Router /public/tags [get]
func (c *TagsController) SearchTags(ctx echo.Context) error {
	name := ctx.QueryParam("name")
	response, err := c.tagsService.SearchTags(name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
