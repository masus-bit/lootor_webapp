package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/core/dto"
	"lootor/internal/core/services"
	"net/http"
)

type WLController struct {
	wlService services.WLService
}

func NewWLController(wlService services.WLService) *WLController {
	return &WLController{wlService: wlService}
}

// AddWishListItem
// @Summary Создание экземпляра вишлиста
// @Tags wishlist
// @Accept  json
// @Produce  json
// @Param createRequest body dto.WishListCreateRequest true "поля создания"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/wishlist [post]
func (c *WLController) AddWishListItem(ctx echo.Context) error {
	var request dto.WishListCreateRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)

	if !ok {
		authUser = ""
	}

	response, err := c.wlService.AddItem(&request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetWishList
// @Summary получить вишлист по юзер логину
// @Tags wishlist
// @Accept  json
// @Produce  json
// @Param userLogin query string true "login"
// @Success 201 {object} dto.WishlistDataResponseSwagger
// @Router /public/wishlist [get]
func (c *WLController) GetWishList(ctx echo.Context) error {
	login := ctx.QueryParam("userLogin")
	response, err := c.wlService.GetAllUserItems(login)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// UpdatePriority
// @Summary обновление приоритета
// @Tags wishlist
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Param priority body dto.WishListItemUpdatePriority true "priority"
// @Success 201 {object} dto.WishlistDataResponseSwagger
// @Router /secured/wishlist/priority [post]
func (c *WLController) UpdatePriority(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	var request dto.WishListItemUpdatePriority
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)

	if !ok {
		authUser = ""
	}

	response, err := c.wlService.UpdatePriority(id, request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Delete
// @Summary удаление айтема
// @Tags wishlist
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/wishlist [delete]
func (c *WLController) Delete(ctx echo.Context) error {
	id := ctx.QueryParam("id")

	response, err := c.wlService.DeleteItem(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Update
// @Summary обновление
// @Tags wishlist
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Param request body dto.WishListItemUpdatePriority true "priority"
// @Success 201 {object} dto.WishlistDataResponseSwagger
// @Router /secured/wishlist/update [post]
func (c *WLController) Update(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	var request dto.WishListUpdateRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)

	if !ok {
		authUser = ""
	}

	response, err := c.wlService.Update(id, request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
