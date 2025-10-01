package controllers

import (
	"net/http"

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

	indices := []string{"users", "tags", "collection_items", "collections"}
	result, err := c.es.SearchInIndices(ctx.Request().Context(), indices, query, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, result)
}
