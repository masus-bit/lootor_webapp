package migrations

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"lootor/internal/pkg/utils"
)

func FillUserIds(db *gorm.DB) error {

	var users []models.Users

	err := db.Model(&models.Users{}).Find(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		utils.Slugify(user.Login)
		err = db.Model(&models.Users{}).
			Where("login = ?", user.Login).
			Update("id", utils.Slugify(user.Login)).
			Error
	}
	return err
}
