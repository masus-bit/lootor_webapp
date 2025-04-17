package controllers

import (
	"github.com/labstack/echo/v4"
	"lootor/internal/pkg/feedback"
	"net/http"
)

type FeedbackController struct {
	feedService *feedback.FeedbackService
}

func NewFeedbackController(feedService *feedback.FeedbackService) *FeedbackController {
	return &FeedbackController{feedService: feedService}
}

func (c *FeedbackController) AddIssue(ctx echo.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid form data"})
	}

	file := form.File["file"]
	if file == nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "no files provided"})
	}
	authUser, ok := ctx.Get("user_login").(string)
	if !ok {
		authUser = ""
	}

	title := ctx.QueryParam("title")
	description := ctx.QueryParam("description")

	if title == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "title обязателен"})
	}
	if description == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "description обязателен"})
	}

	var fileBuffer []byte
	src, err := file[0].Open()
	if err != nil {
	}
	defer src.Close()

	buf := make([]byte, file[0].Size)
	if _, err = src.Read(buf); err != nil {

	}
	fileBuffer = buf
	response, err := c.feedService.AddIssue(title, description, fileBuffer, authUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response)
}
