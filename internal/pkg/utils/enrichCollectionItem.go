package utils

import (
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/models"
	"slices"
)

type EnrichedCI interface {
	Update(id string, dto *models.CollectionItemsRequestUpdate) (*models.CollectionItems, error)
	GetById(id string, authUser string) (*models.CollectionItems, error)
	GetByEntityAndType(entity string, itemType string, limit string, offset string, search string, orderBy string, order string, authUser string) ([]models.CollectionItems, *models.Entities, int64, error)
}

type EnrichingCIService struct {
	service EnrichedCI
}

func NewEnrichedCIService(service EnrichedCI) *EnrichingCIService {
	return &EnrichingCIService{service: service}
}

func (s *EnrichingCIService) enrichItem(item *models.CollectionItems, authUser string) (*models.CollectionItemsResponse, error) {
	var response models.CollectionItemsResponse
	if err := mapstructure.Decode(item, &response); err != nil {
		return nil, err
	}

	response.Collection = item.Collections[0].Id
	response.Owner = item.Owner
	response.LikesCount = int64(len(item.Likes))
	response.IsOwner = authUser == item.Owner.Login
	response.CollectionTransliteration = item.Collections[0].Transliteration
	response.CollectionName = item.Collections[0].Name

	if authUser != "" {
		response.CanLike = authUser != item.Owner.Login && !slices.Contains(item.Likes, authUser)
	}

	return &response, nil
}

func (s *EnrichingCIService) enrichAnyItems(items []models.CollectionItems, authUser string) ([]models.CollectionItemsResponse, error) {
	var resultCollectionItems []models.CollectionItemsResponse
	for _, item := range items {
		var temp models.CollectionItemsResponse
		err := mapstructure.Decode(item, &temp)
		if err != nil {
			return nil, err
		}
		temp.Collection = item.Collections[0].Id
		temp.Owner = item.Owner
		temp.LikesCount = int64(len(item.Likes))
		temp.CanLike = true
		temp.IsOwner = authUser == item.Owner.Login
		temp.CollectionTransliteration = item.Collections[0].Transliteration
		temp.CollectionName = item.Collections[0].Name

		if authUser != "" {
			if authUser == item.Owner.Login {
				temp.CanLike = true
			} else {
				temp.CanLike = !slices.Contains(item.Likes, authUser)
			}

		}
		resultCollectionItems = append(resultCollectionItems, temp)
	}
	return resultCollectionItems, nil
}

func (s *EnrichingCIService) GetById(id string, authUser string) (*models.CollectionItemsDataResponse, error) {
	item, err := s.service.GetById(id, authUser)
	if err != nil {
		return nil, err
	}
	enriched, err := s.enrichItem(item, authUser)
	return &models.CollectionItemsDataResponse{Data: *enriched}, err
}

func (s *EnrichingCIService) Update(id string, dto *models.CollectionItemsRequestUpdate, authUser string) (*models.CollectionItemsDataResponse, error) {
	item, err := s.service.Update(id, dto)
	if err != nil {
		return nil, err
	}
	enriched, err := s.enrichItem(item, authUser)
	return &models.CollectionItemsDataResponse{Data: *enriched}, err
}

func (s *EnrichingCIService) GetByEntityAndType(entity string, itemType string, limit string, offset string, search string, orderBy string, order string, authUser string) (*models.CollectionItemsDataResponseWithCount, error) {
	items, entityModel, total, err := s.service.GetByEntityAndType(entity, itemType, limit, offset, search, orderBy, order, authUser)
	if err != nil {
		return nil, err
	}
	enriched, err := s.enrichAnyItems(items, authUser)
	return &models.CollectionItemsDataResponseWithCount{Data: enriched, Total: total, Entity: *entityModel}, err
}
