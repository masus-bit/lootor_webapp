package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func MoveEvents(db *gorm.DB) error {
	sourceSchema := "loot"
	targetSchema := "loot_events"
	tableName := "events"

	createSchemaSQL := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, targetSchema)
	if err := db.Exec(createSchemaSQL).Error; err != nil {
		return fmt.Errorf("failed to create target schema '%s': %w", targetSchema, err)
	}

	moveTableSQL := fmt.Sprintf(`ALTER TABLE %s.%s SET SCHEMA %s;`, sourceSchema, tableName, targetSchema)
	if err := db.Exec(moveTableSQL).Error; err != nil {
		return fmt.Errorf("failed to move table '%s' to target schema '%s': %w", tableName, targetSchema, err)
	}

	return nil
}
