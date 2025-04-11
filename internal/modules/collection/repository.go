package collection

import (
	"errors"
	"gorm.io/gorm"
	"lootor/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(collection *models.Collections) error {
	return r.db.Create(collection).Error
}

func (r *Repository) FindAll() ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.Find(&collections).Error
	return collections, err
}

func (r *Repository) GetById(id string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("Users").
		Preload("Tags").
		Preload("CollectionItems").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *Repository) GetByShareString(shareString string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("Users").
		Preload("Tags").
		Preload("CollectionItems").
		Where("share_string = ? AND deleted = ?", shareString, false).
		First(&collection).Error
	return &collection, err
}

func (r *Repository) GetByIdWithoutCollectionItems(id string) (*models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("Users").
		Preload("Tags").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *Repository) GetOneByTransliteration(login string, transliteration string) (*models.Collections, error) {
	var collection models.Collections

	err := r.db.
		Where("collections.transliteration = ?", transliteration).
		Where("collections.deleted = ?", false).
		Joins("User", func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(users.login) = LOWER(?)", login)
		}).
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

func (r *Repository) GetByUserId(login string) ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.
		Joins("User", func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(users.login) = LOWER(?)", login)
		}).
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Order("collection.created DESC").
		Find(&collections).Error
	return collections, err
}

func (r *Repository) GetByUserIdWithoutCollectionItems(login string) ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.
		Joins("User", func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(users.login) = LOWER(?)", login)
		}).
		Preload("Tags").
		Order("collection.created DESC").
		Find(&collections).Error
	return collections, err
}

func (r *Repository) GetByUserIdWithoutPrivates(login string) ([]models.Collections, error) {
	var collections []models.Collections
	err := r.db.
		Joins("User", func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(users.login) = LOWER(?)", login)
		}).
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Where("collections.is_private = ?", false).
		Order("collection.created DESC").
		Find(&collections).Error
	return collections, err
}

func (r *Repository) GetByIdWithoutUser(id string) (models.Collections, error) {
	var collection models.Collections
	err := r.db.
		Preload("Tags").
		Where("id = ? AND deleted = ?", id, false).
		Order("collection.created DESC").
		First(&collection).Error
	return collection, err
}

func (r *Repository) GetByTag(tag string) ([]models.Collections, error) {
	var collections []models.Collections

	err := r.db.
		Joins("JOIN tags_collections_collections ON tags_collections_collections.collections_id = collections.id").
		Joins("JOIN tags ON tags.id = tags_collections_collections.tags_id AND tags.name = ?", tag).
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Where("collections.is_private = ?", false).
		Order("collections.created ASC").
		Find(&collections).Error

	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (r *Repository) Update(existsCollection models.Collections, updated models.Collections) (*models.Collections, error) {
	result := r.db.Model(existsCollection).Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedCollection models.Collections
	err := r.db.First(&updatedCollection, existsCollection.Id).Error
	return &updatedCollection, err
}

func (r *Repository) GetCountByUserLogin(login string) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.Collections{}).
		Joins("User", func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(users.login) = LOWER(?)", login)
		}).
		Where("collections.deleted = ?", false).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) Delete(id string) error {
	err := r.db.Delete(&models.Collections{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
