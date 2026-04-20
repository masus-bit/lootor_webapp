package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type UserSettingsRepository struct {
	db *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) *UserSettingsRepository {
	return &UserSettingsRepository{db: db}
}

func (r *UserSettingsRepository) CreateRecord(sets *models.UserSettings) (*models.UserSettings, error) {

	err := r.db.Create(sets)
	if err.Error != nil {
		return nil, err.Error
	}

	return sets, nil
}

func (r *UserSettingsRepository) FindRecordByUserLogin(userLogin string) (*models.UserSettings, error) {
	var userSettings *models.UserSettings

	_ = r.db.Model(&userSettings).Where("user_login = ?", userLogin).First(&userSettings)

	return userSettings, nil
}

func (r *UserSettingsRepository) UpdateFull(existsSets *models.UserSettings) (*models.UserSettings, error) {
	err := r.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Model(existsSets).Where(
				"user_login = ?",
				existsSets.UserLogin,
			).Select("*").Updates(existsSets).Error; err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	var result models.UserSettings

	if err = r.db.Model(&models.UserSettings{}).
		First(&result, "id = ?", existsSets.ID).
		Error; err != nil {
		return nil, err
	}

	return &result, nil
}
