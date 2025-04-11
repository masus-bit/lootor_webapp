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

func (r *Repository) Create(collection *models.Collection) error {
	return r.db.Create(collection).Error
}

func (r *Repository) FindAll() ([]models.Collection, error) {
	var collections []models.Collection
	err := r.db.Find(&collections).Error
	return collections, err
}

func (r *Repository) GetById(id string) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItem").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *Repository) GetByShareString(shareString string) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItem").
		Where("share_string = ? AND deleted = ?", shareString, false).
		First(&collection).Error
	return &collection, err
}

func (r *Repository) GetByIdWithoutCollectionItems(id string) (*models.Collection, error) {
	var collection models.Collection
	err := r.db.
		Preload("User").
		Preload("Tags").
		Where("id = ? AND deleted = ?", id, false).
		First(&collection).Error
	return &collection, err
}

func (r *Repository) GetOneByTransliteration(login string, transliteration string) (*models.Collection, error) {
	var collection models.Collection

	err := r.db.
		Where("transliteration = ? AND deleted = ?", transliteration, false).
		Joins("JOIN users ON collections.user_id = users.login AND LOWER(users.login) = LOWER(?)", login).
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

// async getByUserId(id: string): Promise<Collection[]> {
// return await this.collectionRepository
// .createQueryBuilder('collection')
// .leftJoinAndSelect('collection.user', 'user')
// .leftJoinAndSelect('collection.tags', 'tags')
// .leftJoinAndSelect('collection.collectionItems', 'collection_item')
// .leftJoinAndSelect('collection_item.platform', 'platforms')
// .leftJoinAndSelect('collection_item.entities', 'entity_model')
// .where('LOWER(user.login) = LOWER(:id)', { id })
// .andWhere('collection.deleted = :deleted', { deleted: false })
// .orderBy('collection.created', 'DESC')
// .getMany();
// }
func (r *Repository) GetByUserId(login string) ([]models.Collection, error) {
	var collections []models.Collection
	err := r.db.
		Preload("User").
		Preload("Tags").
		Preload("CollectionItems").
		Preload("CollectionItems.Platform").
		Preload("CollectionItems.Entities").
		Where("user_id = ? AND deleted = ?", login, false).
		Find(&collections).Error
	return collections, err
}
