package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/feedback"
	"net/http"
)

type FeedbackController struct {
	feedService *feedback.FBService
}

func NewFeedbackController(feedService *feedback.FBService) *FeedbackController {
	return &FeedbackController{feedService: feedService}
}

// AddIssue
// @Summary Создать обращение
// @Description Загрузка файла через multipart/form-data
// @Tags feedback
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Param title query string true "Заголовок тикета"
// @Param description query string true "Описание тикета"
// @Success 200 {object} dto.CommonResponse
// @Router /public/feedback [post]
func (c *FeedbackController) AddIssue(ctx echo.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid form data"})
	}

	file := form.File["file"]

	authInfo := ctx.Get("auth_info").(struct {
		IsAuthenticated bool
		UserLogin       string
	})

	title := ctx.QueryParam("title")
	description := ctx.QueryParam("description")

	if title == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "title обязателен"})
	}
	if description == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "description обязателен"})
	}

	var fileBuffer []byte
	if file == nil {
		fileBuffer = nil
	} else {
		src, err := file[0].Open()
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "file хуевый"})
		}
		defer src.Close()

		buf := make([]byte, file[0].Size)
		if _, err = src.Read(buf); err != nil {

		}
		fileBuffer = buf
	}
	response, err := c.feedService.AddIssue(title, description, fileBuffer, authInfo.UserLogin)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response)
}
