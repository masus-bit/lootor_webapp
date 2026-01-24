package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/dto"
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
// @Success 201 {object} dto.TagsDataResponse
// @Router /secured/tags [get]
func (c *TagsController) SearchTags(ctx echo.Context) error {
	name := ctx.QueryParam("name")
	response, err := c.tagsService.SearchTags(name)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// SearchTagsV2
// @Summary Поиск по тегам
// @Tags tags
// @Accept  json
// @Produce  json
// @Param name query string true "name"
// @Param limit query string true "limit"
// @Success 201 {object} dto.TagsDataResponse
// @Router /secured/tags/search [get]
func (c *TagsController) SearchTagsV2(ctx echo.Context) error {
	name := ctx.QueryParam("name")
	limit := ctx.QueryParam("limit")
	response, err := c.tagsService.SearchSmartTags(name, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetSuggestions
// @Summary Получить предположения
// @Tags tags
// @Accept  json
// @Produce  json
// @Param name query string true "name"
// @Param limit query string true "limit"
// @Success 201 {object} dto.TagsDataResponse
// @Router /secured/tags/suggestions [get]
func (c *TagsController) GetSuggestions(ctx echo.Context) error {
	name := ctx.QueryParam("name")
	limit := ctx.QueryParam("limit")
	response, err := c.tagsService.GetSuggestions(name, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// RecordChoice
// @Summary Выбобр юзера записать
// @Tags tags
// @Accept  json
// @Produce  json
// @Param addRequest body dto.UserChoiceRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/record_choice [post]
func (c *TagsController) RecordChoice(ctx echo.Context) error {
	var request dto.UserChoiceRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.RecordUserChoice(&request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// CreateTag
// @Summary Создание тега
// @Tags tags
// @Accept  json
// @Produce  json
// @Param createRequest body dto.TagCreateRequest true "поля создания (entityID, entityType - необязательные)"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags [post]
func (c *TagsController) CreateTag(ctx echo.Context) error {
	var request dto.TagCreateRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
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
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}
	request.Author = authUser
	response, err := c.tagsService.CreateTag(&request, authUser)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
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
// @Success 201 {object} dto.TagDataResponse
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
	response, err := c.tagsService.FindAllEntitiesByTag(
		tagSlug,
		entityType,
		limit,
		offset,
		authInfo.UserLogin,
		ciFilter,
	)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// AddTagToEntity
// @Summary добавить тег к сущности
// @Tags tags
// @Accept  json
// @Produce  json
// @Param addRequest body dto.AddTagToEntityRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/entity [post]
func (c *TagsController) AddTagToEntity(ctx echo.Context) error {
	var request dto.AddTagToEntityRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	role, ok := ctx.Get("role").(string)
	if !ok {
		role = ""
	}
	if role != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}
	response, err := c.tagsService.AddTagToEntity(&request, authUser, role)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// RemoveTagsFromEntity
// @Summary удалить теги у сущностей
// @Tags tags
// @Accept  json
// @Produce  json
// @Param removeRequest body dto.RemoveTagsRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/entity/remove [post]
func (c *TagsController) RemoveTagsFromEntity(ctx echo.Context) error {
	var request dto.RemoveTagsRequest
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	role, ok := ctx.Get("role").(string)
	if !ok {
		role = ""
	}
	if role != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.RemoveTagsFromEntity(&request, authUser, role)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// MergeTags
// @Summary смержить теги к одному праймари тегу
// @Tags tags
// @Accept  json
// @Produce  json
// @Param mergeRequest body dto.MergeTagsRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/merge [post]
func (c *TagsController) MergeTags(ctx echo.Context) error {
	var request dto.MergeTagsRequest

	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.MergeTags(&request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// MergeSeries
// @Summary смержить теги к одному главному тегу серии
// @Tags tags
// @Accept  json
// @Produce  json
// @Param mergeRequest body dto.MergeTagsRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/merge_series [post]
func (c *TagsController) MergeSeries(ctx echo.Context) error {
	var request dto.MergeTagsRequest

	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.MergeSeries(&request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// DeleteTags
// @Summary смержить теги к одному главному тегу серии
// @Tags tags
// @Accept  json
// @Produce  json
// @Param mergeRequest body dto.DeleteTags true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/delete [post]
func (c *TagsController) DeleteTags(ctx echo.Context) error {
	var request dto.DeleteTags

	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.DeleteTags(&request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// UpdateTag
// @Summary обновить инфу по тегу
// @Tags tags
// @Accept  json
// @Produce  json
// @Param updateRequest body dto.TagUpdateRequest true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/update/{id} [post]
func (c *TagsController) UpdateTag(ctx echo.Context) error {
	var request dto.TagUpdateRequest

	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.UpdateTag(&request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetAllTags
// @Summary Получить все теги
// @Tags tags
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.TagsShortDataResponse
// @Router /secured/tags/all [get]
func (c *TagsController) GetAllTags(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	response, err := c.tagsService.GetAllTags(limit, offset, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// Subscribe
// @Summary подписка на тег
// @Tags tags
// @Accept  json
// @Produce  json
// @Param id path string true "tag id"
// @Success 201 {object} dto.CommonResponse
// @Router /public/tags/subscribe/{id} [get]
func (c *TagsController) Subscribe(ctx echo.Context) error {
	id := ctx.Param("id")
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.tagsService.Subscribe(id, authUser)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// MoveTagLinks
// @Summary переместить связи тега
// @Tags tags
// @Accept  json
// @Produce  json
// @Param mergeRequest body dto.MoveTagLinks true "необходимые поля"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/tags/adm/move [post]
func (c *TagsController) MoveTagLinks(ctx echo.Context) error {
	var request dto.MoveTagLinks
	userRole, ok := ctx.Get("role").(string)
	if !ok {
		userRole = ""
	}

	if userRole != "admin" {
		return ctx.JSON(
			http.StatusForbidden, map[string]string{
				"error": "Forbidden",
			},
		)
	}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.MoveTagLinks(&request)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// PublicUpdate
// @Summary Обновить доп поля
// @Tags tags
// @Accept  json
// @Produce  json
// @Param mergeRequest body dto.PublicUpdateTag true "необходимые поля"
// @Success 201 {object} dto.TagDataResponse
// @Router /secured/tags/update [put]
func (c *TagsController) PublicUpdate(ctx echo.Context) error {
	var request dto.PublicUpdateTag

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(
			http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			},
		)
	}
	response, err := c.tagsService.PublicUpdate(&request)

	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetOnModerateTags
// @Summary получить теги на модерации
// @Tags tags
// @Accept  json
// @Produce  json
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} dto.TagsDataResponse
// @Router /secured/tags/adm/moderate [get]
func (c *TagsController) GetOnModerateTags(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")

	response, err := c.tagsService.GetOnModerationTags(
		&dto.OnModerationTagsRequest{
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *TagsController) TriggerMoveTags(ctx echo.Context) error {
	err := c.tagsService.TriggerMoveTags()
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, "OK")
}
