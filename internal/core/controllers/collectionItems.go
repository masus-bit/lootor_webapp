package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
	"lootor/internal/core/services"
	"lootor/internal/pkg/utils"
	"net/http"
)

type CIController struct {
	ciService       services.CiService
	enrichedService utils.EnrichingCIService
}

func NewCIController(ciService services.CiService, enrichedService utils.EnrichingCIService) *CIController {
	return &CIController{ciService: ciService, enrichedService: enrichedService}
}

// CreateCollectionItem
// @Summary Создание КИ
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param createRequest body models.CollectionItemsRequestCreate true "поля создания"
// @Success 201 {object} dto.CollectionItemsDataResponseSwagger
// @Router /secured/collection_item [post]
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

// DeleteCollectionItem
// @Summary удаление КИ
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/collection_item [delete]
func (c *CIController) DeleteCollectionItem(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}

	value, err := c.ciService.Delete(id, ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, value)
}

// UpdateCollectionItem
// @Summary обновление КИ
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param createRequest body models.CollectionItemsRequestCreate true "поля создания"
// @Success 201 {object} dto.CollectionItemsDataResponseSwagger
// @Router /secured/collection_item/update [post]
func (c *CIController) UpdateCollectionItem(ctx echo.Context) error {
	var request models.CollectionItemsRequestCreate
	id := ctx.QueryParam("id")

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	value, err := c.enrichedService.Update(id, &request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, value)
}

// GetCollectionItem
// @Summary получение одного КИ
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param id query string false "id"
// @Success 201 {object} dto.CollectionItemsDataResponseSwagger
// @Router /public/collection_item [get]
func (c *CIController) GetCollectionItem(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Id parameter is required",
		})
	}
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.enrichedService.GetById(id, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// CopyOrMoveCollectionItem
// @Summary копирование/перемещение КИ
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param createRequest body models.CollectionItemsCopyOrMoveRequest true "polya"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/collection_item/copy [post]
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

// Like
// @Summary like/dislike
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/collection_item/like [get]
func (c *CIController) Like(ctx echo.Context) error {
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
	response, err := c.ciService.Like(id, authUser)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetByEntity
// @Summary получить КИ по кнтити
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param entity query string true "entity"
// @Param limit query string true "limit"
// @Success 201 {object} dto.AllCollectionsDataResponseSwagger
// @Router /public/collection_item/entity [get]
func (c *CIController) GetByEntity(ctx echo.Context) error {
	entity := ctx.QueryParam("entity")
	limit := ctx.QueryParam("limit")

	if entity == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Entity parameter is required",
		})
	}

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.ciService.GetByEntity(entity, authInfo.UserLogin, limit)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetByEntityAndType
// @Summary получить КИ по entity и типу
// @Tags collection items
// @Accept  json
// @Produce  json
// @Param entity query string true "entity"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Param orderBy query string false "поле сортировки% может быть  name, purchaseDate, purchasePrice, platform, type, rating"
// @Param order query string false "порядок сортировки asc или desc"
// @Param search query string false "поиск среди коллекций по entity"
// @Param type query string false "type"
// @Success 201 {object} dto.AllCollectionsDataResponseSwagger
// @Router /public/collection_item/entity/type [get]
func (c *CIController) GetByEntityAndType(ctx echo.Context) error {
	entity := ctx.QueryParam("entity")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	orderBy := ctx.QueryParam("orderBy")
	order := ctx.QueryParam("order")
	search := ctx.QueryParam("search")
	ciType := ctx.QueryParam("type")
	if entity == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Entity parameter is required",
		})
	}

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.enrichedService.GetByEntityAndType(entity, ciType, limit, offset, search, orderBy, order, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
