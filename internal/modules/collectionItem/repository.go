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

func (r *Repository) Create(ci *models.CollectionItem) error {
	return r.db.Create(ci).Error
}

func (r *Repository) FindAll() ([]models.CollectionItem, error) {
	var items []models.CollectionItem
	err := r.db.Find(&items).Error
	return items, err
}

func (r *Repository) GetByID(id string) (*models.CollectionItem, error) {
	var item models.CollectionItem
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
		Model(&models.CollectionItem{}).
		Where("LOWER(owner) = LOWER(?) AND deleted = ?", login, false).
		Count(&count).Error

	return count, err
}

func (r *Repository) Update(item *models.CollectionItem) error {
	return r.db.Model(&models.CollectionItem{}).
		Where("id = ?", item.ID).
		Updates(item).Error
}

func (r *Repository) Sum(collectionID string) (float64, error) {
	var sum float64
	err := r.db.
		Model(&models.CollectionItem{}).
		Select("SUM(purchase_price)").
		Joins("INNER JOIN collection_collection_items_collection_item ON collection_collection_items_collection_item.collectionItemId = collection_item.id").
		Where("collection_collection_items_collection_item.collectionId = ? AND deleted = ?", collectionID, false).
		Scan(&sum).Error

	return sum, err
}

func (r *Repository) SumShippingCost(collectionID string) (float64, error) {
	var sum struct {
		Sum float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItem{}).
		Select("SUM(collection_item.shipping_cost) as sum").
		Joins("INNER JOIN collection_collection_items_collection_item ON collection_collection_items_collection_item.collectionItemId = collection_item.id").
		Where("collection_collection_items_collection_item.collectionId = ? AND collection_item.deleted = ?",
			collectionID, false).
		Scan(&sum).Error

	return sum.Sum, err
}

func (r *Repository) SumByUserLogin(userLogin string) (float64, error) {
	var result struct {
		Sum float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItem{}).
		Select("SUM(purchase_price) as sum").
		Where("owner = ?", userLogin).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.Sum, nil
}

func (r *Repository) Delete(id string) error {
	err := r.db.Delete(&models.CollectionItem{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetCount(collectionID string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.CollectionItem{}).
		Joins("INNER JOIN collection_collection_items_collection_item ON collection_collection_items_collection_item.collectionItemId = collection_item.id").
		Where("collection_collection_items_collection_item.collectionId = ? AND collection_item.deleted = ?", collectionID, false).
		Count(&count).Error

	return count, err
}
