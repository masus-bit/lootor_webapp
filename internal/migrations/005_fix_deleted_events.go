package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func MarkEventsDeletedForDeletedTargets(db *gorm.DB) error {
	return db.Transaction(
		func(tx *gorm.DB) error {
			queries := []string{
				`UPDATE loot_events.events 
             SET deleted = true 
             WHERE target_post_id IN (
                 SELECT target_post_id 
                 FROM loot_events.events 
                 WHERE action = 'delete' 
                 AND event_target_type = 'post'
                 AND target_post_id IS NOT NULL
             )`,

				`UPDATE loot_events.events 
             SET deleted = true 
             WHERE target_item_id IN (
                 SELECT target_item_id 
                 FROM loot_events.events 
                 WHERE action = 'delete' 
                 AND event_target_type = 'collectionItem'
                 AND target_item_id IS NOT NULL
             )`,

				`UPDATE loot_events.events 
             SET deleted = true 
             WHERE target_collection_id IN (
                 SELECT target_collection_id 
                 FROM loot_events.events 
                 WHERE action = 'delete' 
                 AND event_target_type = 'collection'
                 AND target_collection_id IS NOT NULL
             )`,

				`UPDATE loot_events.events 
             SET deleted = true 
             WHERE target_wish_list_item_id IN (
                 SELECT target_wish_list_item_id 
                 FROM loot_events.events 
                 WHERE action = 'delete' 
                 AND event_target_type = 'wishListItem'
                 AND target_wish_list_item_id IS NOT NULL
             )`,
			}

			for i, query := range queries {
				if err := tx.Exec(query).Error; err != nil {
					return fmt.Errorf("failed query %d: %w", i, err)
				}
			}

			return nil
		},
	)
}
