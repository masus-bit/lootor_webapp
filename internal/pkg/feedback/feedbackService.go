package feedback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io"
	"lootor/internal/pkg/dto"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type FeedbackService struct {
}

func NewFeedbackService() *FeedbackService {
	_ = godotenv.Load()
	return &FeedbackService{}
}

func (s *FeedbackService) AddIssue(title string, description string, file []byte, authUserLogin string) (*dto.CommonResponse, error) {

	description = description + "\n\bЛогин пользователя: " + authUserLogin

	params := url.Values{}
	params.Add("idList", os.Getenv("TRELLO_LIST_ID"))
	params.Add("key", os.Getenv("TRELLO_API_KEY"))
	params.Add("token", os.Getenv("TRELLO_TOKEN"))

	bodyReq := url.Values{}
	bodyReq.Add("name", title)
	bodyReq.Add("desc", description)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("TRELLO_API_URL")+"?"+params.Encode(),
		strings.NewReader(bodyReq.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("trello request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("trello error: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)

	var cardId = struct {
		Id string `json:"id"`
	}{}

	err = json.Unmarshal(body, &cardId)

	if err != nil {
		return nil, err
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", uuid.New().String())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(file))

	writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	resUrl := fmt.Sprintf("%s/%s/attachments?key=%s&token=%s",
		os.Getenv("TRELLO_API_URL"),
		cardId.Id,
		os.Getenv("TRELLO_API_KEY"),
		os.Getenv("TRELLO_TOKEN"),
	)

	request, err := http.NewRequest("POST", resUrl, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client = &http.Client{}
	resp, err = client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("trello api error: %s, response: %s", resp.Status, string(body))
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}
