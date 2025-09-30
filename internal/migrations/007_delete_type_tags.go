package migrations

import (
	"gorm.io/gorm"
)

func DeleteTypeTags(db *gorm.DB) error {
	return db.Migrator().DropColumn("loot_tags.tags", "type")
}
