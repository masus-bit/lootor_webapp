package services

import (
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
)

type PlatformsService struct {
	platformsRepo *repositories.PlatformsRepository
}

func NewPlatformsService(platformsRepo *repositories.PlatformsRepository) *PlatformsService {
	return &PlatformsService{
		platformsRepo: platformsRepo,
	}
}

func (s *PlatformsService) GetAllPlatforms() (*models.PlatformsDataResponse, error) {
	var platforms []models.Platforms

	platforms, err := s.platformsRepo.FindAllPlatforms()
	if err != nil {
		return nil, err
	}
	return &models.PlatformsDataResponse{Data: platforms}, nil
}
