package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/utils"
	"reflect"
	"slices"
	"strconv"
	"time"
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
		"id":          ci.Id.String(),
		"name":        ci.Name,
		"description": ci.Description,
		"images":      ci.Images,
	}

	if err := r.es.IndexDocument(context.Background(), "collection_items", doc); err != nil {
		log.Printf("Failed to index collection item: %v", err)
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
		"id":          updated.Id.String(),
		"name":        updated.Name,
		"description": updated.Description,
		"images":      updated.Images,
	}

	if err := r.es.IndexDocument(context.Background(), "collection_items", doc); err != nil {
		log.Printf("Failed to index collection item: %v", err)
		// Не возвращаем ошибку, чтобы не ломать основной flow
	}

	var result models.CollectionItems
	err := r.db.Preload("Platform").Preload("Collections").Preload("ItemType").First(&result, existsItem.Id).Error
	return &result, err

}

func (r *CiRepository) UpdateCIFull(existsItem *models.CollectionItems) (*models.CollectionItems, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result models.CollectionItems
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(existsItem).Select("*").Updates(existsItem).Error; err != nil {
			return fmt.Errorf("failed to update base fields: %w", err)
		}

		associations := []struct {
			name  string
			value interface{}
		}{
			{"Platform", existsItem.Platform},
			{"Collections", existsItem.Collections},
			{"ItemType", existsItem.ItemType},
		}

		for _, assoc := range associations {
			if assoc.value != nil {
				if reflect.ValueOf(assoc.value).Kind() == reflect.Slice {
					if reflect.ValueOf(assoc.value).Len() == 0 {
						if err := tx.Model(existsItem).Association(assoc.name).Clear(); err != nil {
							return fmt.Errorf("failed to clear association %s: %w", assoc.name, err)
						}
						continue
					}
				}
				if err := tx.Model(existsItem).Association(assoc.name).Replace(assoc.value); err != nil {
					return fmt.Errorf("failed to update association %s: %w", assoc.name, err)
				}
			}
		}

		return tx.Select("*").
			Preload("Platform").
			Preload("Collections").
			Preload("ItemType").
			Preload("Owner").
			First(&result, "id = ?", existsItem.Id).
			Error
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	go func() {
		doc := map[string]interface{}{
			"id":          result.Id.String(),
			"name":        result.Name,
			"description": result.Description,
			"images":      result.Images,
		}
		if err := r.es.IndexDocument(context.Background(), "collection_items", doc); err != nil {
			log.Printf("Failed to index collection item: %v", err)
		}
	}()

	return &result, nil
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
		Where("collection_items.deleted_at IS NULL").
		Scan(&result).Error

	if result.Sum == nil {
		return 0, err
	}
	return *result.Sum, err
}

