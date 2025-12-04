package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"net/http"
)

type AchievementsController struct {
	achievementsService services.AchievementsService
}

func NewAchievementsController(achievementsService services.AchievementsService) *AchievementsController {
	return &AchievementsController{achievementsService: achievementsService}
}

// GetAllItems
// @Summary Получить объекты достижений
// @Tags achievements
// @Accept  json
// @Produce  json
// @Success 201 {object} models.AchievementsItemsResponse
// @Router /public/achievements/items [get]
func (c *AchievementsController) GetAllItems(ctx echo.Context) error {
	response, err := c.achievementsService.GetAllAchievementsItems()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetAllUsersAchievements
// @Summary Получить все достижения юзера
// @Tags achievements
// @Accept  json
// @Produce  json
// @Param login query string true "login"
// @Success 201 {object} models.AchievementsResponse
// @Router /public/achievements [get]
func (c *AchievementsController) GetAllUsersAchievements(ctx echo.Context) error {
	userLogin := ctx.QueryParam("login")
	response, err := c.achievementsService.GetUserAchievements(userLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetOneAchievement
// @Summary Получить одну ачиву юзера
// @Tags achievements
// @Accept  json
// @Produce  json
// @Param login query string true "login"
// @Param code query string true "code"
// @Success 201 {object} models.AchievementResponse
// @Router /public/achievements/one [get]
func (c *AchievementsController) GetOneAchievement(ctx echo.Context) error {
	userLogin := ctx.QueryParam("login")
	code := ctx.QueryParam("code")
	response, err := c.achievementsService.GetOneAchievement(userLogin, code)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
