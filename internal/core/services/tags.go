package services

import (
	"fmt"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
)

type TagsService struct {
	tagRepo *repositories.TagsRepository
}

func NewTagsService(tagRepo *repositories.TagsRepository) *TagsService {
	return &TagsService{
		tagRepo: tagRepo,
	}
}

func (s *TagsService) SearchTags(name string) ([]*models.Tags, error) {
	searchTerm := fmt.Sprintf("%%%s%%", name)

	tags, err := s.tagRepo.SearchTagsByName(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}

	return tags, nil
}
