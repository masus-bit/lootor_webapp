package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/auth"
	"net/http"
)

type RecaptchaController struct {
	recaptchaService *auth.RecaptchaService
}

func NewRecaptchaController(recaptchaService *auth.RecaptchaService) *RecaptchaController {
	return &RecaptchaController{recaptchaService: recaptchaService}
}

func (c *RecaptchaController) CheckRecaptcha(ctx echo.Context) error {
	token := ctx.QueryParam("token")

	response, err := c.recaptchaService.CheckRecaptcha(token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response)
}
