package migrations

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

func SetAllRoles(db *gorm.DB) error {
	adminLogins := []string{"seymor", "gseymor", "godSeymor", "eg_rif_ykkur_i_bita"}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Users{}).
			Where("role IS NULL OR role = ''").
			Update("role", "user").Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Users{}).
			Where("login IN ?", adminLogins).
			Update("role", "admin").Error; err != nil {
			return err
		}

		return nil
	})
}
