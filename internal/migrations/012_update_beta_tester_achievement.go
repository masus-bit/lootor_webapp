package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func UpdateBetaTesterAchievement(db *gorm.DB) error {
	return db.Transaction(
		func(tx *gorm.DB) error {

			log.Println("Начинаем миграцию: betaTester для пользователей с level = 0")

			const (
				AchieveBetaTester = "betaTester"
				XPBetaTester      = 100
			)

			releaseDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

			var users []struct {
				Login   string `gorm:"column:login"`
				Created string `gorm:"column:created"`
			}

			if err := tx.Raw(
				`
			SELECT u.login, u.created
			FROM loot.users u
			JOIN loot_achievements.achievements_links al
				ON al.user_login = u.login
			WHERE u.deleted_at IS NULL
			  AND al.deleted_at IS NULL
			  AND al.achievement_code = ?
			  AND COALESCE(al.level, 0) = 0
		`, AchieveBetaTester,
			).Scan(&users).Error; err != nil {
				return fmt.Errorf("ошибка получения пользователей: %w", err)
			}

			log.Printf("Найдено %d пользователей с betaTester.level = 0", len(users))

			parseDate := func(dateStr string) (time.Time, error) {
				dateStr = strings.TrimSpace(dateStr)
				if dateStr == "" {
					return time.Time{}, fmt.Errorf("пустая строка даты")
				}

				formats := []string{
					"2006-01-02 15:04:05.999999",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05Z",
					"2006-01-02",
					time.RFC3339,
					time.RFC3339Nano,
				}

				for _, format := range formats {
					if t, err := time.Parse(format, dateStr); err == nil {
						return t, nil
					}
				}

				return time.Time{}, fmt.Errorf("не удалось распарсить дату: %s", dateStr)
			}

			for i, user := range users {

				if i%50 == 0 {
					log.Printf("Обрабатываем %d/%d: %s", i+1, len(users), user.Login)
				}

				if user.Created == "" {
					continue
				}

				createdAt, err := parseDate(user.Created)
				if err != nil {
					log.Printf("Ошибка парсинга даты для %s: %v", user.Login, err)
					continue
				}

				// Проверка бета-тестера
				if !createdAt.Before(releaseDate) {
					continue
				}

				res := tx.Exec(
					`
				UPDATE loot_achievements.achievements_links
				SET exp = ?,
				    level = 1,
				    current_value_bool = true,
				    updated_at = NOW()
				WHERE user_login = ?
				  AND achievement_code = ?
				  AND deleted_at IS NULL
				  AND COALESCE(level, 0) = 0
			`, XPBetaTester, user.Login, AchieveBetaTester,
				)

				if res.Error != nil {
					log.Printf("Ошибка обновления betaTester для %s: %v", user.Login, res.Error)
					continue
				}

				if res.RowsAffected == 0 {
					continue
				}

				if err := tx.Exec(
					`
				UPDATE loot.users
				SET exp = COALESCE(exp, 0) + ?
				WHERE login = ?
				  AND deleted_at IS NULL
			`, XPBetaTester, user.Login,
				).Error; err != nil {
					log.Printf("Ошибка обновления exp для %s: %v", user.Login, err)
				} else {
					log.Printf("Начислен betaTester для %s: %d XP", user.Login, XPBetaTester)
				}
			}

			log.Println("Миграция betaTester завершена")
			return nil
		},
	)
}
