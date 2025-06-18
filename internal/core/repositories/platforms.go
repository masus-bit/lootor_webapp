package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type PlatformsRepository struct {
	db *gorm.DB
}

func NewPlatformsRepository(db *gorm.DB) *PlatformsRepository {
	return &PlatformsRepository{db: db}
}

func (r *PlatformsRepository) FindAllPlatforms() ([]models.Platforms, error) {
	var platforms []models.Platforms
	err := r.db.Find(&platforms).Error
	return platforms, err
}

func (r *PlatformsRepository) GetPlatformById(id string) (*models.Platforms, error) {
	var platform models.Platforms
	err := r.db.Where("id = ?", id).First(&platform).Error
	return &platform, err
}
