package feedback

import (
	"fmt"
	"github.com/joho/godotenv"
	"lootor/internal/core/dto"
	"lootor/internal/pkg/utils"
	"os"
)

type FBService struct {
}

func NewFeedbackService() *FBService {
	_ = godotenv.Load()
	return &FBService{}
}

func (s *FBService) AddIssue(title string, description string, file []byte, authUserLogin, commMethod string) (
	*dto.CommonResponse,
	error,
) {

	if authUserLogin != "" {
		description = description + "\n\n\n__________\n\nЛогин пользователя: " + authUserLogin
	}

	if commMethod != "" {
		description = description + "\n\n\n__________\n\nСпособ связи: " + commMethod
	}

	response, err := s.SendTextTicket(title, description, false)
	if err != nil {
		return nil, err
	}

	if file != nil {

		attachReqURL := fmt.Sprintf(
			"%s/%s/attachments",
			os.Getenv("TRELLO_API_URL"),
			response.ID,
		)

		attachReqParams := map[string]string{
			"key":   os.Getenv("TRELLO_API_KEY"),
			"token": os.Getenv("TRELLO_TOKEN"),
		}

		_, attachErr := utils.SendRequest[any](
			utils.RequestOptions{
				Method:      "POST",
				URL:         attachReqURL,
				Headers:     map[string]string{},
				QueryParams: attachReqParams,
				Body:        map[string]string{},
				File:        file,
				BasicAuth:   nil,
			},
		)

		if attachErr != nil {
			return nil, fmt.Errorf("failed to request attach: %w", err)
		}
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *FBService) SendTextTicket(title string, description string, isReports bool) (
	*struct {
		ID string `json:"id"`
	}, error,
) {

	var listID string
	if isReports {
		listID = os.Getenv("TRELLO_REPORTS_LIST_ID")
	} else {
		listID = os.Getenv("TRELLO_LIST_ID")
	}

	reqParams := map[string]string{
		"idList": listID,
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
		ID string `json:"id"`
	}](
		utils.RequestOptions{
			Method:      "POST",
			URL:         os.Getenv("TRELLO_API_URL"),
			Headers:     headers,
			QueryParams: reqParams,
			Body:        reqBody,
			File:        []byte{},
			BasicAuth:   nil,
		},
	)

	if err != nil {
		return nil, err
	}
	return &response, nil
}
