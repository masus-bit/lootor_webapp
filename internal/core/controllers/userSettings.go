package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/dto"
	"lootor/internal/core/services"
	"net/http"
)

type UserSettingsController struct {
	userSettingsService services.UserSettingsService
}

func NewUserSettingsController(userSettingsService services.UserSettingsService) *UserSettingsController {
	return &UserSettingsController{userSettingsService: userSettingsService}
}

// CreateOrUpdateSettings
// @Summary Создание/обновление настроек юзера
// @Tags settings
// @Accept  json
// @Produce  json
// @Param updateRequest body dto.UserSettingsRequest true "поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/settings [post]
func (c *UserSettingsController) CreateOrUpdateSettings(ctx echo.Context) error {
	var request dto.UserSettingsRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.userSettingsService.UpdateSettings(authUser, &request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}
