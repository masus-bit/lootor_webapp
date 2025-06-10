package repositories

import (
	"context"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
	"time"
)

type UsersRepository struct {
	es *elasticsearch.ElasticService
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB, es *elasticsearch.ElasticService) *UsersRepository {
	return &UsersRepository{db: db, es: es}
}

func (r *UsersRepository) CreateUser(user *models.Users) error {

	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	doc := map[string]interface{}{
		"id":        user.Login,
		"login":     user.Login,
		"userName":  user.UserName,
		"email":     user.Email,
		"avatarUrl": user.AvatarUrl,
		"vkId":      user.VkId,
		"tgId":      user.TelegramId,
	}

	if err := r.es.IndexDocument(context.Background(), "users", doc); err != nil {
		log.Printf("Failed to index user: %v", err)
	}

	return nil
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

	err := r.db.Where("LOWER(login) = LOWER(?)", login).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("user_name = ?", username).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByEmail(email string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("email = ?", email).Preload("WishListItems").First(&user).Error

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

func (r *UsersRepository) GetByResetToken(token string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("reset_token = ?", token).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByVkId(vkId string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("vk_id = ?", vkId).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByTgId(tgId string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("telegram_id = ?", tgId).Preload("WishListItems").First(&user).Error

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
	err := r.db.Where("login = ?", existsUser.Login).Preload("WishListItems").First(&updatedUser).Error

	doc := map[string]interface{}{
		"id":        updatedUser.Login,
		"login":     updatedUser.Login,
		"userName":  updatedUser.UserName,
		"email":     updatedUser.Email,
		"avatarUrl": updatedUser.AvatarUrl,
		"vkId":      updatedUser.VkId,
		"tgId":      updatedUser.TelegramId,
	}

	if err := r.es.IndexDocument(context.Background(), "users", doc); err != nil {
		log.Printf("Failed to index user: %v", err)
	}

	return &updatedUser, err
}

func (r *UsersRepository) UpdateUserFull(existsUser *models.Users) (*models.Users, error) {
	result := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(existsUser)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser models.Users
	err := r.db.Where("login = ?", existsUser.Login).Preload("WishListItems").First(&updatedUser).Error

	doc := map[string]interface{}{
		"id":        updatedUser.Login,
		"login":     updatedUser.Login,
		"userName":  updatedUser.UserName,
		"email":     updatedUser.Email,
		"avatarUrl": updatedUser.AvatarUrl,
		"vkId":      updatedUser.VkId,
		"tgId":      updatedUser.TelegramId,
	}

	if err := r.es.IndexDocument(context.Background(), "users", doc); err != nil {
		log.Printf("Failed to index user: %v", err)
	}

	return &updatedUser, err
}

func (r *UsersRepository) UpdateLogin(existsUser *models.Users, updated *models.Users) (*models.Users, error) {
	result := r.db.Model(existsUser).Select("*").Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser models.Users
	err := r.db.Where("login = ?", existsUser.Login).Preload("WishListItems").First(&updatedUser).Error

	doc := map[string]interface{}{
		"id":        updatedUser.Login,
		"login":     updatedUser.Login,
		"userName":  updatedUser.UserName,
		"email":     updatedUser.Email,
		"avatarUrl": updatedUser.AvatarUrl,
		"vkId":      updatedUser.VkId,
		"tgId":      updatedUser.TelegramId,
	}

	if err := r.es.IndexDocument(context.Background(), "users", doc); err != nil {
		log.Printf("Failed to index user: %v", err)
	}

	return &updatedUser, err
}

func (r *UsersRepository) ActivatePremium(existsUser *models.Users, updatedUser *models.Users) (bool, error) {
	result := r.db.Model(existsUser).Select("*").Updates(updatedUser)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *UsersRepository) CheckPremiumStatus(userLogin string) (bool, error) {
	var user models.Users
	if err := r.db.Where("login = ?", userLogin).First(&user).Error; err != nil {
		return false, err
	}
	return user.IsPremium && user.PremiumUntil.After(time.Now()), nil
}

func (r *UsersRepository) DeactivatePremium() error {
	result := r.db.Model(&models.Users{}).
		Where("is_premium = ? AND premium_until < ?", true, time.Now()).
		Updates(map[string]interface{}{"is_premium": false})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
