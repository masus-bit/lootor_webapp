package controllers

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/elasticsearch"
)

type SearchController struct {
	es *elasticsearch.ElasticService
}

func NewSearchController(es *elasticsearch.ElasticService) *SearchController {
	return &SearchController{es: es}
}

// Search
// @Summary Поиск по коллекциям, ентити, КИ, тегам, юзерам
// @Tags search
// @Accept  json
// @Produce  json
// @Param search query string true "поисковая строка"
// @Param limit query string true "количество результатов"
// @Param offset query string true "ebal"
// @Success 201 {object} elasticsearch.SearchResult
// @Router /public/search [get]
func (c *SearchController) Search(ctx echo.Context) error {
	query := ctx.QueryParam("search")
	if query == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "search query is обязателен"})
	}
	limit := ctx.QueryParam("limit")
	if limit == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "limit query is обязателей"})
	}
	offset := ctx.QueryParam("offset")
	if offset == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "offset query is обязателей"})
	}
	searchType := ctx.QueryParam("type")

	indices := []string{"users", "tags", "collection_items", "collections"}

	var finalIndices []string

	if searchType != "" {
		index := slices.Index(indices, searchType)
		finalIndices = append(finalIndices, indices[index])
	} else {
		finalIndices = indices
	}

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	result, err := c.es.SearchInIndices(ctx.Request().Context(), finalIndices, query, intLimit, intOffset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, result)
}
