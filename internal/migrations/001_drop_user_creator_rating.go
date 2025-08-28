package migrations

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

func DropUserCreatorRating(db *gorm.DB) error {
	return db.Migrator().DropColumn(&models.Users{}, "creator_rating")
}
