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
	{"003_move_events", migrations.MoveEvents},
	{"004_add_deleted_flag_events", migrations.AddDeletedFlagEvents},
	{"005_fix_deleted_events", migrations.MarkEventsDeletedForDeletedTargets},
	{"006_move_tags_and_entities", migrations.MoveTagsAndEntities},
	{"006_5_set_all_roles", migrations.SetAllRoles},
	{"006_6_entity_id_to_string", migrations.EntityIDToString},
	{"007_uppercase_tags_names", migrations.UppercaseTagsNames},
	{"008_delete_tags_duplicates", migrations.DeleteTagsDuplicates},
	{"010_add_achievements_to_all_users", migrations.AddAchievementsToAllUsers},
	{"011_update_user_achievements_from_existing_data", migrations.UpdateUserAchievementsFromExistingData},
	{"012_update_beta_tester_achievement", migrations.UpdateBetaTesterAchievement},
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
