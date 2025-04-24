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

// CreateCollection
// @Summary Создание коллекции
// @Tags collections
// @Accept  json
// @Produce  json
// @Param createRequest body models.CollectionCreateRequest true "поля создания"
// @Success 201 {object} dto.CollectionDataResponseSwagger
// @Router /secured/collections [post]
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

	value, err := c.colService.Delete(id, ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, value)
}

// UpdateCollection
// @Summary обновление коллекции
// @Tags collections
// @Accept  json
// @Produce  json
// @Param createRequest body models.CollectionCreateRequest true "поля создания"
// @Success 201 {object} dto.CollectionDataResponseSwagger
// @Router /secured/collections/update [post]
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

// GetByUserLogin
// @Summary получение коллекций по логину
// @Tags collections
// @Accept  json
// @Produce  json
// @Param userLogin query string true "логин"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.AllCollectionsDataResponseSwagger
// @Router /public/collections/all [get]
func (c *CollectionsController) GetByUserLogin(ctx echo.Context) error {
	login := ctx.QueryParam("userLogin")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")

	if login == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Login parameter is required",
		})
	}
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	response, err := c.colService.GetByUserLogin(login, authInfo.UserLogin, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetOneByFewParams
// @Summary получение коллекций по id или по логину + транслит
// @Tags collections
// @Accept  json
// @Produce  json
// @Param userLogin query string false "логин"
// @Param id query string false "id"
// @Param transliteration query string false "translit"
// @Param ciLimit query string true "limit"
// @Param ciOffset query string true "offset"
// @Success 201 {object} dto.CollectionDataResponseSwagger
// @Router /public/collections [get]
func (c *CollectionsController) GetOneByFewParams(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	transliteration := ctx.QueryParam("transliteration")
	shareString := ctx.QueryParam("shareString")
	userLogin := ctx.QueryParam("userLogin")
	limit := ctx.QueryParam("ciLimit")
	offset := ctx.QueryParam("ciOffset")

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.colService.GetOne(authInfo.UserLogin, id, transliteration, userLogin, shareString, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// Like
// @Summary like/dislike
// @Tags collections
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/collections/like [get]
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

// Subscribe
// @Summary подписка на кололекцию
// @Tags collections
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Param isSubscribe query boolean true "flag"
// @Success 201 {object} dto.CommonResponse
// @Router /public/collections/subscribe [get]
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

// GetByTag
// @Summary получить коллекции по тегу
// @Tags collections
// @Accept  json
// @Produce  json
// @Param tag query string true "tag"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.AllCollectionsDataResponseSwagger
// @Router /secured/collections/tag [get]
func (c *CollectionsController) GetByTag(ctx echo.Context) error {
	tag := ctx.QueryParam("tag")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	if tag == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tag parameter is required",
		})
	}

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.colService.GetByTag(tag, authInfo.UserLogin, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
