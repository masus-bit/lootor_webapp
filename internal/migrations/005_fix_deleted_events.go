package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func FixDeletedFlagEvents(db *gorm.DB) error {
	if err := db.Exec(`UPDATE loot_events.events SET deleted = false`).Error; err != nil {
		return fmt.Errorf("failed to activate all events: %w", err)
	}

	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = true 
        WHERE 
            (target_user_login IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.users WHERE login = target_user_login AND deleted_at IS NULL
            ))
            OR (target_collection_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.collections WHERE id = target_collection_id AND deleted_at IS NULL
            ))
            OR (target_item_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.collection_items WHERE id = target_item_id AND deleted_at IS NULL
            ))
            OR (target_wish_list_item_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.wish_list_items WHERE id = target_wish_list_item_id AND deleted_at IS NULL
            ))
            OR (target_post_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot_posts.posts WHERE id = target_post_id AND deleted_at IS NULL
            ))
    `).Error; err != nil {
		return fmt.Errorf("failed to mark orphaned events as deleted: %w", err)
	}

	return nil
}
