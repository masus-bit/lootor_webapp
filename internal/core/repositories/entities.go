package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type EntitiesRepository struct {
	db *gorm.DB
}

func NewEntitiesRepository(db *gorm.DB) *EntitiesRepository {
	return &EntitiesRepository{db: db}
}

func (r *EntitiesRepository) CreateEntity(entity *models.Entities) (*models.Entities, error) {
	err := r.db.Create(entity)
	if err != nil {
		return nil, err.Error
	}
	return entity, nil
}

func (r *EntitiesRepository) GetEntityByName(name string) (*models.Entities, error) {
	var entity models.Entities
	err := r.db.Where("name = ?", name).First(&entity).Error
	return &entity, err
}

func (r *EntitiesRepository) GetAllEntities() ([]models.Entities, error) {
	var entities []models.Entities
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *EntitiesRepository) SearchByEntityName(searchTerm string) ([]*models.Entities, error) {
	var entities []*models.Entities

	err := r.db.Where("name ILIKE ?", searchTerm).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
