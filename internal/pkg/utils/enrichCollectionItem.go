package utils

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/infrastructure/tagsclient"
	"slices"
)

type EnrichedCI interface {
	Update(id string, dto *models.CollectionItemsRequestUpdate) (*models.CollectionItems, error)
	GetById(id string, authUser string) (*models.CollectionItems, error)
}

type EnrichingCIService struct {
	service    EnrichedCI
	tagsClient *tagsclient.GRPCTagsClient
}

func NewEnrichedCIService(service EnrichedCI, tagsClient *tagsclient.GRPCTagsClient) *EnrichingCIService {
	return &EnrichingCIService{service: service, tagsClient: tagsClient}
}

func (s *EnrichingCIService) enrichItem(item *models.CollectionItems, authUser string) (*models.CollectionItemsResponse, error) {
	var response models.CollectionItemsResponse
	protoTags, err := s.tagsClient.GetTagsByEntityId(context.Background(), &microservices.GetTagsByEntityIdRequest{EntityId: item.ID.String()})
	if err != nil {
		return nil, err
	}
	if err = mapstructure.Decode(item, &response); err != nil {
		return nil, err
	}
	var resultTags []models.ShortTags
	if protoTags != nil {
		resultTags = NormalizeTagsShort(protoTags.GetTags())
	}

	response.Tags = resultTags

	response.Collection = item.Collections[0].ID
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

		protoTags, err := s.tagsClient.GetTagsByEntityId(context.Background(), &microservices.GetTagsByEntityIdRequest{EntityId: item.ID.String()})
		if err != nil {
			return nil, err
		}

		var resultTags []models.ShortTags = NormalizeTagsShort(protoTags.GetTags())
		temp.Tags = resultTags

		temp.Collection = item.Collections[0].ID
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
	if err != nil {
		return nil, err
	}
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
