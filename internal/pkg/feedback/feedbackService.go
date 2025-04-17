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

	//attachReqHeaders := map[string]string{
	//	"Content-Type": writer.FormDataContentType(),
	//}

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

	//var requestBody bytes.Buffer
	//writer := multipart.NewWriter(&requestBody)
	//
	//part, err := writer.CreateFormFile("file", uuid.New().String())
	//if err != nil {
	//	return nil, err
	//}
	//
	//_, err = io.Copy(part, bytes.NewReader(file))
	//
	//writer.Close()
	//if err != nil {
	//	return nil, fmt.Errorf("error: %s", string(body))
	//}
	//
	//resUrl := fmt.Sprintf("%s/%s/attachments?key=%s&token=%s",
	//	os.Getenv("TRELLO_API_URL"),
	//	response.Id,
	//	os.Getenv("TRELLO_API_KEY"),
	//	os.Getenv("TRELLO_TOKEN"),
	//)
	//
	//request, err := http.NewRequest("POST", resUrl, &requestBody)
	if attachErr != nil {
		return nil, fmt.Errorf("failed to request attach: %w", err)
	}

	//request.Header.Set("Content-Type", writer.FormDataContentType())

	//client = &http.Client{}
	//resp, err = client.Do(request)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to send request: %w", err)
	//}
	//defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	body, _ := io.ReadAll(resp.Body)
	//	return nil, fmt.Errorf("trello api error: %s, response: %s", resp.Status, string(body))
	//}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}
