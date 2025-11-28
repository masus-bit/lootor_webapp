package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/s3"
	"net/http"
	"strconv"
)

type FilesController struct {
	s3 *s3.StorageService
}

func NewFilesController(s3 *s3.StorageService) *FilesController {
	return &FilesController{s3: s3}
}

// Upload
// @Summary Загрузить изображения в хранилище
// @Tags images
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "файлы для загрузки"
// @Param width query string false "ширина картинки"
// @Param height query string false "высота"
// @Success 200 {object} dto.ImagesResponse
// @Router /secured/images/upload [post]
func (c *FilesController) Upload(ctx echo.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid form data"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "no files provided"})
	}

	width, _ := strconv.Atoi(ctx.QueryParam("width"))
	height, _ := strconv.Atoi(ctx.QueryParam("height"))

	var fileBuffers [][]byte
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		buf := make([]byte, file.Size)
		if _, err = src.Read(buf); err != nil {
			continue
		}
		fileBuffers = append(fileBuffers, buf)
	}

	keys, err := c.s3.UploadOptimizedImages(ctx.Request().Context(), fileBuffers, s3.UploadOptions{
		Width:   width,
		Height:  height,
		Quality: 70,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"keys": keys})
}

// GetFile
// @Summary получить изображение
// @Tags images
// @Produce image/png
// @Produce image/jpeg
// @Produce image/webp
// @Param key path string true "id"
// @Success 200 {file} byte "image"
// @Router /public/images/{key} [get]
func (c *FilesController) GetFile(ctx echo.Context) error {
	key := ctx.Param("key")
	if key == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "key is required"})
	}

	file, err := c.s3.GetFile(ctx.Request().Context(), key)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return ctx.Blob(http.StatusOK, "image/webp", file)
}

// Delete
// @Summary удалить изображения из хранилища
// @Tags images
// @Accept json
// @Produce json
// @Param createRequest body s3.DeleteFilesRequest true "поля"
// @Success 200 {object} s3.DeleteFilesResponse
// @Router /secured/images/delete [post]
func (c *FilesController) Delete(ctx echo.Context) error {
	var req s3.DeleteFilesRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if len(req.Keys) == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "at least one key required"})
	}

	result, err := c.s3.DeleteFiles(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, result)
}

// GetThumbnail
// @Summary получить миниатюру изображения
// @Tags images
// @Produce image/webp
// @Param filename path string true "имя файла"
// @Param width query int false "ширина"
// @Param height query int false "высота"
// @Param quality query int false "качество (1-100)"
// @Success 200 {file} byte "image"
// @Router /public/thumbs/{key} [get]
func (c *FilesController) GetThumbnail(ctx echo.Context) error {
	filename := ctx.Param("key")
	if filename == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "filename is required"})
	}

	width, _ := strconv.Atoi(ctx.QueryParam("width"))
	height, _ := strconv.Atoi(ctx.QueryParam("height"))
	quality, _ := strconv.Atoi(ctx.QueryParam("quality"))

	if width == 0 && height == 0 {
		width = 300
	}
	if quality == 0 {
		quality = 85
	}

	thumbData, err := c.s3.GenerateThumbnail(
		ctx.Request().Context(),
		filename,
		width,
		height,
		quality,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("thumbnail generation failed: %v", err),
		})
	}

	ctx.Response().Header().Set("Cache-Control", "public, max-age=31536000")
	return ctx.Blob(http.StatusOK, "image/webp", thumbData)
}
