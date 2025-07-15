package repositories

import (
	"context"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
	"strconv"
)

type EntitiesRepository struct {
	db *gorm.DB
	es *elasticsearch.ElasticService
}

func NewEntitiesRepository(db *gorm.DB, es *elasticsearch.ElasticService) *EntitiesRepository {
	return &EntitiesRepository{db: db, es: es}
}

func (r *EntitiesRepository) CreateEntity(entity *models.Entities) (*models.Entities, error) {
	err := r.db.Create(entity)
	if err.Error != nil {
		return nil, err.Error
	}

	doc := map[string]interface{}{
		"id":              entity.Id.String(),
		"name":            entity.Name,
		"transliteration": entity.Transliteration,
		// другие поля
	}

	if err := r.es.IndexDocument(context.Background(), "entities", doc); err != nil {
		log.Printf("Failed to index entity: %v", err)
	}

	return entity, nil
}

func (r *EntitiesRepository) GetEntityByName(name string) (*models.Entities, error) {
	var entity models.Entities
	err := r.db.Where("name = ?", name).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, err
}

func (r *EntitiesRepository) GetEntityByTranslit(translit string) (*models.Entities, error) {
	var entity models.Entities
	err := r.db.Where("transliteration = ?", translit).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, err
}

func (r *EntitiesRepository) GetAllEntities() ([]models.Entities, error) {
	var entities []models.Entities
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *EntitiesRepository) GetAll(search, limit, offset string) ([]models.Entities, error) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)
	var entities []models.Entities
	query := r.db.
		Order("entities.name ASC").
		Limit(intLimit).
		Offset(intOffset)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&entities).Error

	return entities, err
}

func (r *EntitiesRepository) SearchByEntityName(searchTerm string) ([]models.Entities, error) {
	var entities []models.Entities

	err := r.db.Where("name ILIKE ?", "%"+searchTerm+"%").Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}
