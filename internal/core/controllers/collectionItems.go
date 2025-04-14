package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
	"lootor/internal/core/services"
	"net/http"
)

type CIController struct {
	ciService services.CiService
}

func NewCIController(ciService services.CiService) *CIController {
	return &CIController{ciService: ciService}
}

func (c *CIController) CreateCollectionItem(ctx echo.Context) error {
	var request models.CollectionItemsRequestCreate

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	value, err := c.ciService.Create(&request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, value)
}

func (c *CIController) DeleteCollectionItem(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}

	value, err := c.ciService.Delete(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, value)
}

func (c *CIController) UpdateCollectionItem(ctx echo.Context) error {
	var request models.CollectionItemsRequestCreate
	id := ctx.QueryParam("id")

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	value, err := c.ciService.Update(id, &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, value)
}

func (c *CIController) GetCollectionItem(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}

	response, err := c.ciService.GetById(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CIController) CopyOrMoveCollectionItem(ctx echo.Context) error {
	var request models.CollectionItemsCopyOrMoveRequest

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	value, err := c.ciService.CopyOrMove(request.Id, request.TargetCollectionIds, request.SourceCollectionId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, value)
}
