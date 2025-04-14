package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
	"lootor/internal/core/services"
	"net/http"
)

type CollectionsController struct {
	colService services.CollectionService
}

func NewCollectionsController(colService services.CollectionService) *CollectionsController {
	return &CollectionsController{colService: colService}
}

func (c *CollectionsController) CreateCollection(ctx echo.Context) error {
	var request models.CollectionCreateRequest

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	result, err := c.colService.Create(&request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (c *CollectionsController) DeleteCollection(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}

	value, err := c.colService.Delete(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, value)
}

func (c *CollectionsController) UpdateCollection(ctx echo.Context) error {
	var request models.CollectionCreateRequest
	id := ctx.QueryParam("id")

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	value, err := c.colService.Update(id, &request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, value)
}

func (c *CollectionsController) GetByUserLogin(ctx echo.Context) error {
	login := ctx.QueryParam("userLogin")
	if login == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Login parameter is required",
		})
	}
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	response, err := c.colService.GetByUserLogin(login, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CollectionsController) GetOneByFewParams(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	transliteration := ctx.QueryParam("transliteration")
	shareString := ctx.QueryParam("shareString")
	userLogin := ctx.QueryParam("userLogin")

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.colService.GetOne(authUser, id, transliteration, userLogin, shareString)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CollectionsController) Like(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.colService.Like(id, authUser)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CollectionsController) Subscribe(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	isSubscribe := ctx.QueryParam("isSubscribe")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.colService.Subscribe(id, authUser, isSubscribe == "true")
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CollectionsController) GetByTag(ctx echo.Context) error {
	tag := ctx.QueryParam("tag")
	if tag == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tag parameter is required",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.colService.GetByTag(tag, authUser)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
