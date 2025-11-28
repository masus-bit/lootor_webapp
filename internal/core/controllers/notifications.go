package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"lootor/internal/pkg/dto"
	"net/http"
)

type NotificationsController struct {
	notificationsService services.NotificationsService
}

func NewNotificationsController(notificationsService services.NotificationsService) *NotificationsController {
	return &NotificationsController{notificationsService: notificationsService}
}

// GetAllByLogin
// @Summary получить уведомления для юзера
// @Tags notifications
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.FeedResponseSwag
// @Router /secured/notifications [get]
func (c *NotificationsController) GetAllByLogin(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	context := ctx.Request().Context()
	response, err := c.notificationsService.GetAllNotifications(context, authUser, limit, offset, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// ReadNotifications
// @Summary прочитать уведомления
// @Tags notifications
// @Accept  json
// @Produce  json
// @Param createRequest body dto.FeedRequestSwag true "поля создания"
// @Success 201 {object} dto.NotificationsReadRequest
// @Router /secured/notifications/read [post]
func (c *NotificationsController) ReadNotifications(ctx echo.Context) error {
	var request dto.NotificationsReadRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	response, err := c.notificationsService.ReadNotification(ctx.Request().Context(), request.IDs)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
