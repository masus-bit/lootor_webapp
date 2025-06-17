package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/feedback"
	"net/http"
)

type ReportsController struct {
	reportsService *feedback.ReportsService
}

func NewReportsController(reportsService *feedback.ReportsService) *ReportsController {
	return &ReportsController{reportsService: reportsService}
}

// Report
// @Summary подать жалобу
// @Tags reports
// @Accept  json
// @Produce  json
// @Param id query string true "id/login"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/reports [get]
func (c *ReportsController) Report(ctx echo.Context) error {

	id := ctx.QueryParam("id")

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "id обязателен"})
	}
	response, err := c.reportsService.ReportAnything(id)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response)
}
