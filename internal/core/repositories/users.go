package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/pkg/elasticsearch"
	"strconv"
	"time"
)

type UsersRepository struct {
	es              *elasticsearch.ElasticService
	db              *gorm.DB
	collectionsRepo *CollectionsRepository
}

func NewUsersRepository(db *gorm.DB, es *elasticsearch.ElasticService, collectionsRepo *CollectionsRepository) *UsersRepository {
	return &UsersRepository{db: db, es: es, collectionsRepo: collectionsRepo}
}

func (r *UsersRepository) CreateUser(user *models.Users) error {

	if user == nil {
		return errors.New("user is nil")
	}

	if r.db == nil {
		return errors.New("database connection is nil")
	}

	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	if r.es == nil {
		log.Printf("ElasticSearch client is nil, skipping indexing")
		return nil
	}

	doc := map[string]interface{}{
		"id":          user.Login,
		"login":       user.Login,
		"userName":    user.UserName,
		"email":       user.Email,
		"avatarUrl":   user.AvatarUrl,
		"vkId":        user.VkId,
		"tgId":        user.TelegramId,
		"isPremium":   user.IsPremium,
		"profileName": user.ProfileName,
	}

	if err := r.es.IndexDocument(context.Background(), "users", doc); err != nil {
		log.Printf("Failed to index user: %v", err)
	}

	return nil
}

