package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetByEmailWithPassword(email string) (*User, error) {
	var user User
	err := r.db.Select("*").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByLoginWithPassword(login string) (*User, error) {
	var user User

	err := r.db.Select("*").Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByLogin(login string) (*User, error) {
	var user User

	err := r.db.Where("login = ?", login).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var user User

	err := r.db.Where("user_name = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	var user User

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByVerificationToken(token string) (*User, error) {
	var user User

	err := r.db.Where("verification_token = ?", token).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByVkId(vkId string) (*User, error) {
	var user User

	err := r.db.Where("vk_id = ?", vkId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByTgId(tgId string) (*User, error) {
	var user User

	err := r.db.Where("telegram_id = ?", tgId).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(existsUser User, updated User) (*User, error) {
	result := r.db.Model(existsUser).Updates(updated)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedUser User
	err := r.db.First(&updatedUser, existsUser.Login).Error
	return &updatedUser, err
}
