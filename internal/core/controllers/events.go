package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/services"
	"net/http"
)

type EventsController struct {
	eventsService services.EventsService
}

func NewEventsController(eventsService services.EventsService) *EventsController {
	return &EventsController{eventsService: eventsService}
}

// GetEvents
// @Summary получить эвенты по юзеру из токена
// @Tags events
// @Accept  json
// @Produce  json
// @Success 201 {object} dto.EventsDataResponseSwagger
// @Router /secured/events [get]
func (c *EventsController) GetEvents(ctx echo.Context) error {
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.eventsService.GetEvents(authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetFilteredEvents
// @Summary получить эвенты filtered
// @Tags events
// @Accept  json
// @Produce  json
// @Param userLogin query string false "login"
// @Param collectionId query string false "collection id"
// @Param collectionItemId query string false "collection item id"
// @Success 201 {object} dto.EventsDataResponseSwagger
// @Router /secured/events/filter [get]
func (c *EventsController) GetFilteredEvents(ctx echo.Context) error {
	userLogin := ctx.QueryParam("userLogin")
	collectionId := ctx.QueryParam("collectionId")
	collectionItemId := ctx.QueryParam("collectionItemId")

	response, err := c.eventsService.GetFilteredEvents(userLogin, collectionId, collectionItemId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
