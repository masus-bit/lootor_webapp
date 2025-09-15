package migrations

import (
	"gorm.io/gorm"
	"lootor/internal/pkg/utils"
	"strconv"
)

func ChangePostIdToUint64AutoIncrement(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var columnType string
		err := tx.Raw(`
			SELECT data_type 
			FROM information_schema.columns 
			WHERE table_schema = 'loot_posts' 
			AND table_name = 'post_reactions' 
			AND column_name = 'post_id'
		`).Scan(&columnType).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			CREATE SEQUENCE IF NOT EXISTS lootor.loot_posts.posts_id_seq
			START WITH 1
			INCREMENT BY 1
			NO MINVALUE
			NO MAXVALUE
			CACHE 1
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.posts 
			ADD COLUMN new_id BIGINT NOT NULL DEFAULT nextval('lootor.loot_posts.posts_id_seq')
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.post_reactions 
			ADD COLUMN new_post_id BIGINT
		`).Error
		if err != nil {
			return err
		}

		if columnType == "uuid" {
			err = tx.Exec(`
				UPDATE lootor.loot_posts.post_reactions pr
				SET new_post_id = p.new_id
				FROM lootor.loot_posts.posts p
				WHERE pr.post_id::text = p.id::text
			`).Error
		} else {
			err = tx.Exec(`
				UPDATE lootor.loot_posts.post_reactions pr
				SET new_post_id = p.new_id
				FROM lootor.loot_posts.posts p
				WHERE pr.post_id = p.id::text
			`).Error
		}
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.post_reactions 
			DROP COLUMN post_id
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.posts 
			DROP COLUMN id
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.post_reactions 
			RENAME COLUMN new_post_id TO post_id
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.posts 
			RENAME COLUMN new_id TO id
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.posts 
			ADD PRIMARY KEY (id)
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.post_reactions 
			ALTER COLUMN post_id SET NOT NULL
		`).Error
		if err != nil {
			return err
		}

		err = tx.Exec(`
			ALTER TABLE lootor.loot_posts.post_reactions 
			ADD CONSTRAINT fk_post_reactions_post 
			FOREIGN KEY (post_id) 
			REFERENCES lootor.loot_posts.posts(id)
			ON DELETE CASCADE
		`).Error
		if err != nil {
			return err
		}

		var allExistingPosts []struct {
			ID    int64  `gorm:"column:id"`
			Title string `gorm:"column:title"`
		}

		if err = tx.Table("lootor.loot_posts.posts").Order("id").Find(&allExistingPosts).Error; err != nil {
			return err
		}

		for i, post := range allExistingPosts {
			translit := utils.Slugify(post.Title) + "_" + strconv.Itoa(i+1)

			err = tx.Exec(`
				UPDATE lootor.loot_posts.posts 
				SET translit = $1 
				WHERE id = $2
			`, translit, post.ID).Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}
