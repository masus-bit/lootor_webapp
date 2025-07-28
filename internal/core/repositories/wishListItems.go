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

	err = r.db.Where("id = ?", item.Id).Preload("User").Preload("CollectionItem").
		Preload("CollectionItem.Platform").
		Preload("CollectionItem.Entities").
		Preload("CollectionItem.Owner").
		Preload("CollectionItem.ItemType").
		Preload("CollectionItem.Collections").
		First(item)
	if err.Error != nil {
		return nil, err.Error
	}
	return item, nil
}

func (r *WLRepository) GetAllByUserLogin(userLogin string) ([]models.WishListItems, int64, error) {
	var items []models.WishListItems
	var totalCount int64

	countQuery := r.db.
		Model(&models.WishListItems{}).
		Where("LOWER(user_login) = LOWER(?)", userLogin).Count(&totalCount).Error

	if countQuery != nil {
		return nil, 0, countQuery
	}

	err := r.db.Where("LOWER(user_login) = LOWER(?)", userLogin).Preload("User").Preload("CollectionItem").
		Preload("CollectionItem.Platform").
		Preload("CollectionItem.Entities").
		Preload("CollectionItem.Owner").
		Preload("CollectionItem.ItemType").
		Preload("CollectionItem.Collections").
		Order("priority ASC").Find(&items)
	if err.Error != nil {
		return nil, 0, err.Error
	}
	return items, totalCount, nil

}

func (r *WLRepository) GetById(id string) (*models.WishListItems, error) {
	var item models.WishListItems
	err := r.db.Where("id = ?", id).Preload("User").Preload("CollectionItem").First(&item).Error

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *WLRepository) UpdatePriority(item *models.WishListItems) ([]models.WishListItems, int64, error) {
	err := r.db.Model(item).Updates(item).Error

	if err != nil {
		return nil, 0, err
	}
	return r.GetAllByUserLogin(item.UserLogin)
}

func (r *WLRepository) DeleteItem(id string) error {
	err := r.db.Delete(&models.WishListItems{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WLRepository) UpdateItem(item *models.WishListItems) (*models.WishListItems, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(item)
	if result.Error != nil {
		return nil, result.Error
	}

	var updated models.WishListItems
	err := r.db.Where("id = ?", item.Id).Preload("User").Preload("CollectionItem").First(&updated).Error

	return &updated, err
}
