package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/utils"
	"slices"
	"strconv"
)

type CollectionsRepository struct {
	db *gorm.DB
	es *elasticsearch.ElasticService
}

func NewCollectionsRepository(db *gorm.DB, es *elasticsearch.ElasticService) *CollectionsRepository {
	return &CollectionsRepository{db: db, es: es}
}

func (r *CollectionsRepository) CreateCollection(collection *models.Collections) (*models.Collections, error) {
	err := r.db.Create(&collection)
	if err.Error != nil {
		return nil, err.Error
	}

	if !collection.IsPrivate {
		doc := map[string]interface{}{
			"id":          collection.ID.String(),
			"name":        collection.Name,
			"description": collection.Description,
			"isPrivate":   collection.IsPrivate,
			"bannerUrl":   collection.BannerURL,
		}

		if err := r.es.IndexDocument(context.Background(), "collections", doc); err != nil {
			log.Printf("Failed to index collection: %v", err)
		}
	}

	return collection, nil
}

func (r *CollectionsRepository) FindAllCollections() ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.Find(&collections).Error
	return collections, err
}

func (r *CollectionsRepository) GetCollectionById(
	id string,
	limit string,
	offset string,
	orderByInput string,
	order string,
	search string,
) (*models.Collections, error) {
	var collection models.Collections

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	err := r.db.
		Preload("User").
		Preload(
			"CollectionItems", func(tx *gorm.DB) *gorm.DB {

				if search != "" {
					tx = tx.Where("name ILIKE ?", "%"+search+"%")
				}

				switch orderByInput {
				case "platform":
					if order == "desc" {
						return tx.
							Offset(intOffset).
							Limit(intLimit).
							Joins("LEFT JOIN platforms ON platforms.id = collection_items.platform_id").
							Order("platforms.name DESC")
					}
					return tx.
						Offset(intOffset).
						Limit(intLimit).
						Joins("LEFT JOIN platforms ON platforms.id = collection_items.platform_id").
						Order("platforms.name ASC")
				case "type":
					if order == "desc" {
						return tx.
							Offset(intOffset).
							Limit(intLimit).
							Joins("LEFT JOIN item_types ON item_types.id = collection_items.item_type_id").
							Order("item_types.ru_name DESC")
					}
					return tx.
						Offset(intOffset).
						Limit(intLimit).
						Joins("LEFT JOIN item_types ON item_types.id = collection_items.item_type_id").
						Order("item_types.ru_name ASC")
				case "name":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("name DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("name ASC")

				case "rating":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("rating DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("rating ASC")

				case "purchaseDate":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date ASC")

				case "purchasePrice":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price ASC")

				default:
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date DESC")
				}
			},
		).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.ItemType").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetCollectionByIdWithoutLimits(id string) (*models.Collections, error) {
	var collection models.Collections

	err := r.db.
		Preload("User").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.ItemType").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetByShareString(
	shareString string,
	limit string,
	offset string,
	orderByInput string,
	order string,
	search string,
) (*models.Collections, error) {
	var collection models.Collections

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	err := r.db.
		Preload("User").
		Preload(
			"CollectionItems", func(tx *gorm.DB) *gorm.DB {

				if search != "" {
					tx = tx.Where("name ILIKE ?", "%"+search+"%")
				}

				switch orderByInput {
				case "platform":
					if order == "desc" {
						return tx.
							Offset(intOffset).
							Limit(intLimit).
							Joins("LEFT JOIN platforms ON platforms.id = collection_items.platform_id").
							Order("platforms.name DESC")
					}
					return tx.
						Offset(intOffset).
						Limit(intLimit).
						Joins("LEFT JOIN platforms ON platforms.id = collection_items.platform_id").
						Order("platforms.name ASC")
				case "type":
					if order == "desc" {
						return tx.
							Offset(intOffset).
							Limit(intLimit).
							Joins("LEFT JOIN item_types ON item_types.id = collection_items.item_type_id").
							Order("item_types.ru_name DESC")
					}
					return tx.
						Offset(intOffset).
						Limit(intLimit).
						Joins("LEFT JOIN item_types ON item_types.id = collection_items.item_type_id").
						Order("item_types.ru_name ASC")
				case "name":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("name DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("name ASC")

				case "rating":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("rating DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("rating ASC")

				case "purchaseDate":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date ASC")

				case "purchasePrice":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price ASC")

				default:
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date DESC")
				}
			},
		).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.ItemType").
		Where("share_string = ?", shareString).
		Where("deleted_at IS NULL").
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetByIdWithoutCollectionItems(id string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("User").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetOneByTransliteration(
	login string,
	transliteration string,
	limit string,
	offset string,
	orderByInput string,
	order string,
	search string,
) (*models.Collections, error) {
	var collection models.Collections

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	err := r.db.
		Where("collections.transliteration = ?", transliteration).
		Where("deleted_at IS NULL").
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Preload(
			"CollectionItems", func(tx *gorm.DB) *gorm.DB {

				if search != "" {
					tx = tx.Where("name ILIKE ?", "%"+search+"%")
				}

				switch orderByInput {
				case "platform":
					if order == "desc" {
						return tx.
							Offset(intOffset).
							Limit(intLimit).
							Joins("LEFT JOIN platforms ON platforms.id = collection_items.platform_id").
							Order("platforms.name DESC, collection_items.id DESC")
					}
					return tx.
						Offset(intOffset).
						Limit(intLimit).
						Joins("LEFT JOIN platforms ON platforms.id = collection_items.platform_id").
						Order("platforms.name ASC, collection_items.id ASC")
				case "type":
					if order == "desc" {
						return tx.
							Offset(intOffset).
							Limit(intLimit).
							Joins("LEFT JOIN item_types ON item_types.id = collection_items.item_type_id").
							Order("item_types.ru_name DESC, collection_items.id DESC")
					}
					return tx.
						Offset(intOffset).
						Limit(intLimit).
						Joins("LEFT JOIN item_types ON item_types.id = collection_items.item_type_id").
						Order("item_types.ru_name ASC, collection_items.id ASC")
				case "name":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("name DESC, collection_items.id DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("name ASC, collection_items.id ASC")

				case "rating":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("rating DESC, collection_items.id DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("rating ASC, collection_items.id ASC")

				case "purchaseDate":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date DESC, collection_items.id DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_date ASC, collection_items.id ASC")

				case "purchasePrice":
					if order == "desc" {
						return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price DESC, collection_items.id DESC")
					}
					return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price ASC")

				default:
					return tx.Offset(intOffset).Limit(intLimit).Order("created_at DESC")
				}
			},
		).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Owner").
		Preload("CollectionItems.ItemType").
		First(&collection).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &collection, nil
}

func (r *CollectionsRepository) GetCollectionByUserId(login string) ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.
		Where("LOWER(user_login) = LOWER(?)", login).Preload("User").
		Preload(
			"CollectionItems", func(tx *gorm.DB) *gorm.DB {
				return tx.Order("created DESC")
			},
		).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Owner").
		Preload("CollectionItems.ItemType").
		Where("deleted_at IS NULL").
		Order("created DESC").
		Find(&collections).Error
	return collections, err
}

func (r *CollectionsRepository) GetByUserIdWithoutCollectionItems(login string, search string) (
	[]models.Collections,
	int64,
	error,
) {
	var collections []models.Collections
	var totalCount int64

	countQuery := r.db.
		Model(&models.Collections{}).
		Where("LOWER(user_login) = LOWER(?)", login).Count(&totalCount).Error

	if countQuery != nil {
		return nil, 0, countQuery
	}

	query := r.db.
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Where("deleted_at IS NULL").
		Order("created DESC")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&collections).Error

	return collections, totalCount, err
}

func (r *CollectionsRepository) GetByUserIdWithoutPrivates(login string, search string) (
	[]models.Collections,
	int64,
	error,
) {
	var collections []models.Collections
	var totalCount int64

	countQuery := r.db.
		Model(&models.Collections{}).
		Where("LOWER(user_login) = LOWER(?)", login).Count(&totalCount).Error

	if countQuery != nil {
		return nil, 0, countQuery
	}

	query := r.db.
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Where("collections.is_private = ?", false).
		Where("collections.deleted_at IS NULL").
		Order("collections.created DESC")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&collections).Error

	return collections, totalCount, err
}

func (r *CollectionsRepository) GetAllWithoutPrivates(search, limit, offset, sortBy, order string, showEmpty bool) (
	[]models.Collections,
	int64,
	error,
) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)
	var totalCount int64

	countQuery := r.db.
		Model(&models.Collections{}).
		Where("collections.deleted_at IS NULL").
		Where("collections.is_private = ?", false)

	if search != "" {
		countQuery = countQuery.Where("collections.name ILIKE ?", "%"+search+"%")
	}

	if !showEmpty {
		countQuery = countQuery.Where("EXISTS (SELECT 1 FROM collections_collection_items_collection_items WHERE collections_id = collections.id)")
	}

	err := countQuery.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	var collections []models.Collections
	query := r.db.
		Preload("User").
		Where("collections.is_private = ?", false).
		Where("collections.deleted_at IS NULL").
		Offset(intOffset).
		Limit(intLimit)

	if !showEmpty {
		query = query.Where("EXISTS (SELECT 1 FROM collections_collection_items_collection_items WHERE collections_id = collections.id)")
	}

	switch sortBy {
	case "likesCount":
		if order == "desc" {
			query = query.Order("array_length(collections.likes, 1) DESC NULLS LAST")
		} else {
			query = query.Order("array_length(collections.likes, 1) ASC NULLS FIRST")
		}
	case "name":
		if order == "desc" {
			query = query.Order("collections.name DESC")
		} else {
			query = query.Order("collections.name ASC")
		}
	case "created":
		if order == "desc" {
			query = query.Order("collections.created DESC")
		} else {
			query = query.Order("collections.created ASC")
		}
	case "totalPrice":
		orderClause := "ASC"
		if order == "desc" {
			orderClause = "DESC"
		}
		query = query.
			Joins(
				"LEFT JOIN (SELECT collections_collection_items_collection_items.collections_id as collection_id, " +
					"COALESCE(SUM(collection_items.purchase_price), 0) as total_price " +
					"FROM collection_items " +
					"JOIN collections_collection_items_collection_items ON collections_collection_items_collection_items.collection_items_id = collection_items.id " +
					"WHERE collection_items.deleted_at IS NULL " +
					"GROUP BY collections_collection_items_collection_items.collections_id" +
					") AS price_sum ON price_sum.collection_id = collections.id",
			).
			Order("price_sum.total_price " + orderClause)
	case "collectionItemsCount":
		orderClause := "ASC"
		if order == "desc" {
			orderClause = "DESC"
		}
		query = query.
			Joins(
				"LEFT JOIN (" +
					"SELECT collections_id, COUNT(*) as items_count " +
					"FROM collections_collection_items_collection_items " +
					"JOIN collection_items ON collection_items.id = collections_collection_items_collection_items.collection_items_id " +
					"WHERE collection_items.deleted_at IS NULL " +
					"GROUP BY collections_id" +
					") AS items_count ON items_count.collections_id = collections.id",
			).
			Order("COALESCE(items_count.items_count, 0) " + orderClause)

	default:
		query = query.Order("array_length(collections.likes, 1) DESC")
	}

	if search != "" {
		query = query.Where("collections.name ILIKE ?", "%"+search+"%")
	}

	err = query.Find(&collections).Error

	return collections, totalCount, err
}

func (r *CollectionsRepository) GetByIdWithoutUser(id string) (models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Where("id = ?", id).
		Where("collections.deleted_at IS NULL").
		Order("collections.created DESC").
		First(&collection).Error
	return collection, err
}

func (r *CollectionsRepository) GetCollectionByTag(
	tag string,
	limit string,
	offset string,
	search string,
) ([]models.Collections, int64, error) {
	var collections []models.Collections
	var totalCount int64

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	countQuery := r.db.
		Model(&models.Collections{}).
		Joins("JOIN tags_collections_collections ON tags_collections_collections.collections_id = collections.id").
		Joins("JOIN tags ON tags.id = tags_collections_collections.tags_id").
		Where("LOWER(tags.name) = LOWER(?)", tag).
		Where("collections.deleted_at IS NULL").
		Where("collections.is_private = ?", false)

	if search != "" {
		countQuery = countQuery.Where("collections.name ILIKE ?", "%"+search+"%")
	}

	err := countQuery.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	query := r.db.
		Joins("JOIN tags_collections_collections ON tags_collections_collections.collections_id = collections.id").
		Joins("JOIN tags ON tags.id = tags_collections_collections.tags_id").
		Where("LOWER(tags.name) = LOWER(?)", tag).
		Preload("User").
		Where("collections.is_private = ?", false).
		Where("collections.deleted_at IS NULL").
		Order("collections.created ASC").
		Limit(intLimit).
		Offset(intOffset)

	if search != "" {
		query = query.Where("collections.name ILIKE ?", "%"+search+"%")
	}

	err = query.Find(&collections).Error

	if err != nil {
		return nil, totalCount, err
	}
	return collections, totalCount, nil
}

func (r *CollectionsRepository) UpdateCollection(
	existsCollection *models.Collections,
	updated *models.Collections,
) (*models.Collections, error) {

	if updated.CollectionItems != nil {
		err := r.db.Model(existsCollection).Association("CollectionItems").Replace(updated.CollectionItems)
		if err != nil {
			return nil, err
		}
	}

	if err := r.db.Model(existsCollection).Updates(updated).Error; err != nil {
		return nil, err
	}

	var result models.Collections
	err := r.db.First(&result, existsCollection.ID).Error

	if !result.IsPrivate {
		doc := map[string]interface{}{
			"id":          result.ID.String(),
			"name":        result.Name,
			"description": result.Description,
			"isPrivate":   result.IsPrivate,
			"bannerUrl":   result.BannerURL,
		}

		if err = r.es.IndexDocument(context.Background(), "collections", doc); err != nil {
			log.Printf("Failed to index collection: %v", err)
		}
	}

	return &result, err
}

func (r *CollectionsRepository) DeleteRelation(sourceCollectionId string, itemId string) (bool, error) {
	var collection models.Collections
	if err := r.db.First(&collection, "id = ?", sourceCollectionId).Error; err != nil {
		return false, err
	}

	var item models.CollectionItems
	if err := r.db.First(&item, "id = ?", itemId).Error; err != nil {
		return false, err
	}

	if err := r.db.Model(&collection).Association("CollectionItems").Delete(&item); err != nil {
		return false, err
	}

	return true, nil
}

func (r *CollectionsRepository) UpdateCollectionFull(existsCollection *models.Collections) (
	*models.Collections,
	error,
) {
	err := r.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Model(existsCollection).Select("*").Updates(existsCollection).Error; err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	var result models.Collections

	if err = r.db.
		Preload("User").
		First(&result, "id = ?", existsCollection.ID).
		Error; err != nil {
		return nil, err
	}
	if !result.IsPrivate {
		doc := map[string]interface{}{
			"id":          result.ID.String(),
			"name":        result.Name,
			"description": result.Description,
			"isPrivate":   result.IsPrivate,
			"bannerUrl":   result.BannerURL,
		}

		if err := r.es.IndexDocument(context.Background(), "collections", doc); err != nil {
			log.Printf("Failed to index collection: %v", err)
		}
	}

	return &result, nil
}

func (r *CollectionsRepository) GetCollectionCountByUserLogin(login string) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.Collections{}).
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Where("collections.deleted_at IS NULL").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CollectionsRepository) DeleteCollection(id string) error {
	var itemIDs []string
	err := r.db.Table("lootor.loot.collections_collection_items_collection_items").
		Where("collections_id = ?", id).
		Pluck("collection_items_id", &itemIDs).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(&models.Collections{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	result := r.db.Exec(
		"DELETE FROM lootor.loot.collections_collection_items_collection_items WHERE collections_id = ?",
		id,
	)
	if result.Error != nil {
		return result.Error
	}
	if len(itemIDs) > 0 {
		err = r.db.Delete(&models.CollectionItems{}, "id IN (?)", itemIDs).Error
		if err != nil {
			return err
		}
	}
	if err = r.es.DeleteDocument(context.Background(), "collections", id); err != nil {
		log.Printf("Failed to delete collection from index: %v", err)
	}
	return nil
}

func (r *CollectionsRepository) DeleteCollectionsByUserLogin(login string) error {
	var itemIDs []string
	var collectionsIDs []string

	err := r.db.Table("lootor.loot.collections").
		Where("user_login = ?", login).
		Pluck("id", &collectionsIDs).Error
	if err != nil {
		return err
	}

	err = r.db.Table("lootor.loot.collections_collection_items_collection_items").
		Where("collections_id IN (?)", collectionsIDs).
		Pluck("collection_items_id", &itemIDs).Error
	if err != nil {
		return err
	}

	resultCollections := r.db.Exec(
		"DELETE FROM lootor.loot.collections_collection_items_collection_items WHERE collections_id = ANY(?)",
		pq.Array(collectionsIDs),
	)

	if resultCollections.Error != nil {
		return resultCollections.Error
	}

	resultCollectionsTags := r.db.Exec(
		"DELETE FROM lootor.loot.tags_collections_collections WHERE collections_id = ANY(?)",
		pq.Array(collectionsIDs),
	)

	if resultCollectionsTags.Error != nil {
		return resultCollectionsTags.Error
	}

	if len(collectionsIDs) > 0 {
		err = r.db.Delete(&models.Collections{}, "id IN (?)", collectionsIDs).Error
		if err != nil {
			return err
		}
	}

	if len(itemIDs) > 0 {
		err = r.db.Delete(&models.CollectionItems{}, "id IN (?)", itemIDs).Error
		if err != nil {
			return err
		}
	}

	for _, id := range collectionsIDs {
		if err = r.es.DeleteDocument(context.Background(), "collections", id); err != nil {
			log.Printf("Failed to delete collection from index: %v", err)
		}
	}

	for _, id := range itemIDs {
		if err = r.es.DeleteDocument(context.Background(), "collection_items", id); err != nil {
			log.Printf("Failed to delete collection item from index: %v", err)
		}

	}

	return nil
}

func (r *CollectionsRepository) GetUniqueByName(name string, login string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.Where("LOWER(name) = LOWER(?)", name).Where("LOWER(user_login) = LOWER(?)", login).First(&collection)

	if err.Error != nil {
		return nil, err.Error
	}

	return &collection, nil
}

func (r *CollectionsRepository) IncrementCommentsCount(id string, amount int) error {
	return r.db.Model(&models.Collections{}).
		Where("id = ?", id).
		Update("comments_count", gorm.Expr("COALESCE(comments_count, 0) + ?", amount)).Error
}

func (r *CollectionsRepository) DecrementCommentsCount(id string, amount int) error {
	return r.db.Model(&models.Collections{}).
		Where("id = ?", id).
		Update("comments_count", gorm.Expr("GREATEST(COALESCE(comments_count, 0) - ?, 0)", amount)).Error
}

func (r *CollectionsRepository) GetCollectionsCount(login string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.Collections{}).
		Where("LOWER(user_login) = LOWER(?)", login).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

func (r *CollectionsRepository) GetCollectionsByIdsMap(
	ids []string,
	counts map[uuid.UUID]int64,
	totalPrices map[uuid.UUID]float64,
	shippingCosts map[uuid.UUID]float64,
	authorizedUser string,
) (map[string]dto.CollectionsResponse, error) {
	var collections []models.Collections
	err := r.db.Where("id IN (?)", ids).Preload("User").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	collectionsMap := make(map[string]dto.CollectionsResponse, len(collections))
	for _, collection := range collections {

		var temp dto.CollectionsResponse
		err = mapstructure.Decode(collection, &temp)
		if err != nil {
			return nil, err
		}
		temp.CollectionItemsCount = counts[collection.ID]
		temp.TotalPrice = totalPrices[collection.ID]
		temp.ShippingTotal = shippingCosts[collection.ID]
		temp.LikesCount = int64(len(collection.Likes))
		temp.CanLike = false
		temp.IsOwner = authorizedUser == temp.User.Login
		if authorizedUser != "" {
			temp.CanLike = !slices.Contains(collection.Likes, authorizedUser)
		}

		collectionsMap[collection.ID.String()] = temp
	}
	return collectionsMap, nil
}

func (r *CollectionsRepository) GetCollectionsByIdsMapForShort(ids []string) (map[string]dto.CollectionShort, error) {
	var collections []models.Collections
	err := r.db.Where("id IN (?)", ids).Preload("User").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	collectionsMap := make(map[string]dto.CollectionShort, len(collections))
	for _, collection := range collections {
		var temp dto.CollectionShort
		err = mapstructure.Decode(collection, &temp)
		if err != nil {
			return nil, err
		}
		collectionsMap[collection.ID.String()] = temp
	}
	return collectionsMap, nil
}

func (r *CollectionsRepository) GetCollectionsByTranslitsMapForShort(
	ids []string,
	userLogin string,
) (map[string]dto.CollectionShort, error) {
	var collections []models.Collections
	err := r.db.Where("id IN (?) AND user_login = ?", ids, userLogin).Preload("User").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	collectionsMap := make(map[string]dto.CollectionShort, len(collections))
	for _, collection := range collections {
		var temp dto.CollectionShort
		err = mapstructure.Decode(collection, &temp)
		if err != nil {
			return nil, err
		}
		collectionsMap[collection.ID.String()] = temp
	}
	return collectionsMap, nil
}

func (r *CollectionsRepository) GetCollectionsByIds(
	ids []string,
	counts map[uuid.UUID]int64,
	totalPrices map[uuid.UUID]float64,
	shippingCosts map[uuid.UUID]float64,
	authorizedUser string,
	tagsMap map[string]*microservices.GetShortsResponse,
) ([]dto.CollectionsResponse, error) {
	var collections []models.Collections
	err := r.db.Where("id IN (?)", ids).Preload("User").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	var result []dto.CollectionsResponse
	for _, collection := range collections {
		var temp dto.CollectionsResponse
		err = mapstructure.Decode(collection, &temp)
		if err != nil {
			return nil, err
		}
		temp.CollectionItemsCount = counts[collection.ID]
		temp.TotalPrice = totalPrices[collection.ID]
		temp.ShippingTotal = shippingCosts[collection.ID]
		temp.LikesCount = int64(len(collection.Likes))
		temp.CanLike = true
		temp.IsOwner = authorizedUser == temp.User.Login
		if authorizedUser != "" {
			temp.CanLike = !slices.Contains(collection.Likes, authorizedUser)
		}
		if len(tagsMap) > 0 {
			tagsProto := tagsMap[collection.ID.String()].Tags
			tempTags := utils.NormalizeTagsShort(tagsProto)
			temp.Tags = tempTags
		}

		result = append(result, temp)
	}
	return result, nil
}

func (r *CollectionsRepository) GetCollectionItemsIds(collectionId string) ([]string, error) {
	var itemIDs []string
	collectionItemsIds := r.db.Table("collections_collection_items_collection_items").Where(
		"collections_id = ?",
		collectionId,
	).Select("collection_items_id").Unscoped().
		Pluck("collection_items_id", &itemIDs)

	if collectionItemsIds.Error != nil {
		return nil, collectionItemsIds.Error
	}

	return itemIDs, nil

}

func (r *CollectionsRepository) GetLikesOfAllCollectionsUser(login string) (int64, error) {
	var totalLikes int64

	query := `
        SELECT COALESCE(SUM(array_length(likes, 1)), 0) as total_likes
        FROM collections 
        WHERE LOWER(user_login) = LOWER($1) 
          AND deleted_at IS NULL
    `

	err := r.db.Raw(query, login).Scan(&totalLikes).Error
	return totalLikes, err
}
