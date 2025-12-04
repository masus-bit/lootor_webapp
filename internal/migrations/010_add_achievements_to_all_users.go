package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

func AddAchievementsToAllUsers(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		log.Println("Начинаем миграцию: заполнение справочника достижений и добавление их пользователям...")

		// 1. Сначала заполняем справочник achievements если он пустой
		var existingCount int64
		if err := tx.Raw(`SELECT COUNT(*) FROM loot_achievements.achievements WHERE deleted_at IS NULL`).
			Scan(&existingCount).Error; err != nil {
			return fmt.Errorf("ошибка проверки существующих достижений: %w", err)
		}

		if existingCount == 0 {
			log.Println("Справочник достижений пуст, заполняем...")

			// Все коды достижений
			achievementCodes := []string{
				"betaTester",
				"betaTesterDonate",
				"firstCollectionCreate",
				"donate",
				"subscribers",
				"yearlyRegister",
				"collectionItemsAdded",
				"photosAdded",
				"tagsCreated",
				"postsCreated",
				"collectionsLikes",
				"collectionItemsLikes",
				"photosLikes",
				"postsReactions",
				"collectionsSum",
				"collectionShipSum",
			}

			// Вставляем достижения
			for _, code := range achievementCodes {
				if err := tx.Exec(`
					INSERT INTO loot_achievements.achievements 
					(id, code, created_at, updated_at)
					VALUES (gen_random_uuid(), ?, NOW(), NOW())
				`, code).Error; err != nil {
					log.Printf("Предупреждение при вставке достижения %s: %v", code, err)
				}
			}

			log.Printf("Попытка добавления %d достижений в справочник завершена", len(achievementCodes))
		} else {
			log.Printf("В справочнике уже есть %d достижений, пропускаем заполнение", existingCount)
		}

		// 2. Проверяем, что таблица achievements заполнена
		var achievements []struct {
			ID   string
			Code string
		}

		if err := tx.Raw(`
            SELECT id::text, code 
            FROM loot_achievements.achievements 
            WHERE deleted_at IS NULL
        `).Scan(&achievements).Error; err != nil {
			return fmt.Errorf("ошибка получения достижений: %w", err)
		}

		if len(achievements) == 0 {
			return fmt.Errorf("таблица achievements пуста после попытки заполнения")
		}

		log.Printf("Найдено %d достижений в справочнике", len(achievements))

		// 3. Получаем количество пользователей
		var userCount int64
		if err := tx.Raw(`SELECT COUNT(*) FROM loot.users WHERE deleted_at IS NULL`).
			Scan(&userCount).Error; err != nil {
			return fmt.Errorf("ошибка подсчета пользователей: %w", err)
		}

		log.Printf("Найдено %d пользователей", userCount)

		// 4. Рассчитываем общее количество записей
		totalRecords := int64(len(achievements)) * userCount
		log.Printf("Будет добавлено до %d записей (пользователи × достижения)", totalRecords)

		// 5. Проверяем, есть ли уже записи в achievements_links
		var existingLinksCount int64
		if err := tx.Raw(`
			SELECT COUNT(*) 
			FROM loot_achievements.achievements_links 
			WHERE deleted_at IS NULL
		`).Scan(&existingLinksCount).Error; err != nil {
			log.Printf("Таблица achievements_links возможно не существует или пуста: %v", err)
			existingLinksCount = 0
		}

		if existingLinksCount > 0 {
			log.Printf("В таблице achievements_links уже есть %d записей", existingLinksCount)
			log.Println("Проверяем, нужно ли добавлять новые...")

			// Удаляем существующие записи чтобы избежать дублей
			// (если ты хочешь сохранить существующие данные, удали эту часть)
			log.Println("Очищаем существующие записи для чистой миграции...")
			if err := tx.Exec(`DELETE FROM loot_achievements.achievements_links WHERE deleted_at IS NULL`).Error; err != nil {
				log.Printf("Не удалось очистить таблицу: %v", err)
			}
		}

		// 6. Выполняем вставку БЕЗ ON CONFLICT
		log.Println("Начинаем вставку записей для пользователей...")

		result := tx.Exec(`
            INSERT INTO loot_achievements.achievements_links (
                id,
                achievement_id,
                achievement_code,
                user_login,
                exp,
                level,
                current_value_int,
                current_value_bool,
                created_at,
                updated_at
            )
            SELECT
                gen_random_uuid(),
                a.id,
                a.code,
                u.login,
                0, -- начальный опыт
                0, -- начальный уровень
                0, -- начальное значение
                FALSE, -- начальное булево значение
                NOW(),
                NOW()
            FROM loot.users u
            CROSS JOIN loot_achievements.achievements a
            WHERE u.deleted_at IS NULL
              AND a.deleted_at IS NULL
        `)

		if result.Error != nil {
			return fmt.Errorf("ошибка вставки: %w", result.Error)
		}

		log.Printf("Успешно добавлено %d записей", result.RowsAffected)

		// 7. Создаем UNIQUE constraint (правильный синтаксис)
		log.Println("Создаем UNIQUE constraint для предотвращения дублей в будущем...")

		// Сначала проверяем, существует ли уже constraint
		var constraintExists bool
		tx.Raw(`
			SELECT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_schema = 'loot_achievements' 
				AND table_name = 'achievements_links' 
				AND constraint_name = 'unique_user_achievement'
			)
		`).Scan(&constraintExists)

		if !constraintExists {
			if err := tx.Exec(`
				ALTER TABLE loot_achievements.achievements_links 
				ADD CONSTRAINT unique_user_achievement 
				UNIQUE (user_login, achievement_code)
			`).Error; err != nil {
				log.Printf("Не удалось создать constraint: %v", err)
				// Не прерываем выполнение, constraint может уже существовать под другим именем
			} else {
				log.Println("UNIQUE constraint успешно создан")
			}
		} else {
			log.Println("UNIQUE constraint уже существует")
		}

		// 8. Удаляем дубликаты если они есть (перед созданием constraint это важно)
		log.Println("Проверяем и удаляем возможные дубликаты...")

		// Создаем временную таблицу для удаления дублей
		if err := tx.Exec(`
			CREATE TEMP TABLE temp_achievement_links AS
			SELECT DISTINCT ON (user_login, achievement_code) *
			FROM loot_achievements.achievements_links
			WHERE deleted_at IS NULL
			ORDER BY user_login, achievement_code, created_at
		`).Error; err != nil {
			log.Printf("Не удалось создать временную таблицу: %v", err)
		} else {
			// Удаляем все записи и вставляем уникальные
			if err := tx.Exec(`DELETE FROM loot_achievements.achievements_links WHERE deleted_at IS NULL`).Error; err != nil {
				log.Printf("Не удалось удалить старые записи: %v", err)
			} else {
				if err := tx.Exec(`
					INSERT INTO loot_achievements.achievements_links 
					SELECT * FROM temp_achievement_links
				`).Error; err != nil {
					log.Printf("Не удалось вставить уникальные записи: %v", err)
				}
			}

			// Удаляем временную таблицу
			tx.Exec(`DROP TABLE IF EXISTS temp_achievement_links`)
		}

		// 9. Создаем индексы для ускорения
		log.Println("Создаем индексы для оптимизации...")

		indexes := []string{
			`CREATE INDEX idx_achievements_links_lookup 
			 ON loot_achievements.achievements_links (user_login, achievement_code) 
			 WHERE deleted_at IS NULL`,

			`CREATE INDEX idx_achievements_links_code 
			 ON loot_achievements.achievements_links (achievement_code) 
			 WHERE deleted_at IS NULL`,

			`CREATE INDEX idx_achievements_links_value 
			 ON loot_achievements.achievements_links (achievement_code, current_value_int DESC) 
			 WHERE deleted_at IS NULL`,

			`CREATE INDEX idx_achievements_links_level 
			 ON loot_achievements.achievements_links (level) 
			 WHERE deleted_at IS NULL AND level > 0`,

			`CREATE INDEX idx_achievements_links_updated 
			 ON loot_achievements.achievements_links (updated_at DESC) 
			 WHERE deleted_at IS NULL`,
		}

		for _, indexSQL := range indexes {
			if err := tx.Exec(indexSQL).Error; err != nil {
				// Игнорируем ошибки создания индексов (они могут уже существовать)
				log.Printf("Предупреждение при создании индекса: %v", err)
			}
		}

		// 10. Также создаем UNIQUE constraint для таблицы achievements
		log.Println("Создаем UNIQUE constraint для таблицы achievements...")

		var achievementsConstraintExists bool
		tx.Raw(`
			SELECT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_schema = 'loot_achievements' 
				AND table_name = 'achievements' 
				AND constraint_name = 'achievements_code_unique'
			)
		`).Scan(&achievementsConstraintExists)

		if !achievementsConstraintExists {
			if err := tx.Exec(`
				ALTER TABLE loot_achievements.achievements 
				ADD CONSTRAINT achievements_code_unique 
				UNIQUE (code)
			`).Error; err != nil {
				log.Printf("Не удалось создать constraint для achievements: %v", err)
			} else {
				log.Println("UNIQUE constraint для achievements успешно создан")
			}
		}

		log.Println("Миграция успешно завершена!")
		return nil
	})
}
