package services

import (
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
)

type EntitiesService struct {
	eRepo *repositories.EntitiesRepository
}

func NewEntitiesService(eRepo *repositories.EntitiesRepository) *EntitiesService {
	return &EntitiesService{
		eRepo: eRepo,
	}
}

func (s *EntitiesService) CreateEntities(dto *models.EntitiesCreateRequest) (*models.EntitiesDataResponse, error) {
	var resultEntities []models.Entities
	for _, entity := range dto.Names {
		dbEntity, _ := s.eRepo.CreateEntity(&models.Entities{Name: entity})
		resultEntities = append(resultEntities, *dbEntity)
	}
	return &models.EntitiesDataResponse{Data: resultEntities}, nil
}

func (s *EntitiesService) SearchEntities(name string) (*models.EntitiesDataResponse, error) {
	var entities []models.Entities
	entities, _ = s.eRepo.SearchByEntityName(name)

	return &models.EntitiesDataResponse{Data: entities}, nil
}
