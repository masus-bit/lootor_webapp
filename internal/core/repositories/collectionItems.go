package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
)

type CiRepository struct {
	db *gorm.DB
	es *elasticsearch.ElasticService
}

func NewCiRepository(db *gorm.DB, es *elasticsearch.ElasticService) *CiRepository {
	return &CiRepository{db: db, es: es}
}

func (r *CiRepository) CreateCI(ci *models.CollectionItems) (*models.CollectionItems, error) {
	err := r.db.Create(ci)
	if err.Error != nil {
		return nil, err.Error
	}

	doc := map[string]interface{}{
		"id":   ci.Id.String(),
		"name": ci.Name,
	}

	if err := r.es.IndexDocument(context.Background(), "collection_items", doc); err != nil {
		log.Printf("Failed to index collection item: %v", err)
		// Не возвращаем ошибку, чтобы не ломать основной flow
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
		Preload("Owner").
		Preload("Platform").
		Preload("Entities").
		Preload("ItemType").
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

	if updated.ItemType != nil {
		err := r.db.Model(existsItem).Association("ItemType").Replace(updated.ItemType)
		if err != nil {
			return nil, err
		}
	}

	if err := r.db.Model(existsItem).Updates(updated).Error; err != nil {
		return nil, err
	}

	doc := map[string]interface{}{
		"id":   updated.Id.String(),
		"name": updated.Name,
	}

	if err := r.es.IndexDocument(context.Background(), "collection_items", doc); err != nil {
		log.Printf("Failed to index collection item: %v", err)
		// Не возвращаем ошибку, чтобы не ломать основной flow
	}

	var result models.CollectionItems
	err := r.db.Preload("Platform").Preload("Collections").Preload("Entities").Preload("ItemType").First(&result, existsItem.Id).Error
	return &result, err

}

func (r *CiRepository) Sum(collectionID uuid.UUID) (float64, error) {
	var result struct {
		Sum *float64 `gorm:"column:sum"` // Используем указатель для обработки NULL
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(purchase_price), 0) as sum"). // Заменяем NULL на 0
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ? AND deleted = ?", collectionID, false).
		Scan(&result).Error

	if result.Sum == nil {
		return 0, err // На случай, если все же получим nil
	}
	return *result.Sum, err
}
func (r *CiRepository) SumShippingCost(collectionID uuid.UUID) (float64, error) {
	var result struct {
		Sum *float64 `gorm:"column:sum"` // Используем указатель для обработки NULL
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(shipping_cost), 0) as sum"). // Заменяем NULL на 0
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ? AND collection_items.deleted = ?",
			collectionID, false).
		Scan(&result).Error

	if result.Sum == nil {
		return 0, err // На случай, если все же получим nil
	}
	return *result.Sum, err
}

func (r *CiRepository) SumByUserLogin(userLogin string) (float64, error) {
	var result struct {
		Sum *float64 `gorm:"column:sum"` // Используем указатель для обработки NULL
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(purchase_price), 0) as sum"). // Заменяем NULL на 0
		Where("user_login = ?", userLogin).                // Используем user_login вместо owner (согласно вашей модели)
		Scan(&result).Error

	if result.Sum == nil {
		return 0, err
	}
	return *result.Sum, err
}

func (r *CiRepository) DeleteCI(id string) error {
	err := r.db.Delete(&models.CollectionItems{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	if err := r.es.DeleteDocument(context.Background(), "collection_items", id); err != nil {
		log.Printf("Failed to delete collection item from index: %v", err)
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
