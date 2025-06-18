package services

import (
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
)

type ItemTypesService struct {
	itemTypesRepo *repositories.ItemTypesRepository
}

func NewItemTypesService(itemTypesRepo *repositories.ItemTypesRepository) *ItemTypesService {
	return &ItemTypesService{
		itemTypesRepo: itemTypesRepo,
	}
}

func (s *ItemTypesService) GetAllTypes() (*models.ItemTypesDataResponse, error) {
	var types []models.ItemTypes

	types, err := s.itemTypesRepo.FindAllTypes()
	if err != nil {
		return nil, err
	}
	return &models.ItemTypesDataResponse{Data: types}, nil
}
