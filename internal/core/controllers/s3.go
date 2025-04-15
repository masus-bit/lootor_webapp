package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/s3"

	"net/http"
	"strconv"
)

type FilesController struct {
	s3 *s3.S3Service
}

func NewFilesController(s3 *s3.S3Service) *FilesController {
	return &FilesController{s3: s3}
}

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
