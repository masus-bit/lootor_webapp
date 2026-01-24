package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func DeleteTagsDuplicates(db *gorm.DB) error {
	return db.Transaction(
		func(tx *gorm.DB) error {

			var duplicateCount int64
			if err := tx.Raw(
				`
			SELECT COUNT(*) 
			FROM (
				SELECT name 
				FROM loot_tags.tags 
				WHERE deleted_at IS NULL 
				GROUP BY name 
				HAVING COUNT(*) > 1
			) AS duplicates
		`,
			).Scan(&duplicateCount).Error; err != nil {
				return fmt.Errorf("ошибка проверки дубликатов: %w", err)
			}

			if duplicateCount > 0 {

				if err := deleteDuplicateTags(tx); err != nil {
					return err
				}
			}

			if err := tx.Exec(
				`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name 
			ON loot_tags.tags (name)
			WHERE deleted_at IS NULL
		`,
			).Error; err != nil {
				return fmt.Errorf("ошибка создания уникального индекса: %w", err)
			}

			return nil
		},
	)
}

func deleteDuplicateTags(tx *gorm.DB) error {

	if err := tx.Exec(
		`
		WITH ranked_tags AS (
			SELECT 
				id,
				name,
				created_at,
				ROW_NUMBER() OVER (PARTITION BY name ORDER BY created_at, id) as rn
			FROM loot_tags.tags
			WHERE deleted_at IS NULL
		)
		UPDATE loot_tags.tags
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id IN (
			SELECT id 
			FROM ranked_tags 
			WHERE rn > 1
		)
	`,
	).Error; err != nil {
		return fmt.Errorf("ошибка удаления дубликатов: %w", err)
	}

	return nil
}
