package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) CreateUser(user *models.Users) error {
	return r.db.Create(user).Error
}

func (r *UsersRepository) FindAllUsers() ([]models.Users, error) {
	var users []models.Users
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UsersRepository) GetByEmailWithPassword(email string) (*models.Users, error) {
	var user models.Users
	err := r.db.Select("*").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepository) GetByLoginWithPassword(login string) (*models.Users, error) {
	var user models.Users

	err := r.db.Select("*").Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByLogin(login string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("user_name = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByEmail(email string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByVerificationToken(token string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("verification_token = ?", token).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByVkId(vkId string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("vk_id = ?", vkId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByTgId(tgId string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("telegram_id = ?", tgId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) UpdateUser(existsUser *models.Users, updated models.Users) (*models.Users, error) {
	result := r.db.Model(existsUser).Select("*").Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser models.Users
	err := r.db.Where("login = ?", existsUser.Login).First(&updatedUser).Error
	return &updatedUser, err
}
