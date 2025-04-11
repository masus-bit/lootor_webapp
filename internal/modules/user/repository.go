package user

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

func (r *Repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *Repository) GetByEmailWithPassword(email string) (*models.User, error) {
	var user models.User
	err := r.db.Select("*").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByLoginWithPassword(login string) (*models.User, error) {
	var user models.User

	err := r.db.Select("*").Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByLogin(login string) (*models.User, error) {
	var user models.User

	err := r.db.Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByUsername(username string) (*models.User, error) {
	var user models.User

	err := r.db.Where("user_name = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByVerificationToken(token string) (*models.User, error) {
	var user models.User

	err := r.db.Where("verification_token = ?", token).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByVkId(vkId string) (*models.User, error) {
	var user models.User

	err := r.db.Where("vk_id = ?", vkId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByTgId(tgId string) (*models.User, error) {
	var user models.User

	err := r.db.Where("telegram_id = ?", tgId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(existsUser models.User, updated models.User) (*models.User, error) {
	result := r.db.Model(existsUser).Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser models.User
	err := r.db.First(&updatedUser, existsUser.Login).Error
	return &updatedUser, err
}
