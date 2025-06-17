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

	description = description + "\n\n\n__________\n\nЛогин пользователя: " + authUserLogin
	response, err := s.SendTextTicket(title, description, false)
	if err != nil {
		return nil, err
	}

	if file != nil {

		attachReqUrl := fmt.Sprintf("%s/%s/attachments",
			os.Getenv("TRELLO_API_URL"),
			response.Id,
		)

		attachReqParams := map[string]string{
			"key":   os.Getenv("TRELLO_API_KEY"),
			"token": os.Getenv("TRELLO_TOKEN"),
		}

		_, attachErr := utils.SendRequest[any](utils.RequestOptions{Method: "POST", URL: attachReqUrl, Headers: map[string]string{}, QueryParams: attachReqParams, Body: map[string]string{}, File: file, BasicAuth: nil})

		if attachErr != nil {
			return nil, fmt.Errorf("failed to request attach: %w", err)
		}
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *FeedbackService) SendTextTicket(title string, description string, isReports bool) (*struct {
	Id string `json:"id"`
}, error) {

	var listId string
	if isReports {
		listId = os.Getenv("TRELLO_REPORTS_LIST_ID")
	} else {
		listId = os.Getenv("TRELLO_LIST_ID")
	}

	reqParams := map[string]string{
		"idList": listId,
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
	}](utils.RequestOptions{Method: "POST", URL: os.Getenv("TRELLO_API_URL"), Headers: headers, QueryParams: reqParams, Body: reqBody, File: []byte{}, BasicAuth: nil})

	if err != nil {
		return nil, err
	}
	return &response, nil
}
