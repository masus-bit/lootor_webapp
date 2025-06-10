package auth

import (
	"github.com/joho/godotenv"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"os"
)

type RecaptchaService struct {
}

func NewRecaptchaService() *RecaptchaService {
	_ = godotenv.Load()
	return &RecaptchaService{}
}

func (s *RecaptchaService) CheckRecaptcha(token string) (*dto.CommonResponse, error) {
	baseUrl := os.Getenv("RECAPTCHA_URL")
	secret := os.Getenv("RECAPTCHA_SECRET")

	reqBody := map[string]string{
		"secret":   secret,
		"response": token,
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := utils.SendRequest[struct {
		Data struct {
			Success bool    `json:"success"`
			Score   float64 `json:"score"`
		}
	}](struct {
		Method      string
		URL         string
		Headers     map[string]string
		QueryParams map[string]string
		Body        map[string]string
		File        []byte
	}{Method: "POST", URL: baseUrl, Headers: headers, QueryParams: nil, Body: reqBody, File: []byte{}})

	if err != nil {
		return nil, err
	}
	if response.Data.Success && response.Data.Score > 0.5 {
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}
	return &dto.CommonResponse{Data: dto.Resp{
		Success: false,
	}}, nil
}
