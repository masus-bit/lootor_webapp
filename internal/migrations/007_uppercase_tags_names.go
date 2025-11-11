package migrations

import (
	"gorm.io/gorm"
)

func UppercaseTagsNames(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Table("loot_tags.tags").
			Update("name", gorm.Expr("UPPER(name)"))

		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}
