package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type TagsRepository struct {
	db *gorm.DB
}

func NewTagsRepository(db *gorm.DB) *TagsRepository {
	return &TagsRepository{db: db}
}

func (r *TagsRepository) AddTag(tag *models.Tags) (*models.Tags, error) {
	err := r.db.Create(tag)
	if err != nil {
		return nil, err.Error
	}
	return tag, nil
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
	return &tag, err
}