func (r *UsersRepository) DeleteUser(login string) error {

	err := r.collectionsRepo.DeleteCollectionsByUserLogin(login)
	if err != nil {
		return err
	}

	resultEvents := r.db.Exec(
		"DELETE FROM lootor.loot.events WHERE initiator_login = ?",
		login,
	)

	if resultEvents.Error != nil {
		return resultEvents.Error
	}
	resultWishListItems := r.db.Exec(
		"DELETE FROM lootor.loot.wish_list_items WHERE user_login = ?",
		login,
	)

	if resultWishListItems.Error != nil {
		return resultWishListItems.Error
	}

	//err = r.db.Delete(&models.Users{}, "login = ?", login).Error
	now := time.Now()
	err = r.db.Model(&models.Users{}).
		Where("login = ?", login).
		Update("deleted_at", now).Error
	if err != nil {
		return err
	}

	if err = r.es.DeleteDocument(context.Background(), "users", login); err != nil {
		log.Printf("Failed to delete user from index: %v", err)
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
	err := r.db.Select("*").Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepository) GetByLoginWithPassword(login string) (*models.Users, error) {
	var user models.Users

	err := r.db.Select("*").Where("login = ? AND deleted_at IS NULL", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByLogin(login string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("LOWER(login) = LOWER(?) AND deleted_at IS NULL", login).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByLoginForSignUp(login string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("LOWER(login) = LOWER(?)", login).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("user_name = ? AND deleted_at IS NULL", username).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByEmail(email string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("email = ? AND deleted_at IS NULL", email).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByEmailForSignUp(email string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("email = ?", email).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByVerificationToken(token string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("verification_token = ? AND deleted_at IS NULL", token).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByResetToken(token string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("reset_token = ? AND deleted_at IS NULL", token).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByVkId(vkId string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("vk_id = ? AND deleted_at IS NULL", vkId).Preload("WishListItems").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetByTgId(tgId string) (*models.Users, error) {
	var user models.Users

	err := r.db.Where("telegram_id = ? AND deleted_at IS NULL", tgId).Preload("WishListItems").First(&user).Error

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
		"id":          updatedUser.Login,
		"login":       updatedUser.Login,
		"userName":    updatedUser.UserName,
		"email":       updatedUser.Email,
		"avatarUrl":   updatedUser.AvatarUrl,
		"vkId":        updatedUser.VkId,
		"tgId":        updatedUser.TelegramId,
		"isPremium":   updatedUser.IsPremium,
		"profileName": updatedUser.ProfileName,
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
		"id":          updatedUser.Login,
		"login":       updatedUser.Login,
		"userName":    updatedUser.UserName,
		"email":       updatedUser.Email,
		"avatarUrl":   updatedUser.AvatarUrl,
		"vkId":        updatedUser.VkId,
		"tgId":        updatedUser.TelegramId,
		"isPremium":   updatedUser.IsPremium,
		"profileName": updatedUser.ProfileName,
	}

	if err := r.es.IndexDocument(context.Background(), "users", doc); err != nil {
		log.Printf("Failed to index user: %v", err)
	}

	return &updatedUser, err
}

func (r *UsersRepository) UpdateLogin(existsUser string, updated *models.Users) (*models.Users, error) {
	result := r.db.Model(&models.Users{}).Where("login = ?", existsUser).Select("*").Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser models.Users
	err := r.db.Where("login = ?", updated.Login).Preload("WishListItems").First(&updatedUser).Error

	doc := map[string]interface{}{
		"id":          updatedUser.Login,
		"login":       updatedUser.Login,
		"userName":    updatedUser.UserName,
		"email":       updatedUser.Email,
		"avatarUrl":   updatedUser.AvatarUrl,
		"vkId":        updatedUser.VkId,
		"tgId":        updatedUser.TelegramId,
		"isPremium":   updatedUser.IsPremium,
		"profileName": updatedUser.ProfileName,
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

func (r *UsersRepository) AddRating(existsUser *models.Users, updatedUser *models.Users) (bool, error) {
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

func (r *UsersRepository) GetAll(search string, limit, offset string) ([]models.Users, error) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)
	var users []models.Users
	query := r.db.
		Order("users.login ASC").
		Limit(intLimit).
		Offset(intOffset)

	if search != "" {
		query = query.Where("login ILIKE ? AND deleted_at IS NULL", "%"+search+"%")
	}

	err := query.Find(&users).Error

	return users, err
}

func (r *UsersRepository) DeactivatePremium() error {
	var expiredPremiumLogins []string
	result := r.db.Model(&models.Users{}).
		Select("login").
		Where("is_premium = ? AND premium_until < ?", true, time.Now()).
		Pluck("login", &expiredPremiumLogins)
	if result.Error != nil {
		return result.Error
	}

	if len(expiredPremiumLogins) > 0 {
		updateResult := r.db.Model(&models.Users{}).
			Where("login IN ?", expiredPremiumLogins).
			Updates(map[string]interface{}{"is_premium": false})
		if updateResult.Error != nil {
			return updateResult.Error
		}

		subscriptionUpdateResult := r.db.Model(&models.Subscription{}).
			Where("user_login IN ?", expiredPremiumLogins).
			Updates(map[string]interface{}{"status": "expired"})
		if subscriptionUpdateResult.Error != nil {
			return subscriptionUpdateResult.Error
		}
	}

	return nil
}

func (r *UsersRepository) DeleteDeletedUsersForever() error {
	var deletedLogins []string
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	result := r.db.Model(&models.Users{}).
		Select("login").
		Where("deleted_at IS NOT NULL AND deleted_at < ?", oneYearAgo).
		Pluck("login", &deletedLogins)
	if result.Error != nil {
		return result.Error
	}

	if len(deletedLogins) > 0 {
		err := r.db.Delete(&models.Users{}, "login IN ?", deletedLogins).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UsersRepository) GetForSubs(users []string) ([]models.SubUsers, error) {
	var subUsers []models.SubUsers
	err := r.db.Model(&models.Users{}).
		Where("login IN (?)", users).
		Select("login", "avatar_url", "profile_name", "is_premium").
		Find(&subUsers).Error
	return subUsers, err
}

func (r *UsersRepository) IncrementExperience(userLogin string, amount int) error {
	return r.db.Model(&models.Users{}).
		Where("login = ?", userLogin).
		Update("exp", gorm.Expr("COALESCE(exp, 0) + ?", amount)).Error
}

func (r *UsersRepository) DecrementExperience(userLogin string, amount int) error {
	return r.db.Model(&models.Users{}).
		Where("login = ?", userLogin).
		Update("exp", gorm.Expr("COALESCE(exp, 0) - ?", amount)).Error
}
