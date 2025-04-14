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
