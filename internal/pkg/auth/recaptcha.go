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
}

func NewRecaptchaService() *RecaptchaService {
	_ = godotenv.Load()
	return &RecaptchaService{}
}

func (s *RecaptchaService) CheckRecaptcha(token string) (*dto.CommonResponse, error) {
	secret := os.Getenv("RECAPTCHA_SECRET")
	projectID := os.Getenv("RECAPTCHA_PROJECT_ID")
	recaptchaKey := secret
	recaptchaAction := "login"

	success := createAssessment(projectID, recaptchaKey, token, recaptchaAction)

	return &dto.CommonResponse{Data: dto.Resp{Success: success}}, nil
}

func createAssessment(projectID string, recaptchaKey string, token string, recaptchaAction string) bool {

	ctx := context.Background()
	client, err := recaptcha.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error creating reCAPTCHA client\n")
	}
	defer client.Close()

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

	response, err := client.CreateAssessment(
		ctx,
		request)

	if err != nil {
		fmt.Printf("Error calling CreateAssessment: %v", err.Error())
	}

	if !response.TokenProperties.Valid {
		fmt.Printf("The CreateAssessment() call failed because the token was invalid for the following reasons: %v",
			response.TokenProperties.InvalidReason)
		return false
	}

	if response.TokenProperties.Action != recaptchaAction {
		fmt.Printf("The action attribute in your reCAPTCHA tag does not match the action you are expecting to score")
		return false
	}

	fmt.Printf("The reCAPTCHA score for this token is:  %v", response.RiskAnalysis.Score)

	if response.RiskAnalysis.Score < 0.5 {
		return false
	}
	return true

	//for _, reason := range response.RiskAnalysis.Reasons {
	//	fmt.Printf(reason.String() + "\n")
	//}
}