func (r *CiRepository) ShippingSumByUserLogin(userLogin string) (float64, error) {
	var result struct {
		Sum *float64 `gorm:"column:sum"`
	}

	err := r.db.
		Model(&models.CollectionItems{}).
		Select("COALESCE(SUM(shipping_cost), 0) as sum").
		Where("user_login = ?", userLogin).
		Where("collection_items.deleted_at IS NULL").
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

func (r *CiRepository) GetAll(limit, offset, search string) ([]models.CollectionItems, error) {
	var collectionItems []models.CollectionItems

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	query := r.db.
		Preload("Collections").
		Preload("Collections.User").
		Preload("Owner").
		Preload("Platform").
		Preload("ItemType").
		Where("collection_items.deleted_at IS NULL").
		Limit(intLimit).
		Offset(intOffset)

	if search != "" {
		query = query.Where("collection_items.name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&collectionItems).Error
	if err != nil {
		return nil, err
	}

	return collectionItems, nil
}

func (r *CiRepository) IncrementCommentsCount(id string, amount int) error {
	return r.db.Model(&models.CollectionItems{}).
		Where("id = ?", id).
		Update("comments_count", gorm.Expr("COALESCE(comments_count, 0) + ?", amount)).Error
}

func (r *CiRepository) GetCollectionItemsByIdsMap(ids []string, authUserLogin string) (map[string]models.CollectionItemsResponse, error) {
	var collectionItems []models.CollectionItems
	err := r.db.Where("id IN (?)", ids).
		Preload("Collections").
		Preload("Collections.User").
		Preload("Owner").
		Preload("Platform").
		Preload("ItemType").
		Find(&collectionItems).Error
	if err != nil {
		return nil, err
	}
	collectionItemsMap := make(map[string]models.CollectionItemsResponse, len(collectionItems))
	for _, collectionItem := range collectionItems {
		var temp models.CollectionItemsResponse
		err = mapstructure.Decode(collectionItem, &temp)
		if err != nil {
			return nil, err
		}

		if len(collectionItem.Collections) > 0 {
			temp.Collection = collectionItem.Collections[0].Id
			temp.CollectionTransliteration = collectionItem.Collections[0].Transliteration
			temp.CollectionName = collectionItem.Collections[0].Name
		} else {
			temp.Collection = uuid.Nil
			temp.CollectionTransliteration = ""
		}

		if collectionItem.Likes != nil {
			temp.LikesCount = int64(len(collectionItem.Likes))
		} else {
			temp.LikesCount = 0
		}

		temp.Owner = collectionItem.Owner
		temp.CanLike = false

		if authUserLogin != "" {
			if authUserLogin == collectionItem.Owner.Login {
				temp.CanLike = false
			} else {
				temp.CanLike = !slices.Contains(collectionItem.Likes, authUserLogin)
			}

		}
		collectionItemsMap[collectionItem.Id.String()] = temp
	}
	return collectionItemsMap, nil
}

func (r *CiRepository) GetCollectionItemsByIDs(ids []string, authUserLogin, filter string, tagsMap map[string]*microservices.GetShortsResponse) ([]models.CollectionItemsResponse, error) {
	var collectionItems []models.CollectionItems
	query := r.db.Where("id IN (?)", ids).
		Preload("Collections").
		Preload("Collections.User").
		Preload("Owner").
		Preload("Platform").
		Preload("ItemType")

	if filter != "" {
		query = query.Where("collection_items.item_type_id = ?", filter)
	}
	err := query.Find(&collectionItems).Error
	if err != nil {
		return nil, err
	}
	var result []models.CollectionItemsResponse
	for _, collectionItem := range collectionItems {
		var temp models.CollectionItemsResponse
		err = mapstructure.Decode(collectionItem, &temp)
		if err != nil {
			return nil, err
		}
		if len(tagsMap) > 0 {
			tagsProto := tagsMap[collectionItem.Id.String()].Tags
			tempTags := utils.NormalizeTagsShort(tagsProto)
			temp.Tags = tempTags
		}

		if len(collectionItem.Collections) > 0 {
			temp.Collection = collectionItem.Collections[0].Id
			temp.CollectionTransliteration = collectionItem.Collections[0].Transliteration
			temp.CollectionName = collectionItem.Collections[0].Name
		} else {
			temp.Collection = uuid.Nil
			temp.CollectionTransliteration = ""
		}

		if collectionItem.Likes != nil {
			temp.LikesCount = int64(len(collectionItem.Likes))
		} else {
			temp.LikesCount = 0
		}

		temp.Owner = collectionItem.Owner
		temp.CanLike = true

		if authUserLogin != "" {
			if authUserLogin == collectionItem.Owner.Login {
				temp.CanLike = true
			} else {
				temp.CanLike = !slices.Contains(collectionItem.Likes, authUserLogin)
			}

		}

		result = append(result, temp)
	}
	return result, nil
}

func (r *CiRepository) GetCICountsByItemType(entityIDs []string) (map[string]int64, error) {
	var results []struct {
		ItemType string
		Count    int64
	}
	err := r.db.Table("collection_items").
		Where("id IN (?)", entityIDs).
		Select("item_type_id as item_type, COUNT(*) as count").
		Group("item_type").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int64, len(results))
	for _, result := range results {
		counts[result.ItemType] = result.Count
	}
	return counts, nil
}
