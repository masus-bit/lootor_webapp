package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type ItemTypesRepository struct {
	db *gorm.DB
}

func NewItemTypesRepository(db *gorm.DB) *ItemTypesRepository {
	return &ItemTypesRepository{db: db}
}

func (r *ItemTypesRepository) FindAllTypes() ([]models.ItemTypes, error) {
	var types []models.ItemTypes
	err := r.db.Find(&types).Error
	return types, err
}

func (r *ItemTypesRepository) GetTypeById(id string) (*models.ItemTypes, error) {
	var itemType models.ItemTypes
	err := r.db.Where("id = ?", id).First(&itemType).Error
	return &itemType, err
}
