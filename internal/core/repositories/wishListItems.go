package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type WLRepository struct {
	db *gorm.DB
}

func NewWLRepository(db *gorm.DB) *WLRepository {
	return &WLRepository{db: db}
}

func (r *WLRepository) AddItem(item *models.WishListItems) (*models.WishListItems, error) {
	err := r.db.Create(item)
	if err.Error != nil {
		return nil, err.Error
	}
	return item, nil
}

func (r *WLRepository) GetAllByUserLogin(userLogin string) ([]models.WishListItems, error) {
	var items []models.WishListItems

	err := r.db.Where("LOWER(user_login) = LOWER(?)", userLogin).Preload("User").Preload("CollectionItem").Order("priority DESC").Find(&items)
	if err.Error != nil {
		return nil, err.Error
	}
	return items, nil

}

func (r *WLRepository) GetById(id string) (*models.WishListItems, error) {
	var item models.WishListItems
	err := r.db.Where("id = ?", id).First(&item)

	if err != nil {
		return nil, err.Error
	}
	return &item, nil
}
