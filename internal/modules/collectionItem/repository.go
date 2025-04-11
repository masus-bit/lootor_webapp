package collectionItem

import (
	"gorm.io/gorm"
	"lootor/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ci *models.CollectionItems) error {
	return r.db.Create(ci).Error
}

func (r *Repository) FindAll() ([]models.CollectionItems, error) {
	var items []models.CollectionItems
	err := r.db.Find(&items).Error
	return items, err
}

func (r *Repository) GetByID(id string) (*models.CollectionItems, error) {
	var item models.CollectionItems
	err := r.db.
		Preload("Collections").
		Preload("Collections.User").
		Preload("Platforms").
		Preload("Entities").
		Where("id = ? AND deleted = ?", id, false).
		First(&item).Error

	return &item, err
}

func (r *Repository) GetCountByUserLogin(login string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.CollectionItems{}).
		Where("LOWER(owner) = LOWER(?) AND deleted = ?", login, false).
		Count(&count).Error

	return count, err
}

func (r *Repository) Update(existsItem models.CollectionItems, updated models.CollectionItems) (*models.CollectionItems, error) {
	result := r.db.Model(existsItem).Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedItem models.CollectionItems
	err := r.db.First(&updatedItem, existsItem.Id).Error
	return &updatedItem, err
}

func (r *Repository) Sum(collectionID string) (float64, error) {
	var sum float64
	err := r.db.
		Model(&models.CollectionItems{}).
		Select("SUM(purchase_price)").
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collectionItemsId = collection_items.id").
		Where("collections_collection_items_collection_items.collectionsId = ? AND deleted = ?", collectionID, false).
		Scan(&sum).Error

	return sum, err
}

func (r *Repository) SumShippingCost(collectionID string) (float64, error) {
	var sum struct {
		Sum float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("SUM(collection_item.shipping_cost) as sum").
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collectionItemsId = collection_items.id").
		Where("collections_collection_items_collection_items.collectionsId = ? AND collection_items.deleted = ?",
			collectionID, false).
		Scan(&sum).Error

	return sum.Sum, err
}

func (r *Repository) SumByUserLogin(userLogin string) (float64, error) {
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

func (r *Repository) Delete(id string) error {
	err := r.db.Delete(&models.CollectionItems{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetCount(collectionID string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.CollectionItems{}).
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collectionItemsId = collection_items.id").
		Where("collections_collection_items_collection_items.collectionsId = ? AND collection_items.deleted = ?", collectionID, false).
		Count(&count).Error

	return count, err
}
