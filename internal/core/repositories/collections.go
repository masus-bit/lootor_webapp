package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
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

	doc := map[string]interface{}{
		"id":   collection.Id.String(),
		"name": collection.Name,
	}

	if err := r.es.IndexDocument(context.Background(), "collections", doc); err != nil {
		log.Printf("Failed to index collection: %v", err)
	}

	return collection, nil
}

func (r *CollectionsRepository) FindAllCollections() ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.Find(&collections).Error
	return collections, err
}

func (r *CollectionsRepository) GetCollectionById(id string, limit string, offset string, orderByInput string, order string, search string) (*models.Collections, error) {
	var collection models.Collections

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems", func(tx *gorm.DB) *gorm.DB {

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
		}).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Preload("CollectionItems.ItemType").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetCollectionByIdWithoutLimits(id string) (*models.Collections, error) {
	var collection models.Collections

	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Preload("CollectionItems.ItemType").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetByShareString(shareString string, limit string, offset string, orderByInput string, order string, search string) (*models.Collections, error) {
	var collection models.Collections

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems", func(tx *gorm.DB) *gorm.DB {

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
		}).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Preload("CollectionItems.ItemType").
		Where("share_string = ? AND deleted = ?", shareString, false).
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetByIdWithoutCollectionItems(id string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("User").
		Preload("Tags").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetOneByTransliteration(login string, transliteration string, limit string, offset string, orderByInput string, order string, search string) (*models.Collections, error) {
	var collection models.Collections

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	err := r.db.
		Where("collections.transliteration = ?", transliteration).
		Where("collections.deleted = ?", false).
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems", func(tx *gorm.DB) *gorm.DB {

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
				return tx.Offset(intOffset).Limit(intLimit).Order("purchase_price ASC, collection_items.id ASC")

			default:
				return tx.Offset(intOffset).Limit(intLimit).Order("created_at DESC, collection_items.id DESC")
			}
		}).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
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
		Preload("Tags").
		Preload("CollectionItems", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("created DESC")
		}).
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Preload("CollectionItems.Owner").
		Preload("CollectionItems.ItemType").
		Order("created DESC").
		Find(&collections).Error
	return collections, err
}

func (r *CollectionsRepository) GetByUserIdWithoutCollectionItems(login string, search string) ([]models.Collections, error) {
	var collections []models.Collections

	query := r.db.
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Preload("Tags").
		Order("created DESC")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&collections).Error

	return collections, err
}

func (r *CollectionsRepository) GetByUserIdWithoutPrivates(login string, search string) ([]models.Collections, error) {
	var collections []models.Collections
	query := r.db.
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Preload("Tags").
		Where("collections.is_private = ?", false).
		Order("collections.created DESC")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Find(&collections).Error

	return collections, err
}

func (r *CollectionsRepository) GetByIdWithoutUser(id string) (models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("Tags").
		Where("id = ? AND deleted = ?", id, false).
		Order("collections.created DESC").
		First(&collection).Error
	return collection, err
}

func (r *CollectionsRepository) GetCollectionByTag(tag string, limit string, offset string, search string) ([]models.Collections, int64, error) {
	var collections []models.Collections
	var totalCount int64

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	countQuery := r.db.
		Model(&models.Collections{}).
		Joins("JOIN tags_collections_collections ON tags_collections_collections.collections_id = collections.id").
		Joins("JOIN tags ON tags.id = tags_collections_collections.tags_id").
		Where("LOWER(tags.name) = LOWER(?)", tag).
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
		Preload("Tags").
		Where("collections.is_private = ?", false).
		Order("collections.created ASC").
		Limit(intLimit).
		Offset(intOffset)

	if search != "" {
		query = query.Where("collections.name ILIKE ?", "%"+search+"%").Where("collections.delete = ?", false)
	}

	err = query.Find(&collections).Error

	if err != nil {
		return nil, totalCount, err
	}
	return collections, totalCount, nil
}

func (r *CollectionsRepository) UpdateCollection(existsCollection *models.Collections, updated *models.Collections) (*models.Collections, error) {
	if updated.Tags != nil {
		err := r.db.Model(existsCollection).Association("Tags").Replace(updated.Tags)
		if err != nil {
			return nil, err
		}
	}

	if err := r.db.Model(existsCollection).Updates(updated).Error; err != nil {
		return nil, err
	}

	var result models.Collections
	err := r.db.Preload("Tags").First(&result, existsCollection.Id).Error

	doc := map[string]interface{}{
		"id":   result.Id.String(),
		"name": result.Name,
	}

	if err := r.es.IndexDocument(context.Background(), "collections", doc); err != nil {
		log.Printf("Failed to index collection: %v", err)
	}

	return &result, err
}

func (r *CollectionsRepository) GetCollectionCountByUserLogin(login string) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.Collections{}).
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Where("collections.deleted = ?", false).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CollectionsRepository) DeleteCollection(id string) error {
	err := r.db.Delete(&models.Collections{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	if err := r.es.DeleteDocument(context.Background(), "collections", id); err != nil {
		log.Printf("Failed to delete collection from index: %v", err)
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
