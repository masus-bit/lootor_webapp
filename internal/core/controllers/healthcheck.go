package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (c *HealthCheckController) Healthcheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
