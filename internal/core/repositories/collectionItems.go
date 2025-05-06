package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/utils"
	"strconv"
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
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&item).Error

	return &item, err
}

func (r *CiRepository) GetCountByUserLogin(login string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.CollectionItems{}).
		Where("LOWER(user_login) = LOWER(?)", login).
		Where("deleted_at IS NULL").
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
		Sum *float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(purchase_price), 0) as sum").
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ?", collectionID).
		Where("deleted_at IS NULL").
		Scan(&result).Error

	if result.Sum == nil {
		return 0, err
	}
	return *result.Sum, err
}
func (r *CiRepository) SumShippingCost(collectionID uuid.UUID) (float64, error) {
	var result struct {
		Sum *float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(shipping_cost), 0) as sum").
		Joins("INNER JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id = ?",
			collectionID).
		Where("collection_items.deleted_at IS NULL").
		Scan(&result).Error

	if result.Sum == nil {
		return 0, err
	}
	return *result.Sum, err
}

func (r *CiRepository) SumByUserLogin(userLogin string) (float64, error) {
	var result struct {
		Sum *float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(purchase_price), 0) as sum").
		Where("user_login = ?", userLogin).
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
		Where("collections_collection_items_collection_items.collections_id = ?", collectionID).
		Where("collection_items.deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

func (r *CiRepository) GetCountCIByIDs(collectionIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	var results []struct {
		CollectionsId uuid.UUID
		Count         int64
	}

	err := r.db.Table("collections_collection_items_collection_items").
		Select("collections_id as collections_id, COUNT(*) as count").
		Joins("JOIN collections ON collections.id = collections_collection_items_collection_items.collections_id").
		Joins("JOIN collection_items ON collection_items.id = collections_collection_items_collection_items.collection_items_id").
		Where("collections_id IN (?)", collectionIDs).
		Where("collections.deleted_at IS NULL").
		Where("collection_items.deleted_at IS NULL").
		Group("collections_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int64)
	for _, result := range results {
		counts[result.CollectionsId] = result.Count
	}

	for _, id := range collectionIDs {
		if _, exists := counts[id]; !exists {
			counts[id] = 0
		}
	}

	return counts, nil
}

func (r *CiRepository) GetSumsByCollectionIDs(collectionIDs []uuid.UUID) (map[uuid.UUID]float64, error) {
	var results []struct {
		CollectionId uuid.UUID
		Sum          float64
	}
	err := r.db.Model(&models.CollectionItems{}).
		Select("collections_collection_items_collection_items.collections_id as collection_id, SUM(collection_items.purchase_price) as sum").
		Joins("JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id IN (?)", collectionIDs).
		Where("collection_items.deleted_at IS NULL").
		Group("collections_collection_items_collection_items.collections_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	sums := make(map[uuid.UUID]float64)
	for _, result := range results {
		sums[result.CollectionId] = result.Sum
	}

	for _, id := range collectionIDs {
		if _, exists := sums[id]; !exists {
			sums[id] = 0
		}
	}

	return sums, nil
}

func (r *CiRepository) GetShippingCostsByCollectionIDs(collectionIDs []uuid.UUID) (map[uuid.UUID]float64, error) {
	var results []struct {
		CollectionId uuid.UUID
		Sum          float64
	}

	err := r.db.Model(&models.CollectionItems{}).
		Select("collections_collection_items_collection_items.collections_id as collection_id, SUM(collection_items.shipping_cost) as sum").
		Joins("JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id").
		Where("collections_collection_items_collection_items.collections_id IN (?)", collectionIDs).
		Where("collection_items.deleted_at IS NULL").
		Group("collections_collection_items_collection_items.collections_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	sums := make(map[uuid.UUID]float64)
	for _, result := range results {
		sums[result.CollectionId] = result.Sum
	}

	for _, id := range collectionIDs {
		if _, exists := sums[id]; !exists {
			sums[id] = 0
		}
	}

	return sums, nil
}

func (r *CiRepository) GetCollectionItemsByEntity(entity string, limit string) (*struct {
	VideoGames         []models.CollectionItems `json:"videoGames"`
	BoardGames         []models.CollectionItems `json:"boardGames"`
	Comics             []models.CollectionItems `json:"comics"`
	GamingHardware     []models.CollectionItems `json:"gamingHardware"`
	CollectibleFigures []models.CollectionItems `json:"collectibleFigures"`
	Books              []models.CollectionItems `json:"books"`
	Vinyl              []models.CollectionItems `json:"vinyl"`
}, int64, error) {
	var GroupedCollectionItems struct {
		VideoGames         []models.CollectionItems `json:"videoGames"`
		BoardGames         []models.CollectionItems `json:"boardGames"`
		Comics             []models.CollectionItems `json:"comics"`
		GamingHardware     []models.CollectionItems `json:"gamingHardware"`
		CollectibleFigures []models.CollectionItems `json:"collectibleFigures"`
		Books              []models.CollectionItems `json:"books"`
		Vinyl              []models.CollectionItems `json:"vinyl"`
	}
	var totalCount int64
	var itemTypesArr = []string{"Video games", "Board games", "Comics", "Gaming hardware", "Collectible figures", "Books", "Vinyl"}

	intLimit, _ := strconv.Atoi(limit)

	err := r.db.
		Model(&models.CollectionItems{}).
		Joins("JOIN entities_collection_item_collection_items ON entities_collection_item_collection_items.collection_items_id = collection_items.id").
		Joins("JOIN entities ON entities.id = entities_collection_item_collection_items.entities_id").
		Where("LOWER(entities.transliteration) = LOWER(?)", entity).
		Where("collection_items.deleted_at IS NULL").
		Count(&totalCount).Error

	if err != nil {
		return nil, 0, err
	}

	for _, itemType := range itemTypesArr {
		query := r.db.
			Joins("JOIN entities_collection_item_collection_items ON entities_collection_item_collection_items.collection_items_id = collection_items.id").
			Joins("JOIN entities ON entities.id = entities_collection_item_collection_items.entities_id").
			Joins("JOIN item_types ON item_types.id = collection_items.item_type_id").
			Where("LOWER(entities.transliteration) = LOWER(?)", entity).
			Where("item_types.name = ?", itemType).
			Preload("Collections").
			Preload("Collections.User").
			Preload("Owner").
			Preload("Platform").
			Preload("Entities").
			Preload("ItemType").
			Where("collection_items.deleted_at IS NULL").
			Order("collection_items.created_at DESC").
			Limit(intLimit)

		switch itemType {
		case "Video games":
			query.Find(&GroupedCollectionItems.VideoGames)
		case "Board games":
			query.Find(&GroupedCollectionItems.BoardGames)
		case "Comics":
			query.Find(&GroupedCollectionItems.Comics)
		case "Gaming hardware":
			query.Find(&GroupedCollectionItems.GamingHardware)
		case "Collectible figures":
			query.Find(&GroupedCollectionItems.CollectibleFigures)
		case "Books":
			query.Find(&GroupedCollectionItems.Books)
		case "Vinyl":
			query.Find(&GroupedCollectionItems.Vinyl)
		}
	}

	return &GroupedCollectionItems, totalCount, nil

}

func (r *CiRepository) GetByEntityAndType(entity string, itemType string, limit string, offset string, search string, orderBy string, order string) ([]models.CollectionItems, int64, error) {
	var collectionItems []models.CollectionItems
	var totalCount int64

	itemTypeNew := utils.GetReformatedItemType(itemType)

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	countQuery := r.db.
		Model(&models.CollectionItems{}).
		Joins("JOIN entities_collection_item_collection_items ON entities_collection_item_collection_items.collection_items_id = collection_items.id").
		Joins("JOIN entities ON entities.id = entities_collection_item_collection_items.entities_id").
		Joins("JOIN item_types ON item_types.id = collection_items.item_type_id").
		Where("LOWER(entities.transliteration) = LOWER(?)", entity).
		Where("item_types.name = ?", itemTypeNew).
		Where("collection_items.deleted_at IS NULL")
	if search != "" {
		countQuery = countQuery.Where("collection_items.name ILIKE ?", "%"+search+"%")
	}

	err := countQuery.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	query := r.db.
		Joins("JOIN entities_collection_item_collection_items ON entities_collection_item_collection_items.collection_items_id = collection_items.id").
		Joins("JOIN entities ON entities.id = entities_collection_item_collection_items.entities_id").
		Joins("JOIN item_types ON item_types.id = collection_items.item_type_id").
		Where("LOWER(entities.transliteration) = LOWER(?)", entity).
		Where("item_types.name = ?", itemTypeNew).
		Preload("Collections").
		Preload("Collections.User").
		Preload("Owner").
		Preload("Platform").
		Preload("Entities").
		Preload("ItemType").
		Where("collection_items.deleted_at IS NULL").
		Order("collection_items.created_at DESC").
		Limit(intLimit).
		Offset(intOffset)

	if search != "" {
		query = query.Where("collection_items.name ILIKE ?", "%"+search+"%")
	}

	ciOrder := utils.GetCIOrderString(orderBy, order)
	if ciOrder.Joins != "" {
		query.Joins(ciOrder.Joins).Order(ciOrder.Order)
	} else {
		query.Order(ciOrder.Order)
	}

	err = query.Find(&collectionItems).Error
	if err != nil {
		return nil, 0, err
	}
	return collectionItems, totalCount, nil
}
