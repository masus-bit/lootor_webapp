package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"net/http"
	"slices"
)

type FeedController struct {
	feedService services.FeedService
}

func NewFeedController(feedService services.FeedService) *FeedController {
	return &FeedController{feedService: feedService}
}

// CreateFeedItem
// @Summary Создание новости
// @Tags feed
// @Accept  json
// @Produce  json
// @Param createRequest body dto.FeedRequestSwag true "поля создания"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/feed [post]
func (c *FeedController) CreateFeedItem(ctx echo.Context) error {
	var request dto.FeedRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	adminList := utils.GetAdminLogins()

	if !slices.Contains(adminList, authUser) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	context := ctx.Request().Context()
	response, err := c.feedService.CreateFeed(context, &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetAll
// @Summary все новости
// @Tags feed
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.FeedDataResponseSwag
// @Router /public/feed [get]
func (c *FeedController) GetAll(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	response, err := c.feedService.GetAllFeed(ctx.Request().Context(), limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetOne
// @Summary одна новость
// @Tags feed
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /public/feed/{id} [get]
func (c *FeedController) GetOne(ctx echo.Context) error {
	id := ctx.Param("id")
	response, err := c.feedService.GetFeed(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Update
// @Summary обновление новости
// @Tags feed
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param updateRequest body dto.FeedRequestSwag true "поля обновления"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/feed/{id} [put]
func (c *FeedController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var request dto.FeedRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	adminList := utils.GetAdminLogins()

	if !slices.Contains(adminList, authUser) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	response, err := c.feedService.UpdateFeed(ctx.Request().Context(), id, &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Delete
// @Summary удаление новости
// @Tags feed
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/feed/{id} [delete]
func (c *FeedController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	adminList := utils.GetAdminLogins()

	if !slices.Contains(adminList, authUser) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	response, err := c.feedService.DeleteFeed(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
