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

func (c *SearchController) Search(ctx echo.Context) error {
	query := ctx.QueryParam("search")
	if query == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "search query is required"})
	}

	indices := []string{"users", "tags", "collection_items", "collections", "entities"}
	result, err := c.es.SearchInIndices(ctx.Request().Context(), indices, query)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, result)
}
