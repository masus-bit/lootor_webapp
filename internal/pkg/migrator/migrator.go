package migrator

import (
	"fmt"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"lootor/internal/migrations"
)

type MigrationFunc func(*gorm.DB) error

var migrationList = []struct {
	name string
	fn   MigrationFunc
}{
	{"001_drop_user_creator_rating", migrations.DropUserCreatorRating},
	{"002_change_post_id_to_uint", migrations.ChangePostIdToUint64AutoIncrement},
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Migrations{}); err != nil {
		return fmt.Errorf("пизда миграции: %w", err)
	}

	for _, m := range migrationList {
		var count int64
		if err := db.Model(&models.Migrations{}).Where("name = ?", m.name).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count == 0 {
			tx := db.Begin()
			if err := m.fn(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("migration %s failed: %w", m.name, err)
			}

			if err := tx.Create(&models.Migrations{Name: m.name}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration: %w", err)
			}

			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("failed to commit migration transaction: %w", err)
			}

			fmt.Printf("Applied migration: %s\n", m.name)
		}
	}

	return nil
}
