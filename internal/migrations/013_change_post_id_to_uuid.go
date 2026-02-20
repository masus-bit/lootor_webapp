package migrations

import "gorm.io/gorm"

func ChangePostIdToUUIDWithNumberID(db *gorm.DB) error {
	return db.Transaction(
		func(tx *gorm.DB) error {

			// ------------------------------------------------
			// uuid extension
			// ------------------------------------------------
			if err := tx.Exec(
				`
			CREATE EXTENSION IF NOT EXISTS pgcrypto
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 1. add uuid column to posts
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ADD COLUMN new_uuid uuid
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			UPDATE loot_posts.posts
			SET new_uuid = gen_random_uuid()
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ALTER COLUMN new_uuid SET NOT NULL
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 2. number_id = old bigint id
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ADD COLUMN number_id BIGINT
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			UPDATE loot_posts.posts
			SET number_id = id
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 3. reactions temp uuid fk
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.post_reactions
			ADD COLUMN new_post_uuid uuid
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			UPDATE loot_posts.post_reactions pr
			SET new_post_uuid = p.new_uuid
			FROM loot_posts.posts p
			WHERE pr.post_id = p.id
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 4. drop BOTH FKs
			// ------------------------------------------------
			tx.Exec(`ALTER TABLE loot_posts.post_reactions DROP CONSTRAINT IF EXISTS fk_post_reactions_post`)
			tx.Exec(`ALTER TABLE loot_posts.post_reactions DROP CONSTRAINT IF EXISTS fk_posts_reactions`)

			// ------------------------------------------------
			// 5. drop unique index using old post_id
			// ------------------------------------------------
			tx.Exec(`DROP INDEX IF EXISTS loot_posts.idx_user_post`)

			// ------------------------------------------------
			// 6. drop old bigint columns
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.post_reactions
			DROP COLUMN post_id
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			DROP CONSTRAINT posts_pkey
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			DROP COLUMN id
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 7. rename uuid columns
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			RENAME COLUMN new_uuid TO id
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.post_reactions
			RENAME COLUMN new_post_uuid TO post_id
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 8. restore PK
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ADD PRIMARY KEY (id)
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 9. restore FK
			// ------------------------------------------------
			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.post_reactions
			ALTER COLUMN post_id SET NOT NULL
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.post_reactions
			ADD CONSTRAINT fk_post_reactions_post
			FOREIGN KEY (post_id)
			REFERENCES loot_posts.posts(id)
			ON DELETE CASCADE
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 10. recreate unique index
			// ------------------------------------------------
			if err := tx.Exec(
				`
			CREATE UNIQUE INDEX idx_user_post
			ON loot_posts.post_reactions(post_id, user_login)
		`,
			).Error; err != nil {
				return err
			}

			// ------------------------------------------------
			// 11. AUTO INCREMENT number_id
			// ------------------------------------------------
			if err := tx.Exec(
				`
			CREATE SEQUENCE IF NOT EXISTS loot_posts.posts_number_id_seq
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ALTER COLUMN number_id
			SET DEFAULT nextval('loot_posts.posts_number_id_seq')
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			SELECT setval(
				'loot_posts.posts_number_id_seq',
				(SELECT COALESCE(MAX(number_id),0) FROM loot_posts.posts)
			)
		`,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`
			ALTER TABLE loot_posts.posts
			ALTER COLUMN number_id SET NOT NULL
		`,
			).Error; err != nil {
				return err
			}

			return nil
		},
	)
}
