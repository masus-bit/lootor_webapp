package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func FixDeletedFlagEvents(db *gorm.DB) error {
	var eventsToMarkDeleted []string
	var eventsToMarkActive []string

	var userEventsToDelete []string
	var userEventsToActivate []string

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.users u ON e.target_user_login = u.login").
		Where("e.target_user_login IS NOT NULL AND (u.login IS NULL OR u.deleted_at IS NOT NULL)").
		Pluck("e.id", &userEventsToDelete)

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.users u ON e.target_user_login = u.login").
		Where("e.target_user_login IS NOT NULL AND u.login IS NOT NULL AND u.deleted_at IS NULL").
		Pluck("e.id", &userEventsToActivate)

	eventsToMarkDeleted = append(eventsToMarkDeleted, userEventsToDelete...)
	eventsToMarkActive = append(eventsToMarkActive, userEventsToActivate...)

	var collectionEventsToDelete []string
	var collectionEventsToActivate []string

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.collections c ON e.target_collection_id = c.id::uuid").
		Where("e.target_collection_id IS NOT NULL AND (c.id IS NULL OR c.deleted_at IS NOT NULL)").
		Pluck("e.id", &collectionEventsToDelete)

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.collections c ON e.target_collection_id = c.id::uuid").
		Where("e.target_collection_id IS NOT NULL AND c.id IS NOT NULL AND c.deleted_at IS NULL").
		Pluck("e.id", &collectionEventsToActivate)

	eventsToMarkDeleted = append(eventsToMarkDeleted, collectionEventsToDelete...)
	eventsToMarkActive = append(eventsToMarkActive, collectionEventsToActivate...)

	var itemEventsToDelete []string
	var itemEventsToActivate []string

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.collection_items ci ON e.target_item_id = ci.id::uuid").
		Where("e.target_item_id IS NOT NULL AND (ci.id IS NULL OR ci.deleted_at IS NOT NULL)").
		Pluck("e.id", &itemEventsToDelete)

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.collection_items ci ON e.target_item_id = ci.id::uuid").
		Where("e.target_item_id IS NOT NULL AND ci.id IS NOT NULL AND ci.deleted_at IS NULL").
		Pluck("e.id", &itemEventsToActivate)

	eventsToMarkDeleted = append(eventsToMarkDeleted, itemEventsToDelete...)
	eventsToMarkActive = append(eventsToMarkActive, itemEventsToActivate...)

	var wishListEventsToDelete []string
	var wishListEventsToActivate []string

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.wish_list_items wli ON e.target_wish_list_item_id = wli.id::uuid").
		Where("e.target_wish_list_item_id IS NOT NULL AND (wli.id IS NULL OR wli.deleted_at IS NOT NULL)").
		Pluck("e.id", &wishListEventsToDelete)

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot.wish_list_items wli ON e.target_wish_list_item_id = wli.id::uuid").
		Where("e.target_wish_list_item_id IS NOT NULL AND wli.id IS NOT NULL AND wli.deleted_at IS NULL").
		Pluck("e.id", &wishListEventsToActivate)

	eventsToMarkDeleted = append(eventsToMarkDeleted, wishListEventsToDelete...)
	eventsToMarkActive = append(eventsToMarkActive, wishListEventsToActivate...)

	var postEventsToDelete []string
	var postEventsToActivate []string

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot_posts.posts p ON e.target_post_id = p.id::text").
		Where("e.target_post_id IS NOT NULL AND (p.id IS NULL OR p.deleted_at IS NOT NULL)").
		Pluck("e.id", &postEventsToDelete)

	db.Table("loot_events.events e").
		Select("e.id").
		Joins("LEFT JOIN loot_posts.posts p ON e.target_post_id = p.id::text").
		Where("e.target_post_id IS NOT NULL AND p.id IS NOT NULL AND p.deleted_at IS NULL").
		Pluck("e.id", &postEventsToActivate)

	eventsToMarkDeleted = append(eventsToMarkDeleted, postEventsToDelete...)
	eventsToMarkActive = append(eventsToMarkActive, postEventsToActivate...)

	uniqueEventsToDelete := make(map[string]bool)
	for _, id := range eventsToMarkDeleted {
		uniqueEventsToDelete[id] = true
	}

	uniqueEventsToActivate := make(map[string]bool)
	for _, id := range eventsToMarkActive {
		uniqueEventsToActivate[id] = true
	}

	for id := range uniqueEventsToDelete {
		delete(uniqueEventsToActivate, id)
	}

	if len(uniqueEventsToDelete) > 0 {
		var eventIDs []string
		for id := range uniqueEventsToDelete {
			eventIDs = append(eventIDs, id)
		}

		if err := db.Table("loot_events.events").
			Where("id IN (?)", eventIDs).
			Update("deleted", true).Error; err != nil {
			return fmt.Errorf("failed to mark events as deleted: %w", err)
		}
	}

	if len(uniqueEventsToActivate) > 0 {
		var eventIDs []string
		for id := range uniqueEventsToActivate {
			eventIDs = append(eventIDs, id)
		}

		if err := db.Table("loot_events.events").
			Where("id IN (?)", eventIDs).
			Update("deleted", false).Error; err != nil {
			return fmt.Errorf("failed to mark events as active: %w", err)
		}
	}

	return nil
}
