package services

import (
	"context"
	"fmt"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/achievementsclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/utils"
)

type AchievementsService struct {
	achievementsClient *achievementsclient.GRPCAchievementsClient
	userRepo           *repositories.UsersRepository
	tagsClient         *tagsclient.GRPCTagsClient
	collectionRepo     *repositories.CollectionsRepository
	eventsService      *EventsService
}

func NewAchievementsService(achievementsClient *achievementsclient.GRPCAchievementsClient, userRepo *repositories.UsersRepository, tagsClient *tagsclient.GRPCTagsClient, collectionRepo *repositories.CollectionsRepository, eventsService *EventsService) *AchievementsService {

	return &AchievementsService{
		achievementsClient: achievementsClient,
		userRepo:           userRepo,
		tagsClient:         tagsClient,
		collectionRepo:     collectionRepo,
		eventsService:      eventsService,
	}
}

func (s *AchievementsService) AddAchievement(code, userLogin string, level, xp int64) error {
	achievement := &microservices.AddOrUpdateAchievementRequest{
		Code:      code,
		UserLogin: userLogin,
		Xp:        xp,
		Level:     level,
	}
	_, err := s.achievementsClient.AddOrUpdateAchievement(context.Background(), achievement)
	if err != nil {
		return fmt.Errorf("AddAchievement: %w", err)
	}
	return nil
}

func (s *AchievementsService) GetUserAchievements(userLogin string) (*models.AchievementsResponse, error) {
	var achievementsResult []models.AchievedAchievement
	resp, err := s.achievementsClient.GetUserAchievements(context.Background(), &microservices.GetUserAchievementsRequest{
		UserLogin: userLogin,
	})
	if err != nil {
		return nil, fmt.Errorf("GetUserAchievements: %w", err)
	}

	for _, a := range resp.GetAchievedAchievements() {
		temp := s.convertAchievement(a)
		achievementsResult = append(achievementsResult, *temp)
	}
	return &models.AchievementsResponse{Data: achievementsResult}, nil
}

func (s *AchievementsService) GetAchievedUserAchievements(userLogin string) (*models.AchievementsResponse, error) {
	var achievementsResult []models.AchievedAchievement
	resp, err := s.achievementsClient.GetAchievedUserAchievements(context.Background(), &microservices.GetUserAchievementsRequest{
		UserLogin: userLogin,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAchievedUserAchievements: %w", err)
	}
	for _, a := range resp.GetAchievedAchievements() {
		temp := s.convertAchievement(a)
		achievementsResult = append(achievementsResult, *temp)
	}
	return &models.AchievementsResponse{Data: achievementsResult}, nil
}

func (s *AchievementsService) GetOneAchievement(userLogin, code string) (*models.AchievementResponse, error) {
	resp, err := s.achievementsClient.GetOneAchievement(context.Background(), &microservices.GetOneAchievementRequest{
		UserLogin: userLogin,
		Id:        code,
	})
	if err != nil {
		return nil, fmt.Errorf("GetOneAchievement: %w", err)
	}
	temp := s.convertAchievement(resp.GetAchievedAchievement())
	return &models.AchievementResponse{Data: *temp}, nil
}

func (s *AchievementsService) GetAllAchievementsItems() (*models.AchievementsItemsResponse, error) {
	var achievements []models.Achievements
	resp, err := s.achievementsClient.GetAllAchievementsItems(context.Background(), &microservices.GetAllAchievementsItemsRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetAllAchievementsItems: %w", err)
	}
	for _, a := range resp.GetAchievements() {
		achievements = append(achievements, models.Achievements{
			ID:   a.GetId(),
			Code: a.GetCode(),
		})
	}
	return &models.AchievementsItemsResponse{Data: achievements}, nil
}

func (s *AchievementsService) convertAchievement(achievement *microservices.AchievedAchievement) *models.AchievedAchievement {
	var achievedDate string
	if achievement.GetLevel() > 1 {
		achievedDate = achievement.GetUpdatedAt()
	} else {
		achievedDate = achievement.GetCreatedAt()
	}

	totalUsers := s.userRepo.GetAllUserCount()
	totalAchievedPercentage := (float64(achievement.TotalAchieved) / float64(totalUsers)) * 100
	threshold, _, _ := utils.GetThreshold(achievement.GetAchievement().GetCode(), int(achievement.GetLevel()))
	achievementResult := &models.AchievedAchievement{
		AchievedDate:    achievedDate,
		ID:              achievement.GetId(),
		AchievementID:   achievement.GetAchievement().GetId(),
		AchievementCode: achievement.GetAchievement().GetCode(),
		UserLogin:       achievement.GetUserLogin(),
		Achievement: models.Achievements{
			ID:   achievement.GetAchievement().GetId(),
			Code: achievement.GetAchievement().GetCode(),
		},
		Exp:                     achievement.GetXp(),
		Level:                   achievement.GetLevel(),
		TotalAchieved:           achievement.GetTotalAchieved(),
		TotalAchievedPercentage: totalAchievedPercentage,
		CurrentProgress:         achievement.GetCurrentValueInt(),
		Threshold:               int64(threshold),
	}
	return achievementResult
}
