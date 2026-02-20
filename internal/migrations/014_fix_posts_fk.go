package migrations

import "gorm.io/gorm"

func FixPostsUUIDDefault(db *gorm.DB) error {
	return db.Transaction(
		func(tx *gorm.DB) error {

			if err := tx.Exec(
				`
			CREATE EXTENSION IF NOT EXISTS pgcrypto
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ALTER COLUMN id
			SET DEFAULT gen_random_uuid()
		`,
			).Error; err != nil {
				return err
			}

			return nil
		},
	)
}
