package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
	"lootor/internal/core/services"
	"net/http"
)

type TagsController struct {
	tagsService services.TagsService
}

func NewTagsController(tagsService services.TagsService) *TagsController {
	return &TagsController{tagsService: tagsService}
}

// SearchTags
// @Summary Поиск по тегам
// @Tags tags
// @Accept  json
// @Produce  json
// @Param name query string true "name"
// @Success 201 {object} models.TagsDataResponse
// @Router /secured/tags [get]
func (c *TagsController) SearchTags(ctx echo.Context) error {
	name := ctx.QueryParam("name")
	response, err := c.tagsService.SearchTags(name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// CreateTag
// @Summary Создание тега
// @Tags tags
// @Accept  json
// @Produce  json
// @Param createRequest body models.TagCreateRequest true "поля создания (entityID, entityType - необязательные)"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags [post]
func (c *TagsController) CreateTag(ctx echo.Context) error {
	var request models.TagCreateRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	if request.Author != "" && userRole != "admin" {
		return ctx.JSON(http.StatusForbidden, map[string]string{
			"error": "Forbidden",
		})
	}
	request.Author = authUser
	response, err := c.tagsService.CreateTag(&request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// FindEntitiesByTag
// @Summary получить все сущности по slug тегa с фильтром
// @Tags tags
// @Accept  json
// @Produce  json
// @Param id path string true "tag id"
// @Param entityType query string true "entity type"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Param ciFilter query string true "filter"
// @Success 201 {object} models.TagDataResponse
// @Router /public/tags/{id} [get]
func (c *TagsController) FindEntitiesByTag(ctx echo.Context) error {
	tagSlug := ctx.Param("id")
	entityType := ctx.QueryParam("entityType")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	ciFilter := ctx.QueryParam("ciFilter")
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	response, err := c.tagsService.FindAllEntitiesByTag(tagSlug, entityType, limit, offset, authInfo.UserLogin, ciFilter)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// AddTagToEntity
// @Summary добавить тег к сущности
// @Tags tags
// @Accept  json
// @Produce  json
// @Param addRequest body models.AddTagToEntityRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/entity [post]
func (c *TagsController) AddTagToEntity(ctx echo.Context) error {
	var request models.AddTagToEntityRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	role, ok := ctx.Get("role").(string)
	if !ok {
		role = ""
	}
	response, err := c.tagsService.AddTagToEntity(&request, authUser, role)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// RemoveTagsFromEntity
// @Summary удалить теги у сущностей
// @Tags tags
// @Accept  json
// @Produce  json
// @Param removeRequest body models.RemoveTagsRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/entity/remove [post]
func (c *TagsController) RemoveTagsFromEntity(ctx echo.Context) error {
	var request models.RemoveTagsRequest
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	role, ok := ctx.Get("role").(string)
	if !ok {
		role = ""
	}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	response, err := c.tagsService.RemoveTagsFromEntity(&request, authUser, role)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// MergeTags
// @Summary смержить теги к одному праймари тегу
// @Tags tags
// @Accept  json
// @Produce  json
// @Param mergeRequest body models.MergeTagsRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/merge [post]
func (c *TagsController) MergeTags(ctx echo.Context) error {
	var request models.MergeTagsRequest

	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(http.StatusForbidden, map[string]string{
			"error": "Forbidden",
		})
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	response, err := c.tagsService.MergeTags(&request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// UpdateTag
// @Summary обновить инфу по тегу
// @Tags tags
// @Accept  json
// @Produce  json
// @Param updateRequest body models.TagUpdateRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/update/{id} [post]
func (c *TagsController) UpdateTag(ctx echo.Context) error {
	var request models.TagUpdateRequest

	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(http.StatusForbidden, map[string]string{
			"error": "Forbidden",
		})
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	response, err := c.tagsService.UpdateTag(&request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
