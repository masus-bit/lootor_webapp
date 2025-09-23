package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
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
// @Param getEvents body models.GetEventsRequest true "фильтры"
// @Success 201 {object} dto.EventsDataResponseSwagger
// @Router /secured/events [post]
func (c *EventsController) GetEvents(ctx echo.Context) error {
	var request models.GetEventsRequestForAll
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.eventsService.GetEvents(authUser, request.Limit, request.Offset, request.EventTargetTypes, request.Actions)
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
// @Param wishListItemId query string false "wish list item id"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.EventsDataResponseSwagger
// @Router /public/events/filter [get]
func (c *EventsController) GetFilteredEvents(ctx echo.Context) error {
	userLogin := ctx.QueryParam("userLogin")
	collectionId := ctx.QueryParam("collectionId")
	collectionItemId := ctx.QueryParam("collectionItemId")
	wishListItemId := ctx.QueryParam("wishListItemId")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.eventsService.GetFilteredEvents(userLogin, collectionId, collectionItemId, wishListItemId, limit, offset, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
