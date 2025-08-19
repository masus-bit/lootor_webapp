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

// SignUp  создает нового пользователя
// @Summary создание пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.SignUpRequest true "User data"
// @Success 201 {object} models.SignUpResponse
// @Router /public/auth/signup [post]
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

// SignIn Авторизация
// @Summary Авторизация
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.SignInRequest true "User data"
// @Success 201 {object} models.SignInResponse
// @Router /public/auth/signin [post]
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

// GetByLogin
// @Summary получение юзера по логину
// @Tags users
// @Accept  json
// @Produce  json
// @Param login query string true "логин"
// @Success 201 {object} dto.DataUserResponseSwagger
// @Router /public/user [get]
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

// GetAll
// @Summary получение юзеров
// @Tags users
// @Accept  json
// @Produce  json
// @Param offset query string true "offset"
// @Param limit query string true "limit"
// @Param search query string true "search"
// @Param order query string true "order"
// @Success 201 {object} dto.DataUserResponseSwagger
// @Router /public/user/all [get]
func (c *UserController) GetAll(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	order := ctx.QueryParam("order")

	response, err := c.userService.GetAll(search, limit, offset, order)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// ChangeRating
// @Summary Лайк/дизлайк
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.LikeUserSwagger true "User data"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/user/rating [post]
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

// ChangePass
// @Summary Смена пароля
// @Tags users
// @Accept  json
// @Produce  json
// @Param password body string true "parolj"
// @Param login query string true "User data"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/user/password [post]
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

// Subscribe
// @Summary Подписка одного юзера на другого
// @Tags users
// @Accept  json
// @Produce  json
// @Param isSubscribe query boolean true "признак подписки"
// @Param login query string true "User data"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/user/subscriptions [post]
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
	codeVerifier := ctx.QueryParam("code_verifier")
	deviceId := ctx.QueryParam("device_id")
	state := ctx.QueryParam("state")
	request := models.VkOauthRequest{
		Code:         code,
		CodeVerifier: codeVerifier,
		DeviceId:     deviceId,
		State:        state,
	}
	switch "" {
	case code:
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Code parameter is required",
		})
	case codeVerifier:
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "CodeVerifier parameter is required",
		})
	case deviceId:
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "DeviceId parameter is required",
		})
	case state:
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "State parameter is required",
		})
	}
	response, err := c.userService.VkOauth(&request)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) TelegramOauth(ctx echo.Context) error {
	var request models.TelegramOauthRequest

	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error бля": err.Error(),
		})
	}

	response, err := c.userService.TelegramOauth(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error бля": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)

}

// ResetPassword
// @Summary сброс пароля
// @Tags users
// @Accept  json
// @Produce  json
// @Param email query string true "email"
// @Success 201 {object} dto.CommonResponse
// @Router /public/auth/reset [get]
func (c *UserController) ResetPassword(ctx echo.Context) error {
	email := ctx.QueryParam("email")

	if email == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "email parameter is required",
		})
	}

	response, err := c.userService.ResetPassword(email)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// ChangeResetPassword
// @Summary смена пароля по токену
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.ChangePasswordReset true "User data"
// @Success 201 {object} dto.CommonResponse
// @Router /public/user/reset [post]
func (c *UserController) ChangeResetPassword(ctx echo.Context) error {
	var request models.ChangePasswordReset
	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	response, err := c.userService.ChangeResetPassword(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// UpdateUser
// @Summary смена некоторых данных пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UserRequestUpdate true "User data"
// @Success 201 {object} dto.DataUserResponseSwagger
// @Router /secured/user/update [post]
func (c *UserController) UpdateUser(ctx echo.Context) error {
	var request models.UserRequestUpdate
	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.userService.UpdateUser(&request, authUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// UpdateLogin
// @Summary смена логина после оауса
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UserRequestUpdateFirstTime true "User data"
// @Success 201 {object} models.SignInResponse
// @Router /secured/user/update/reg [post]
func (c *UserController) UpdateLogin(ctx echo.Context) error {
	var request models.UserRequestUpdateFirstTime
	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.userService.UpdateOnlyOnce(&request, authUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// DeleteUser
// @Summary удаление пользователя
// @Tags users
// @Accept  json
// @Produce  json
// @Success 201 {object} dto.CommonResponse
// @Router /secured/user/delete [get]
func (c *UserController) DeleteUser(ctx echo.Context) error {
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.userService.DeleteUser(authUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
