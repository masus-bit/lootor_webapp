package migrations

import (
	"gorm.io/gorm"
)

func EntityIDToString(db *gorm.DB) error {
	return db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Exec(
				`
			ALTER TABLE loot_tags.tag_links 
			ADD COLUMN entity_id_str VARCHAR(255)
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			UPDATE loot_tags.tag_links 
			SET entity_id_str = entity_id::text
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_tags.tag_links 
			DROP COLUMN entity_id
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_tags.tag_links 
			RENAME COLUMN entity_id_str TO entity_id
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_tags.tag_links 
			ALTER COLUMN entity_id SET NOT NULL
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_tags.tag_links 
			ADD PRIMARY KEY (tag_id, entity_type, entity_id)
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			CREATE INDEX idx_tag_links_entity ON loot_tags.tag_links(entity_type, entity_id)
		`,
			).Error; err != nil {
				return err
			}

			return nil
		},
	)
}
