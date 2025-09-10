package migrations

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

func DropPostsCount(db *gorm.DB) error {
	return db.Migrator().DropColumn(&models.Users{}, "post_count")
}
