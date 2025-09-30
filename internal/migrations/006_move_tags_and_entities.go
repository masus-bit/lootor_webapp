package migrations

import (
	"fmt"
	"gorm.io/gorm"

	"lootor/internal/pkg/utils"
)

func MoveTagsAndEntities(db *gorm.DB) error {
	sourceSchema := "loot"
	targetSchema := "loot_tags"

	createSchemaSQL := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, targetSchema)
	if err := db.Exec(createSchemaSQL).Error; err != nil {
		return fmt.Errorf("failed to create target schema '%s': %w", targetSchema, err)
	}

	moveTagsSQL := fmt.Sprintf(`INSERT INTO %s.tags (id, name, slug, is_primary, created_at, type, updated_at) 
		SELECT id, name, name as slug, true as is_primary, created_at, 'collection' as type, updated_at 
		FROM %s.tags;`, targetSchema, sourceSchema)

	if err := db.Exec(moveTagsSQL).Error; err != nil {
		return fmt.Errorf("failed to move tags from %s.tags to %s.tags: %w", sourceSchema, targetSchema, err)
	}

	type Tag struct {
		ID   string
		Name string
	}

	var tags []Tag
	if err := db.Table(fmt.Sprintf("%s.tags", targetSchema)).Where("type = ?", "collection").Select("id, name").Find(&tags).Error; err != nil {
		return fmt.Errorf("failed to select tags from %s.tags: %w", targetSchema, err)
	}

	for _, tag := range tags {
		slug := utils.Slugify(tag.Name)
		uniqueSlug := generateUniqueSlug(db, targetSchema, slug, tag.ID)

		updateSQL := fmt.Sprintf(`UPDATE %s.tags SET slug = ? WHERE id = ?`, targetSchema)
		if err := db.Exec(updateSQL, uniqueSlug, tag.ID).Error; err != nil {
			return fmt.Errorf("failed to update slug for tag %s: %w", tag.ID, err)
		}
	}

	moveCollectionLinksSQL := fmt.Sprintf(`INSERT INTO %s.tag_links (tag_id, entity_type, entity_id, created_at) 
		SELECT tags_id as tag_id, 'collection' as entity_type, collections_id as entity_id, NOW() as created_at 
		FROM %s.tags_collections_collections;`, targetSchema, sourceSchema)

	if err := db.Exec(moveCollectionLinksSQL).Error; err != nil {
		return fmt.Errorf("failed to move collection links to %s.tag_links: %w", targetSchema, err)
	}

	moveEntitiesSQL := fmt.Sprintf(`INSERT INTO %s.tags (id, name, slug, author, is_primary, type, created_at, updated_at) 
		SELECT id, name, id::text as slug, author, true as is_primary, 'collectionItem' as type, created_at, updated_at 
		FROM %s.entities;`, targetSchema, sourceSchema)

	if err := db.Exec(moveEntitiesSQL).Error; err != nil {
		return fmt.Errorf("failed to move entities from %s.entities to %s.tags: %w", sourceSchema, targetSchema, err)
	}

	type Entity struct {
		ID              string
		Name            string
		Transliteration *string
	}

	var entities []Entity
	if err := db.Table(fmt.Sprintf("%s.entities", sourceSchema)).Select("id, name, transliteration").Find(&entities).Error; err != nil {
		return fmt.Errorf("failed to select entities from %s.entities: %w", sourceSchema, err)
	}

	for _, entity := range entities {
		slugSource := entity.Name
		if entity.Transliteration != nil && *entity.Transliteration != "" {
			slugSource = *entity.Transliteration
		}

		slug := utils.Slugify(slugSource)
		uniqueSlug := generateUniqueSlug(db, targetSchema, slug, entity.ID)

		updateSQL := fmt.Sprintf(`UPDATE %s.tags SET slug = ? WHERE id = ?`, targetSchema)
		if err := db.Exec(updateSQL, uniqueSlug, entity.ID).Error; err != nil {
			return fmt.Errorf("failed to update slug for entity %s: %w", entity.ID, err)
		}
	}

	moveEntityLinksSQL := fmt.Sprintf(`INSERT INTO %s.tag_links (tag_id, entity_type, entity_id, created_at) 
		SELECT entities_id as tag_id, 'collection_item' as entity_type, collection_items_id as entity_id, NOW() as created_at 
		FROM %s.entities_collection_item_collection_items;`, targetSchema, sourceSchema)

	if err := db.Exec(moveEntityLinksSQL).Error; err != nil {
		return fmt.Errorf("failed to move entity links to %s.tag_links: %w", targetSchema, err)
	}

	return nil
}

func generateUniqueSlug(db *gorm.DB, schema string, baseSlug string, currentID string) string {
	slug := baseSlug
	counter := 1

	for {
		var count int64
		checkSQL := fmt.Sprintf(`SELECT COUNT(*) FROM %s.tags WHERE slug = ? AND id != ?`, schema)
		db.Raw(checkSQL, slug, currentID).Scan(&count)

		if count == 0 {
			break
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		if counter > 1000 {
			slug = fmt.Sprintf("%s-%s", baseSlug, currentID[:8])
			break
		}
	}

	return slug
}
