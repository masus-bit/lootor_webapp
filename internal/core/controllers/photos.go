package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"lootor/internal/core/dto"
	"lootor/internal/core/services"
	"net/http"
)

type PhotosController struct {
	photosService services.PhotosService
}

func NewPhotosController(photosService services.PhotosService) *PhotosController {
	return &PhotosController{photosService: photosService}
}

// CreatePhotos
// @Summary Создать фотки
// @Tags photos
// @Accept  json
// @Produce  json
// @Param addRequest body models.PhotoCreate true "необходимые поля"
// @Success 201 {object} models.PhotosDataResponse
// @Router /secured/photos [post]
func (c *PhotosController) CreatePhotos(ctx echo.Context) error {
	var request dto.PhotoCreate
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	req := dto.PhotoCreateRequest{
		CollectionID: request.CollectionID,
		Paths:        request.Paths,
		Author:       authUser,
	}
	response, err := c.photosService.CreatePhoto(context.Background(), &req, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetByUserLogin
// @Summary получить все фото по логину юзера
// @Tags photos
// @Accept  json
// @Produce  json
// @Param userLogin query string true "user login"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} models.PhotosDataResponse
// @Router /public/photos/by_user [get]
func (c *PhotosController) GetByUserLogin(ctx echo.Context) error {
	userLogin := ctx.QueryParam("userLogin")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.photosService.GetByUser(context.Background(), userLogin, authInfo.UserLogin, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Delete
// @Summary удаление photo
// @Tags photos
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/photos/{id} [delete]
func (c *PhotosController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	response, err := c.photosService.DeletePhoto(context.Background(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetByID
// @Summary получить по ид
// @Tags photos
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} models.PhotoDataResponse
// @Router /public/photos/{id} [get]
func (c *PhotosController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.photosService.FindOneByID(context.Background(), id, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// FindAllByCollectionID
// @Summary получить по ид collection
// @Tags photos
// @Accept  json
// @Produce  json
// @Param collectionId query string true "collectionId"
// @Param limit query string true "limit"
// @Param offset query string true "offset"
// @Success 201 {object} models.PhotosDataResponse
// @Router /public/photos/by_collection [get]
func (c *PhotosController) FindAllByCollectionID(ctx echo.Context) error {
	collectionId := ctx.QueryParam("collectionId")
	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	response, err := c.photosService.FindAllByCollectionID(context.Background(), collectionId, limit, offset, authInfo.UserLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Like
// @Summary like
// @Tags photos
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 201 {object} dto.CommonResponse
// @Router /secured/photos/like/{id} [get]
func (c *PhotosController) Like(ctx echo.Context) error {
	id := ctx.Param("id")

	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.photosService.LikePhoto(context.Background(), id, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}

// Update
// @Summary update
// @Tags photos
// @Accept  json
// @Produce  json
// @Param addRequest body models.PhotoUpdateRequest true "необходимые поля"
// @Success 201 {object} models.PhotoDataResponse
// @Router /secured/photos [patch]
func (c *PhotosController) Update(ctx echo.Context) error {
	var request dto.PhotoUpdateRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	response, err := c.photosService.UpdatePhoto(context.Background(), &request, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, response)
}
