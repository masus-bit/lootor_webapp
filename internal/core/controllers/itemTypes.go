package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"net/http"
)

type ItemTypesController struct {
	itemTypesService services.ItemTypesService
}

func NewItemTypesController(itemTypesService services.ItemTypesService) *ItemTypesController {
	return &ItemTypesController{itemTypesService: itemTypesService}
}

// GetItemTypes
// @Summary получить список типов КИ
// @Tags itemTypes
// @Accept  json
// @Produce  json
// @Success 201 {object} dto.ItemTypesDataResponse
// @Router /secured/item_types [get]
func (c *ItemTypesController) GetItemTypes(ctx echo.Context) error {
	response, err := c.itemTypesService.GetAllTypes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
