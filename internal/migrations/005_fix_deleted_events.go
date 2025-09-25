package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func FixDeletedFlagEvents(db *gorm.DB) error {
	eventsToMarkDeleted := make(map[string]bool)
	eventsToMarkActive := make(map[string]bool)

	var userEvents []struct {
		ID     string
		Exists bool
	}
	db.Raw(`
        SELECT 
            e.id,
            CASE WHEN u.login IS NOT NULL AND u.deleted_at IS NULL THEN true ELSE false END as exists
        FROM loot_events.events e
        LEFT JOIN loot.users u ON e.target_user_login = u.login
        WHERE e.target_user_login IS NOT NULL
    `).Scan(&userEvents)

	for _, event := range userEvents {
		if event.Exists {
			eventsToMarkActive[event.ID] = true
		} else {
			eventsToMarkDeleted[event.ID] = true
		}
	}

	var collectionEvents []struct {
		ID     string
		Exists bool
	}
	db.Raw(`
        SELECT 
            e.id,
            CASE WHEN c.id IS NOT NULL AND c.deleted_at IS NULL THEN true ELSE false END as exists
        FROM loot_events.events e
        LEFT JOIN loot.collections c ON e.target_collection_id = c.id::uuid
        WHERE e.target_collection_id IS NOT NULL
    `).Scan(&collectionEvents)

	for _, event := range collectionEvents {
		if event.Exists {
			eventsToMarkActive[event.ID] = true
		} else {
			eventsToMarkDeleted[event.ID] = true
		}
	}

	var itemEvents []struct {
		ID     string
		Exists bool
	}
	db.Raw(`
        SELECT 
            e.id,
            CASE WHEN ci.id IS NOT NULL AND ci.deleted_at IS NULL THEN true ELSE false END as exists
        FROM loot_events.events e
        LEFT JOIN loot.collection_items ci ON e.target_item_id = ci.id::uuid
        WHERE e.target_item_id IS NOT NULL
    `).Scan(&itemEvents)

	for _, event := range itemEvents {
		if event.Exists {
			eventsToMarkActive[event.ID] = true
		} else {
			eventsToMarkDeleted[event.ID] = true
		}
	}

	var wishListEvents []struct {
		ID     string
		Exists bool
	}
	db.Raw(`
        SELECT 
            e.id,
            CASE WHEN wli.id IS NOT NULL AND wli.deleted_at IS NULL THEN true ELSE false END as exists
        FROM loot_events.events e
        LEFT JOIN loot.wish_list_items wli ON e.target_wish_list_item_id = wli.id::uuid
        WHERE e.target_wish_list_item_id IS NOT NULL
    `).Scan(&wishListEvents)

	for _, event := range wishListEvents {
		if event.Exists {
			eventsToMarkActive[event.ID] = true
		} else {
			eventsToMarkDeleted[event.ID] = true
		}
	}

	var postEvents []struct {
		ID     string
		Exists bool
	}
	db.Raw(`
        SELECT 
            e.id,
            CASE WHEN p.id IS NOT NULL AND p.deleted_at IS NULL THEN true ELSE false END as exists
        FROM loot_events.events e
        LEFT JOIN loot_posts.posts p ON e.target_post_id = p.id::text
        WHERE e.target_post_id IS NOT NULL
    `).Scan(&postEvents)

	for _, event := range postEvents {
		if event.Exists {
			eventsToMarkActive[event.ID] = true
		} else {
			eventsToMarkDeleted[event.ID] = true
		}
	}

	var eventsWithoutTargets []string
	db.Table("loot_events.events").
		Where("target_user_login IS NULL AND target_collection_id IS NULL AND target_item_id IS NULL AND target_wish_list_item_id IS NULL AND target_post_id IS NULL").
		Pluck("id", &eventsWithoutTargets)

	for _, id := range eventsWithoutTargets {
		eventsToMarkActive[id] = true
	}

	for id := range eventsToMarkDeleted {
		delete(eventsToMarkActive, id)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if len(eventsToMarkDeleted) > 0 {
			eventIDs := make([]string, 0, len(eventsToMarkDeleted))
			for id := range eventsToMarkDeleted {
				eventIDs = append(eventIDs, id)
			}

			if err := tx.Table("loot_events.events").
				Where("id IN (?)", eventIDs).
				Update("deleted", true).Error; err != nil {
				return fmt.Errorf("failed to mark events as deleted: %w", err)
			}
		}

		if len(eventsToMarkActive) > 0 {
			eventIDs := make([]string, 0, len(eventsToMarkActive))
			for id := range eventsToMarkActive {
				eventIDs = append(eventIDs, id)
			}

			if err := tx.Table("loot_events.events").
				Where("id IN (?)", eventIDs).
				Update("deleted", false).Error; err != nil {
				return fmt.Errorf("failed to mark events as active: %w", err)
			}
		}

		return nil
	})
}
