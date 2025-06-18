package auth

import (
	"github.com/joho/godotenv"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"os"
)

type ValidationResponse struct {
	Status string `json:"status"`
}

type RecaptchaService struct {
}

func NewRecaptchaService() *RecaptchaService {
	_ = godotenv.Load()
	return &RecaptchaService{}
}

func (s *RecaptchaService) CheckCaptcha(token string) (*dto.CommonResponse, error) {
	secret := os.Getenv("SMART_CAPTCHA_SERVER_KEY")
	url := os.Getenv("SMART_CAPTCHA_URL")
	reqParams := map[string]string{}

	reqBody := map[string]string{
		"secret": secret,
		"token":  token,
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	response, err := utils.SendRequest[ValidationResponse](utils.RequestOptions{Method: "POST", URL: url, Headers: headers, QueryParams: reqParams, Body: reqBody, File: []byte{}, BasicAuth: nil})

	if err != nil {
		return nil, err
	}
	if response.Status == "ok" {
		return &dto.CommonResponse{Data: dto.Resp{
			Success: true,
		}}, nil
	}
	return &dto.CommonResponse{Data: dto.Resp{
		Success: false,
	}}, nil
}
