package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/models"
	"lootor/internal/core/services"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) SignUp(ctx echo.Context) error {
	var request models.SignUpRequest

	reqCtx := ctx.Request().Context()

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if _, err := c.userService.SignUp(&request, reqCtx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, request)
}

func (c *UserController) SignIn(ctx echo.Context) error {
	var request models.SignInRequest

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

func (c *UserController) GetByLogin(ctx echo.Context) error {
	login := ctx.QueryParam("login")
	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})
	if login == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Login parameter is required",
		})
	}

	response, err := c.userService.GetByLogin(login, authInfo.UserLogin, authInfo.IsAuthenticated)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) ChangeRating(ctx echo.Context) error {
	var request models.ChangeRatingRequest
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}
	login := ctx.QueryParam("login")
	if login == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Login parameter is required",
		})
	}

	response, err := c.userService.ChangeRating(request.IsLike, login)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) ChangePass(ctx echo.Context) error {
	var request models.ChangePasswordRequest
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}
	login := ctx.QueryParam("login")
	if login == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Login parameter is required",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.userService.ChangePassword(request.Password, login, authUser)

	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) Subscribe(ctx echo.Context) error {
	login := ctx.QueryParam("login")
	isSubscribe := ctx.QueryParam("isSubscribe")
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}
	response, err := c.userService.Subscribe(login, authUser, isSubscribe == "true")
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) Verification(ctx echo.Context) error {
	token := ctx.QueryParam("verificationToken")

	response, err := c.userService.Verification(token)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) RenewTokens(ctx echo.Context) error {
	var request models.RenewTokensRequest
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}
	response, err := c.userService.RefreshTokens(request.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) VkOauth(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	codeVerifier := ctx.QueryParam("codeVerifier")
	deviceId := ctx.QueryParam("deviceId")
	state := ctx.QueryParam("state")
	codeChallenge := ctx.QueryParam("codeChallenge")
	request := models.VkOauthRequest{
		Code:          code,
		CodeVerifier:  codeVerifier,
		DeviceId:      deviceId,
		State:         state,
		CodeChallenge: codeChallenge,
	}
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}
	response, err := c.userService.VkOauth(&request)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
