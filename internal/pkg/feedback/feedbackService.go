package feedback

import (
	"fmt"
	"github.com/joho/godotenv"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"os"
)

type FeedbackService struct {
}

func NewFeedbackService() *FeedbackService {
	_ = godotenv.Load()
	return &FeedbackService{}
}

func (s *FeedbackService) AddIssue(title string, description string, file []byte, authUserLogin string) (*dto.CommonResponse, error) {

	description = description + "\n\bЛогин пользователя: " + authUserLogin

	reqParams := map[string]string{
		"idList": os.Getenv("TRELLO_LIST_ID"),
		"key":    os.Getenv("TRELLO_API_KEY"),
		"token":  os.Getenv("TRELLO_TOKEN"),
	}

	reqBody := map[string]string{
		"name": title,
		"desc": description,
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	response, err := utils.SendRequest[struct {
		Id string `json:"id"`
	}](struct {
		Method      string
		URL         string
		Headers     map[string]string
		QueryParams map[string]string
		Body        map[string]string
		File        []byte
	}{Method: "POST", URL: os.Getenv("TRELLO_API_URL"), Headers: headers, QueryParams: reqParams, Body: reqBody, File: []byte{}})

	attachReqUrl := fmt.Sprintf("%s/%s/attachments",
		os.Getenv("TRELLO_API_URL"),
		response.Id,
	)

	attachReqParams := map[string]string{
		"key":   os.Getenv("TRELLO_API_KEY"),
		"token": os.Getenv("TRELLO_TOKEN"),
	}

	_, attachErr := utils.SendRequest[any](struct {
		Method      string
		URL         string
		Headers     map[string]string
		QueryParams map[string]string
		Body        map[string]string
		File        []byte
	}{Method: "POST", URL: attachReqUrl, Headers: map[string]string{}, QueryParams: attachReqParams, Body: map[string]string{}, File: file})

	if attachErr != nil {
		return nil, fmt.Errorf("failed to request attach: %w", err)
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}
