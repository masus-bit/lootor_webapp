package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
	"lootor/internal/core/services"
	"net/http"
)

type EntitiesController struct {
	entitiesService services.EntitiesService
}

func NewEntitiesController(entitiesService services.EntitiesService) *EntitiesController {
	return &EntitiesController{entitiesService: entitiesService}
}

func (c *EntitiesController) CreateEntities(ctx echo.Context) error {
	var request models.EntitiesCreateRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	response, err := c.entitiesService.CreateEntities(&request)
	fmt.Println(response)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *EntitiesController) SearchEntities(ctx echo.Context) error {
	query := ctx.QueryParam("name")
	response, err := c.entitiesService.SearchEntities(query)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
