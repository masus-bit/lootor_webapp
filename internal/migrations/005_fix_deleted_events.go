package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func FixDeletedFlagEvents(db *gorm.DB) error {
	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = false 
        WHERE target_user_login IS NOT NULL 
        AND EXISTS (
            SELECT 1 FROM loot.users 
            WHERE login = target_user_login 
            AND deleted_at IS NULL
        )
    `).Error; err != nil {
		return fmt.Errorf("failed to activate user events: %w", err)
	}

	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = false 
        WHERE target_collection_id IS NOT NULL 
        AND EXISTS (
            SELECT 1 FROM loot.collections 
            WHERE id::uuid = target_collection_id 
            AND deleted_at IS NULL
        )
    `).Error; err != nil {
		return fmt.Errorf("failed to activate collection events: %w", err)
	}

	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = false 
        WHERE target_item_id IS NOT NULL 
        AND EXISTS (
            SELECT 1 FROM loot.collection_items 
            WHERE id::uuid = target_item_id 
            AND deleted_at IS NULL
        )
    `).Error; err != nil {
		return fmt.Errorf("failed to activate item events: %w", err)
	}

	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = false 
        WHERE target_wish_list_item_id IS NOT NULL 
        AND EXISTS (
            SELECT 1 FROM loot.wish_list_items 
            WHERE id::uuid = target_wish_list_item_id 
            AND deleted_at IS NULL
        )
    `).Error; err != nil {
		return fmt.Errorf("failed to activate wishlist events: %w", err)
	}

	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = false 
        WHERE target_post_id IS NOT NULL 
        AND EXISTS (
            SELECT 1 FROM loot_posts.posts 
            WHERE id::text = target_post_id 
            AND deleted_at IS NULL
        )
    `).Error; err != nil {
		return fmt.Errorf("failed to activate post events: %w", err)
	}

	if err := db.Exec(`
        UPDATE loot_events.events 
        SET deleted = true 
        WHERE (
            (target_user_login IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.users 
                WHERE login = target_user_login 
                AND deleted_at IS NULL
            ))
            OR (target_collection_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.collections 
                WHERE id::uuid = target_collection_id 
                AND deleted_at IS NULL
            ))
            OR (target_item_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.collection_items 
                WHERE id::uuid = target_item_id 
                AND deleted_at IS NULL
            ))
            OR (target_wish_list_item_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot.wish_list_items 
                WHERE id::uuid = target_wish_list_item_id 
                AND deleted_at IS NULL
            ))
            OR (target_post_id IS NOT NULL AND NOT EXISTS (
                SELECT 1 FROM loot_posts.posts 
                WHERE id::text = target_post_id 
                AND deleted_at IS NULL
            ))
        )
    `).Error; err != nil {
		return fmt.Errorf("failed to mark orphaned events as deleted: %w", err)
	}

	return nil
}
