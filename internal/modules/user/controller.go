package user

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	userService Service
}

func NewController(userService Service) *Controller {
	return &Controller{userService: userService}
}

func (c *Controller) SignUp(ctx echo.Context) error {
	var request SignUpRequest

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if _, err := c.userService.SignUp(&request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, request)
}

func (c *Controller) SignIn(ctx echo.Context) error {
	var request SignInRequest

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	response, err := c.userService.SignIn(&request)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetByLogin(ctx echo.Context) error {
	login := ctx.Param("login")
	if login == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Login parameter is required",
		})
	}

	response, err := c.userService.GetByLogin(login)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
