package auth

import (
	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"lootor/internal/pkg/dto"
	"os"
)

type RecaptchaService struct {
	client *recaptcha.Client
}

func NewRecaptchaService() (*RecaptchaService, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	ctx := context.Background()
	client, err := recaptcha.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reCAPTCHA client: %v", err)
	}

	return &RecaptchaService{client: client}, nil
}

func (s *RecaptchaService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *RecaptchaService) CheckRecaptcha(token, action string) (*dto.CommonResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("reCAPTCHA client not initialized")
	}

	secret := os.Getenv("RECAPTCHA_SECRET")
	projectID := os.Getenv("RECAPTCHA_PROJECT_ID")
	if secret == "" || projectID == "" {
		return nil, fmt.Errorf("reCAPTCHA configuration missing")
	}

	success, err := s.createAssessment(projectID, secret, token, action)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: success}}, nil
}

func (s *RecaptchaService) createAssessment(projectID, recaptchaKey, token, recaptchaAction string) (bool, error) {
	ctx := context.Background()

	event := &recaptchapb.Event{
		Token:   token,
		SiteKey: recaptchaKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
	}

	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", projectID),
	}

	response, err := s.client.CreateAssessment(ctx, request)
	if err != nil {
		return false, fmt.Errorf("error calling CreateAssessment: %v", err)
	}

	if response == nil {
		return false, fmt.Errorf("empty response from reCAPTCHA API")
	}

	if response.TokenProperties == nil {
		return false, fmt.Errorf("missing token properties in response")
	}

	if !response.TokenProperties.Valid {
		return false, fmt.Errorf("invalid token: %v", response.TokenProperties.InvalidReason)
	}

	if response.TokenProperties.Action != recaptchaAction {
		return false, fmt.Errorf("action mismatch: expected %s, got %s",
			recaptchaAction, response.TokenProperties.Action)
	}

	if response.RiskAnalysis == nil {
		return false, fmt.Errorf("missing risk analysis in response")
	}

	score := response.RiskAnalysis.Score
	fmt.Printf("The reCAPTCHA score for this token is: %v\n", score)

	return score >= 0.5, nil
}
