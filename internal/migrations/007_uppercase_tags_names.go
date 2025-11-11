package migrations

import (
	"gorm.io/gorm"
)

func UppercaseTagsNames(db *gorm.DB) error {
	return db.Exec("UPDATE loot_tags.tags SET name = UPPER(name) WHERE id IS NOT NULL").Error
}
