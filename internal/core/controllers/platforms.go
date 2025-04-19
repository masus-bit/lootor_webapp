package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"net/http"
)

type PlatformsController struct {
	platformsService services.PlatformsService
}

func NewPlatformsController(platformsService services.PlatformsService) *PlatformsController {
	return &PlatformsController{platformsService: platformsService}
}

// GetPlatforms
// @Summary получить список платформ
// @Tags platforms
// @Accept  json
// @Produce  json
// @Success 201 {object} models.PlatformsDataResponse
// @Router /secured/platforms [get]
func (c *PlatformsController) GetPlatforms(ctx echo.Context) error {
	response, err := c.platformsService.GetAllPlatforms()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
