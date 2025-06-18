package repositories

import (
	"context"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
)

type TagsRepository struct {
	db *gorm.DB
	es *elasticsearch.ElasticService
}

func NewTagsRepository(db *gorm.DB, es *elasticsearch.ElasticService) *TagsRepository {
	return &TagsRepository{db: db, es: es}
}

func (r *TagsRepository) AddTag(tag *models.Tags) (*models.Tags, error) {
	result := r.db.Create(tag)
	if result.Error != nil {
		return nil, result.Error
	}

	doc := map[string]interface{}{
		"id":   tag.Id.String(),
		"name": tag.Name,
	}

	if err := r.es.IndexDocument(context.Background(), "tags", doc); err != nil {
		log.Printf("Failed to index tag: %v", err)
	}

	var createdTag models.Tags
	if err := r.db.First(&createdTag, tag.Id).Error; err != nil {
		return nil, err
	}

	return &createdTag, nil
}

func (r *TagsRepository) SearchTagsByName(searchTerm string) ([]*models.Tags, error) {
	var tags []*models.Tags

	err := r.db.Where("name ILIKE ?", searchTerm).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *TagsRepository) GetTagByName(name string) (*models.Tags, error) {
	var tag models.Tags
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagsRepository) FindAllTags() ([]*models.Tags, error) {
	var tags []*models.Tags
	err := r.db.Find(&tags).Error
	return tags, err
}
