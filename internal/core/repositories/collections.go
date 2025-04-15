package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"strings"
)

type CollectionsRepository struct {
	db *gorm.DB
}

func NewCollectionsRepository(db *gorm.DB) *CollectionsRepository {
	return &CollectionsRepository{db: db}
}

func (r *CollectionsRepository) CreateCollection(collection *models.Collections) (*models.Collections, error) {
	err := r.db.Create(&collection)
	if err.Error != nil {
		return nil, err.Error
	}
	return collection, nil
}

func (r *CollectionsRepository) FindAllCollections() ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.Find(&collections).Error
	return collections, err
}

func (r *CollectionsRepository) GetCollectionById(id string) (*models.Collections, error) {
	fmt.Println(id)
	var collection models.Collections
	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *CollectionsRepository) GetByShareString(shareString string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
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

func (r *CollectionsRepository) GetOneByTransliteration(login string, transliteration string) (*models.Collections, error) {
	var collection models.Collections

	err := r.db.
		Where("collections.transliteration = ?", transliteration).
		Where("collections.deleted = ?", false).
		Where("user_login = LOWER(?)", strings.ToLower(login)).
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
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
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Order("created DESC").
		Find(&collections).Error
	return collections, err
}

func (r *CollectionsRepository) GetByUserIdWithoutCollectionItems(login string) ([]models.Collections, error) {
	var collections []models.Collections

	err := r.db.
		Where("LOWER(user_login) = LOWER(?)", login).
		Preload("User").
		Preload("Tags").
		Order("created DESC").
		Find(&collections).Error

	return collections, err
}

func (r *CollectionsRepository) GetByUserIdWithoutPrivates(login string) ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.
		Where("user_login = LOWER(?)", strings.ToLower(login)).
		Preload("User").
		Preload("Tags").
		Where("collections.is_private = ?", false).
		Order("collections.created DESC").
		Find(&collections).Error
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

func (r *CollectionsRepository) GetCollectionByTag(tag string) ([]models.Collections, error) {
	var collections []models.Collections

	err := r.db.
		Joins("JOIN tags_collections_collections ON tags_collections_collections.collectionsId = collections.id").
		Joins("JOIN tags ON tags.id = tags_collections_collections.tagsId AND tags.name = ?", tag).
		Preload("User").
		Preload("Tags").
		Where("collections.is_private = ?", false).
		Order("collections.created ASC").
		Find(&collections).Error

	if err != nil {
		return nil, err
	}
	return collections, nil
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
	return &result, err
}

func (r *CollectionsRepository) GetCollectionCountByUserLogin(login string) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.Collections{}).
		Where("user_login = LOWER(?)", strings.ToLower(login)).
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

	return nil
}
