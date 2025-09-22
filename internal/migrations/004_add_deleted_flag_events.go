package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func AddDeletedFlagEvents(db *gorm.DB) error {
	var eventsToDelete []string

	var userEvents []string
	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.users u ON e.target_user_login = u.login AND u.deleted_at IS NULL").
		Where("e.deleted = false AND e.target_user_login IS NOT NULL AND u.login IS NULL").
		Pluck("e.id", &userEvents)
	eventsToDelete = append(eventsToDelete, userEvents...)

	var collectionEvents []string
	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.collections c ON e.target_collection_id = c.id::uuid AND c.deleted_at IS NULL").
		Where("e.deleted = false AND e.target_collection_id IS NOT NULL AND c.id IS NULL").
		Pluck("e.id", &collectionEvents)
	eventsToDelete = append(eventsToDelete, collectionEvents...)

	var itemEvents []string
	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.collection_items ci ON e.target_item_id = ci.id::uuid AND ci.deleted_at IS NULL").
		Where("e.deleted = false AND e.target_item_id IS NOT NULL AND ci.id IS NULL").
		Pluck("e.id", &itemEvents)
	eventsToDelete = append(eventsToDelete, itemEvents...)

	var wishListEvents []string
	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.wish_list_items wli ON e.target_wish_list_item_id = wli.id::uuid AND wli.deleted_at IS NULL").
		Where("e.deleted = false AND e.target_wish_list_item_id IS NOT NULL AND wli.id IS NULL").
		Pluck("e.id", &wishListEvents)
	eventsToDelete = append(eventsToDelete, wishListEvents...)

	var postEvents []string
	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot_posts.posts p ON e.target_post_id = p.id::text AND p.deleted_at IS NULL").
		Where("e.deleted = false AND e.target_post_id IS NOT NULL AND p.id IS NULL").
		Pluck("e.id", &postEvents)
	eventsToDelete = append(eventsToDelete, postEvents...)

	uniqueEvents := make(map[string]bool)
	for _, id := range eventsToDelete {
		uniqueEvents[id] = true
	}

	if len(uniqueEvents) > 0 {
		var eventIDs []string
		for id := range uniqueEvents {
			eventIDs = append(eventIDs, id)
		}

		if err := db.Table("loot_events.events").
			Where("id IN (?)", eventIDs).
			Update("deleted", true).Error; err != nil {
			return fmt.Errorf("failed to update events: %w", err)
		}
	}

	return nil
}
