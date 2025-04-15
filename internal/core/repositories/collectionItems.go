package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type CiRepository struct {
	db *gorm.DB
}

func NewCiRepository(db *gorm.DB) *CiRepository {
	return &CiRepository{db: db}
}

func (r *CiRepository) CreateCI(ci *models.CollectionItems) (*models.CollectionItems, error) {
	err := r.db.Create(ci)
	if err.Error != nil {
		return nil, err.Error
	}
	return ci, nil
}

func (r *CiRepository) FindAllCI() ([]models.CollectionItems, error) {
	var items []models.CollectionItems
	err := r.db.Find(&items).Error
	return items, err
}

func (r *CiRepository) GetCIByID(id string) (*models.CollectionItems, error) {
	var item models.CollectionItems
	err := r.db.
		Preload("Collections").
		Preload("Collections.User").
		Preload("Platform").
		Preload("Entities").
		Where("id = ? AND deleted = ?", id, false).
		First(&item).Error

	return &item, err
}

func (r *CiRepository) GetCountByUserLogin(login string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.CollectionItems{}).
		Where("LOWER(owner) = LOWER(?) AND deleted = ?", login, false).
		Count(&count).Error

	return count, err
}

func (r *CiRepository) UpdateCI(existsItem *models.CollectionItems, updated *models.CollectionItems) (*models.CollectionItems, error) {
	if updated.Platform != nil {
		err := r.db.Model(existsItem).Association("Platform").Replace(updated.Platform)
		if err != nil {
			return nil, err
		}
	}

	if updated.Collections != nil {
		err := r.db.Model(existsItem).Association("Collections").Replace(updated.Collections)
		if err != nil {
			return nil, err
		}
	}

	if updated.Entities != nil {
		err := r.db.Model(existsItem).Association("Entities").Replace(updated.Entities)
		if err != nil {
			return nil, err
		}
	}

	if err := r.db.Model(existsItem).Updates(updated).Error; err != nil {
		return nil, err
	}

	var result models.CollectionItems
	err := r.db.Preload("Platform").Preload("Collections").Preload("Entities").First(&result, existsItem.Id).Error
	return &result, err

}

func (r *CiRepository) Sum(collectionID uuid.UUID) (float64, error) {
	var sum float64
	err := r.db.
		Model(&models.CollectionItems{}).
		Select("SUM(purchase_price)").
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ? AND deleted = ?", collectionID, false).
		Scan(&sum).Error

	return sum, err
}

func (r *CiRepository) SumShippingCost(collectionID uuid.UUID) (float64, error) {
	var sum struct {
		Sum float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("SUM(collection_items.shipping_cost) as sum").
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ? AND collection_items.deleted = ?",
			collectionID, false).
		Scan(&sum).Error

	return sum.Sum, err
}

func (r *CiRepository) SumByUserLogin(userLogin string) (float64, error) {
	var result struct {
		Sum float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("SUM(purchase_price) as sum").
		Where("owner = ?", userLogin).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.Sum, nil
}

func (r *CiRepository) DeleteCI(id string) error {
	err := r.db.Delete(&models.CollectionItems{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CiRepository) GetCountCI(collectionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.CollectionItems{}).
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ? AND collection_items.deleted = ?", collectionID, false).
		Count(&count).Error

	return count, err
}
