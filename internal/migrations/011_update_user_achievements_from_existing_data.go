package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

func UpdateUserAchievementsFromExistingData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		log.Println("Начинаем миграцию: обновление достижений пользователей из существующих данных...")

		// Константы
		const (
			// Коды достижений
			AchieveBetaTester            = "betaTester"
			AchieveFirstCollectionCreate = "firstCollectionCreate"
			AchieveSubscribers           = "subscribers"
			AchieveYearlyRegister        = "yearlyRegister"
			AchieveCollectionItemsAdded  = "collectionItemsAdded"
			AchievePhotosAdded           = "photosAdded"
			AchieveTagsCreated           = "tagsCreated"
			AchievePostsCreated          = "postsCreated"
			AchieveCollectionsLikes      = "collectionsLikes"
			AchieveCollectionItemsLikes  = "collectionItemsLikes"
			AchievePhotosLikes           = "photosLikes"
			AchievePostsReactions        = "postsReactions"
			AchieveCollectionsSum        = "collectionsSum"
			AchieveCollectionShipSum     = "collectionShipSum"

			// Опыт за достижения
			XPBetaTester                 = 100
			XPFirstCollectionCreate      = 25
			XPSubscribersLevel1          = 10
			XPSubscribersLevel2          = 200
			XPSubscribersLevel3          = 500
			XPYearlyRegister             = 1000
			XPCollectionItemsAddLevel1   = 20
			XPCollectionItemsAddLevel2   = 250
			XPCollectionItemsAddLevel3   = 3000
			XPPhotosAddLevel1            = 20
			XPPhotosAddLevel2            = 250
			XPPhotosAddLevel3            = 3000
			XPTagsAddLevel1              = 50
			XPTagsAddLevel2              = 400
			XPTagsAddLevel3              = 1500
			XPPostsCreateLevel1          = 25
			XPPostsCreateLevel2          = 300
			XPPostsCreateLevel3          = 1000
			XPCollectionsLikesLevel1     = 20
			XPCollectionsLikesLevel2     = 150
			XPCollectionsLikesLevel3     = 500
			XPCollectionItemsLikesLevel1 = 50
			XPCollectionItemsLikesLevel2 = 300
			XPCollectionItemsLikesLevel3 = 1500
			XPPhotosLikesLevel1          = 50
			XPPhotosLikesLevel2          = 300
			XPPhotosLikesLevel3          = 1500
			XPPostsReactionsLevel1       = 50
			XPPostsReactionsLevel2       = 300
			XPPostsReactionsLevel3       = 1500
			XPCollectionsSumLevel1       = 50
			XPCollectionsSumLevel2       = 300
			XPCollectionsSumLevel3       = 1500
			XPCollectionsShipSumLevel1   = 50
			XPCollectionsShipSumLevel2   = 250
			XPCollectionsShipSumLevel3   = 1000

			// Пороги для уровней
			SubscribersLevel1          = 1
			SubscribersLevel2          = 20
			SubscribersLevel3          = 50
			CollectionItemsAddLevel1   = 10
			CollectionItemsAddLevel2   = 75
			CollectionItemsAddLevel3   = 200
			PhotosAddLevel1            = 10
			PhotosAddLevel2            = 75
			PhotosAddLevel3            = 200
			TagsAddLevel1              = 5
			TagsAddLevel2              = 30
			TagsAddLevel3              = 100
			PostsCreateLevel1          = 1
			PostsCreateLevel2          = 15
			PostsCreateLevel3          = 75
			CollectionsLikesLevel1     = 10
			CollectionsLikesLevel2     = 50
			CollectionsLikesLevel3     = 200
			CollectionItemsLikesLevel1 = 50
			CollectionItemsLikesLevel2 = 250
			CollectionItemsLikesLevel3 = 1000
			PhotosLikesLevel1          = 50
			PhotosLikesLevel2          = 250
			PhotosLikesLevel3          = 1000
			PostsReactionsLevel1       = 50
			PostsReactionsLevel2       = 250
			PostsReactionsLevel3       = 1000
			CollectionsSumLevel1       = 50000
			CollectionsSumLevel2       = 250000
			CollectionsSumLevel3       = 1000000
			CollectionsShipSumLevel1   = 4000
			CollectionsShipSumLevel2   = 20000
			CollectionsShipSumLevel3   = 80000
		)

		// Дата релиза для бета-тестеров (5 декабря 2025 года)
		releaseDate := time.Date(2025, 12, 5, 0, 0, 0, 0, time.UTC)

		// 1. Получаем всех активных пользователей
		log.Println("Получаем список пользователей...")
		var users []struct {
			Login   string `gorm:"column:login"`
			Created string `gorm:"column:created"`
		}

		if err := tx.Table("loot.users").
			Select("login, created").
			Where("deleted_at IS NULL").
			Find(&users).Error; err != nil {
			return fmt.Errorf("ошибка получения пользователей: %w", err)
		}

		log.Printf("Найдено %d пользователей", len(users))

		// Функция для парсинга даты из строки
		parseDate := func(dateStr string) (time.Time, error) {
			dateStr = strings.TrimSpace(dateStr)
			if dateStr == "" {
				return time.Time{}, fmt.Errorf("пустая строка даты")
			}

			// Пробуем разные форматы дат
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

			// Пробуем Unix timestamp
			if unixTime, err := strconv.ParseInt(dateStr, 10, 64); err == nil {
				// Если число меньше 10^10, считаем секундами, иначе миллисекундами
				if unixTime < 10000000000 {
					return time.Unix(unixTime, 0), nil
				}
				return time.Unix(unixTime/1000, (unixTime%1000)*1000000), nil
			}

			return time.Time{}, fmt.Errorf("не удалось распарсить дату: %s", dateStr)
		}

		// 2. Для каждого пользователя собираем статистику и обновляем достижения
		for i, user := range users {
			if i%50 == 0 {
				log.Printf("Обрабатываем пользователя %d из %d: %s", i+1, len(users), user.Login)
			}

			userLogin := user.Login

			// Собираем статистику для пользователя
			var stats struct {
				// Бета-тестер
				IsBetaTester bool

				// Подписчики
				SubscribersCount int64

				// Годы с регистрации
				YearsRegistered int64

				// Коллекции
				HasCollections   bool
				CollectionsCount int64
				CollectionsLikes int64

				// Экземпляры
				CollectionItemsCount int64
				CollectionItemsLikes int64
				CollectionsSum       int64
				CollectionShipSum    int64

				// Фото
				PhotosCount int64
				PhotosLikes int64

				// Теги
				TagsCount int64

				// Посты
				PostsCount     int64
				PostsReactions int64
			}

			// 2.1 Парсим дату регистрации и проверяем бета-тестера
			var createdAt time.Time
			if user.Created != "" {
				parsedDate, err := parseDate(user.Created)
				if err != nil {
					log.Printf("Ошибка парсинга даты для пользователя %s: %v", userLogin, err)
				} else {
					createdAt = parsedDate
					// Проверяем бета-тестера
					stats.IsBetaTester = createdAt.Before(releaseDate)

					// Годы с регистрации
					years := int64(time.Since(createdAt).Hours() / 24 / 365)
					if years > 0 {
						stats.YearsRegistered = years
					}
				}
			}

			// 2.2 Подписчики - используем COALESCE для обработки NULL
			var subscribersCount int64
			if err := tx.Raw(`
				SELECT COALESCE(array_length(subscribers_logins, 1), 0) 
				FROM loot.users 
				WHERE login = ? AND deleted_at IS NULL
			`, userLogin).Scan(&subscribersCount).Error; err != nil {
				log.Printf("Ошибка получения подписчиков для %s: %v", userLogin, err)
			} else {
				stats.SubscribersCount = subscribersCount
			}

			// 2.3 Коллекции
			// Есть ли коллекции
			var hasCollections int64
			if err := tx.Table("loot.collections").
				Where("user_login = ? AND deleted_at IS NULL", userLogin).
				Count(&hasCollections).Error; err != nil {
				log.Printf("Ошибка проверки коллекций для %s: %v", userLogin, err)
			} else {
				stats.HasCollections = hasCollections > 0
			}

			// Количество коллекций
			if err := tx.Table("loot.collections").
				Where("user_login = ? AND deleted_at IS NULL", userLogin).
				Count(&stats.CollectionsCount).Error; err != nil {
				log.Printf("Ошибка подсчета коллекций для %s: %v", userLogin, err)
			}

			// Лайки на коллекциях
			var collectionsLikes int64
			if err := tx.Raw(`
				SELECT COALESCE(SUM(array_length(likes, 1)), 0) 
				FROM loot.collections 
				WHERE user_login = ? AND deleted_at IS NULL
			`, userLogin).Scan(&collectionsLikes).Error; err != nil {
				log.Printf("Ошибка подсчета лайков на коллекциях для %s: %v", userLogin, err)
			} else {
				stats.CollectionsLikes = collectionsLikes
			}

			// 2.4 Экземпляры коллекций
			// Количество экземпляров
			if err := tx.Table("loot.collection_items").
				Where("user_login = ? AND deleted_at IS NULL", userLogin).
				Count(&stats.CollectionItemsCount).Error; err != nil {
				log.Printf("Ошибка подсчета экземпляров для %s: %v", userLogin, err)
			}

			// Лайки на экземплярах
			var itemsLikes int64
			if err := tx.Raw(`
				SELECT COALESCE(SUM(array_length(likes, 1)), 0) 
				FROM loot.collection_items 
				WHERE user_login = ? AND deleted_at IS NULL
			`, userLogin).Scan(&itemsLikes).Error; err != nil {
				log.Printf("Ошибка подсчета лайков на экземплярах для %s: %v", userLogin, err)
			} else {
				stats.CollectionItemsLikes = itemsLikes
			}

			// Сумма покупок и доставок
			var priceStats []struct {
				PurchasePrice int64 `gorm:"column:purchase_price"`
				ShippingCost  int64 `gorm:"column:shipping_cost"`
			}
			if err := tx.Table("loot.collection_items").
				Select("COALESCE(purchase_price, 0) as purchase_price, COALESCE(shipping_cost, 0) as shipping_cost").
				Where("user_login = ? AND deleted_at IS NULL", userLogin).
				Find(&priceStats).Error; err != nil {
				log.Printf("Ошибка получения цен для %s: %v", userLogin, err)
			} else {
				for _, ps := range priceStats {
					stats.CollectionsSum += ps.PurchasePrice
					stats.CollectionShipSum += ps.ShippingCost
				}
			}

			// 2.5 Фото
			// Количество фото
			if err := tx.Table("loot_photos.photos").
				Where("author = ? AND deleted_at IS NULL", userLogin).
				Count(&stats.PhotosCount).Error; err != nil {
				log.Printf("Ошибка подсчета фото для %s: %v", userLogin, err)
			}

			// Лайки на фото
			var photosLikes int64
			if err := tx.Raw(`
				SELECT COALESCE(SUM(array_length(likes, 1)), 0) 
				FROM loot_photos.photos 
				WHERE author = ? AND deleted_at IS NULL
			`, userLogin).Scan(&photosLikes).Error; err != nil {
				log.Printf("Ошибка подсчета лайков на фото для %s: %v", userLogin, err)
			} else {
				stats.PhotosLikes = photosLikes
			}

			// 2.6 Теги
			if err := tx.Table("loot_tags.tags").
				Where("author = ? AND deleted_at IS NULL", userLogin).
				Count(&stats.TagsCount).Error; err != nil {
				log.Printf("Ошибка подсчета тегов для %s: %v", userLogin, err)
			}

			// 2.7 Посты
			// Количество постов
			if err := tx.Table("loot_posts.posts").
				Where("author = ? AND deleted_at IS NULL", userLogin).
				Count(&stats.PostsCount).Error; err != nil {
				log.Printf("Ошибка подсчета постов для %s: %v", userLogin, err)
			}

			// Реакции на посты
			var reactionCount int64
			if err := tx.Raw(`
				SELECT COUNT(*) 
				FROM loot_posts.post_reactions pr
				JOIN loot_posts.posts p ON pr.post_id = p.id
				WHERE p.author = ? AND p.deleted_at IS NULL AND pr.deleted_at IS NULL
			`, userLogin).Scan(&reactionCount).Error; err != nil {
				log.Printf("Ошибка подсчета реакций на посты для %s: %v", userLogin, err)
			} else {
				stats.PostsReactions = reactionCount
			}

			// 3. Обновляем достижения пользователя в таблице achievements_links
			updateAchievement := func(code string, exp, level, currentValue int64, currentBool bool) error {
				return tx.Exec(`
					UPDATE loot_achievements.achievements_links 
					SET exp = ?, 
						level = ?, 
						current_value_int = ?,
						current_value_bool = ?,
						updated_at = NOW()
					WHERE user_login = ? 
						AND achievement_code = ?
						AND deleted_at IS NULL
				`, exp, level, currentValue, currentBool, userLogin, code).Error
			}

			// Переменная для подсчета общего опыта пользователя
			totalUserExp := int64(0)

			// 3.1 Бета-тестер
			if stats.IsBetaTester {
				if err := updateAchievement(AchieveBetaTester, XPBetaTester, 1, 0, true); err != nil {
					log.Printf("Ошибка обновления betaTester для %s: %v", userLogin, err)
				} else {
					log.Printf("Начислено бета-тестер для %s: %d XP", userLogin, XPBetaTester)
					totalUserExp += XPBetaTester
				}
			}

			// 3.2 Первая коллекция
			if stats.HasCollections {
				if err := updateAchievement(AchieveFirstCollectionCreate, XPFirstCollectionCreate, 1, 0, true); err != nil {
					log.Printf("Ошибка обновления firstCollectionCreate для %s: %v", userLogin, err)
				} else {
					log.Printf("Начислено первая коллекция для %s: %d XP", userLogin, XPFirstCollectionCreate)
					totalUserExp += XPFirstCollectionCreate
				}
			}

			// 3.3 Подписчики
			if stats.SubscribersCount > 0 {
				var subsExp, subsLevel int64
				switch {
				case stats.SubscribersCount >= SubscribersLevel3:
					subsExp = XPSubscribersLevel3
					subsLevel = 3
				case stats.SubscribersCount >= SubscribersLevel2:
					subsExp = XPSubscribersLevel2
					subsLevel = 2
				case stats.SubscribersCount >= SubscribersLevel1:
					subsExp = XPSubscribersLevel1
					subsLevel = 1
				}

				if subsExp > 0 {
					if err := updateAchievement(AchieveSubscribers, subsExp, subsLevel, stats.SubscribersCount, false); err != nil {
						log.Printf("Ошибка обновления subscribers для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено подписчики для %s: %d XP (уровень %d)", userLogin, subsExp, subsLevel)
						totalUserExp += subsExp
					}
				}
			}

			// 3.4 Годы с регистрации
			if stats.YearsRegistered >= 1 {
				yearlyExp := stats.YearsRegistered * XPYearlyRegister
				if err := updateAchievement(AchieveYearlyRegister, yearlyExp, 1, stats.YearsRegistered, false); err != nil {
					log.Printf("Ошибка обновления yearlyRegister для %s: %v", userLogin, err)
				} else {
					log.Printf("Начислено годы регистрации для %s: %d лет, %d XP", userLogin, stats.YearsRegistered, yearlyExp)
					totalUserExp += yearlyExp
				}
			}

			// 3.5 Добавленные экземпляры
			if stats.CollectionItemsCount > 0 {
				var itemsExp, itemsLevel int64
				switch {
				case stats.CollectionItemsCount >= CollectionItemsAddLevel3:
					itemsExp = XPCollectionItemsAddLevel3
					itemsLevel = 3
				case stats.CollectionItemsCount >= CollectionItemsAddLevel2:
					itemsExp = XPCollectionItemsAddLevel2
					itemsLevel = 2
				case stats.CollectionItemsCount >= CollectionItemsAddLevel1:
					itemsExp = XPCollectionItemsAddLevel1
					itemsLevel = 1
				}

				if itemsExp > 0 {
					if err := updateAchievement(AchieveCollectionItemsAdded, itemsExp, itemsLevel, stats.CollectionItemsCount, false); err != nil {
						log.Printf("Ошибка обновления collectionItemsAdded для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено экземпляры для %s: %d XP (уровень %d)", userLogin, itemsExp, itemsLevel)
						totalUserExp += itemsExp
					}
				}
			}

			// 3.6 Добавленные фото
			if stats.PhotosCount > 0 {
				var photosExp, photosLevel int64
				switch {
				case stats.PhotosCount >= PhotosAddLevel3:
					photosExp = XPPhotosAddLevel3
					photosLevel = 3
				case stats.PhotosCount >= PhotosAddLevel2:
					photosExp = XPPhotosAddLevel2
					photosLevel = 2
				case stats.PhotosCount >= PhotosAddLevel1:
					photosExp = XPPhotosAddLevel1
					photosLevel = 1
				}

				if photosExp > 0 {
					if err := updateAchievement(AchievePhotosAdded, photosExp, photosLevel, stats.PhotosCount, false); err != nil {
						log.Printf("Ошибка обновления photosAdded для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено фото для %s: %d XP (уровень %d)", userLogin, photosExp, photosLevel)
						totalUserExp += photosExp
					}
				}
			}

			// 3.7 Созданные теги
			if stats.TagsCount > 0 {
				var tagsExp, tagsLevel int64
				switch {
				case stats.TagsCount >= TagsAddLevel3:
					tagsExp = XPTagsAddLevel3
					tagsLevel = 3
				case stats.TagsCount >= TagsAddLevel2:
					tagsExp = XPTagsAddLevel2
					tagsLevel = 2
				case stats.TagsCount >= TagsAddLevel1:
					tagsExp = XPTagsAddLevel1
					tagsLevel = 1
				}

				if tagsExp > 0 {
					if err := updateAchievement(AchieveTagsCreated, tagsExp, tagsLevel, stats.TagsCount, false); err != nil {
						log.Printf("Ошибка обновления tagsCreated для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено теги для %s: %d XP (уровень %d)", userLogin, tagsExp, tagsLevel)
						totalUserExp += tagsExp
					}
				}
			}

			// 3.8 Созданные посты
			if stats.PostsCount > 0 {
				var postsExp, postsLevel int64
				switch {
				case stats.PostsCount >= PostsCreateLevel3:
					postsExp = XPPostsCreateLevel3
					postsLevel = 3
				case stats.PostsCount >= PostsCreateLevel2:
					postsExp = XPPostsCreateLevel2
					postsLevel = 2
				case stats.PostsCount >= PostsCreateLevel1:
					postsExp = XPPostsCreateLevel1
					postsLevel = 1
				}

				if postsExp > 0 {
					if err := updateAchievement(AchievePostsCreated, postsExp, postsLevel, stats.PostsCount, false); err != nil {
						log.Printf("Ошибка обновления postsCreated для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено посты для %s: %d XP (уровень %d)", userLogin, postsExp, postsLevel)
						totalUserExp += postsExp
					}
				}
			}

			// 3.9 Лайки на коллекциях
			if stats.CollectionsLikes > 0 {
				var likesExp, likesLevel int64
				switch {
				case stats.CollectionsLikes >= CollectionsLikesLevel3:
					likesExp = XPCollectionsLikesLevel3
					likesLevel = 3
				case stats.CollectionsLikes >= CollectionsLikesLevel2:
					likesExp = XPCollectionsLikesLevel2
					likesLevel = 2
				case stats.CollectionsLikes >= CollectionsLikesLevel1:
					likesExp = XPCollectionsLikesLevel1
					likesLevel = 1
				}

				if likesExp > 0 {
					if err := updateAchievement(AchieveCollectionsLikes, likesExp, likesLevel, stats.CollectionsLikes, false); err != nil {
						log.Printf("Ошибка обновления collectionsLikes для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено лайки на коллекциях для %s: %d XP (уровень %d)", userLogin, likesExp, likesLevel)
						totalUserExp += likesExp
					}
				}
			}

			// 3.10 Лайки на экземплярах
			if stats.CollectionItemsLikes > 0 {
				var likesExp, likesLevel int64
				switch {
				case stats.CollectionItemsLikes >= CollectionItemsLikesLevel3:
					likesExp = XPCollectionItemsLikesLevel3
					likesLevel = 3
				case stats.CollectionItemsLikes >= CollectionItemsLikesLevel2:
					likesExp = XPCollectionItemsLikesLevel2
					likesLevel = 2
				case stats.CollectionItemsLikes >= CollectionItemsLikesLevel1:
					likesExp = XPCollectionItemsLikesLevel1
					likesLevel = 1
				}

				if likesExp > 0 {
					if err := updateAchievement(AchieveCollectionItemsLikes, likesExp, likesLevel, stats.CollectionItemsLikes, false); err != nil {
						log.Printf("Ошибка обновления collectionItemsLikes для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено лайки на экземплярах для %s: %d XP (уровень %d)", userLogin, likesExp, likesLevel)
						totalUserExp += likesExp
					}
				}
			}

			// 3.11 Лайки на фото
			if stats.PhotosLikes > 0 {
				var likesExp, likesLevel int64
				switch {
				case stats.PhotosLikes >= PhotosLikesLevel3:
					likesExp = XPPhotosLikesLevel3
					likesLevel = 3
				case stats.PhotosLikes >= PhotosLikesLevel2:
					likesExp = XPPhotosLikesLevel2
					likesLevel = 2
				case stats.PhotosLikes >= PhotosLikesLevel1:
					likesExp = XPPhotosLikesLevel1
					likesLevel = 1
				}

				if likesExp > 0 {
					if err := updateAchievement(AchievePhotosLikes, likesExp, likesLevel, stats.PhotosLikes, false); err != nil {
						log.Printf("Ошибка обновления photosLikes для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено лайки на фото для %s: %d XP (уровень %d)", userLogin, likesExp, likesLevel)
						totalUserExp += likesExp
					}
				}
			}

			// 3.12 Реакции на посты
			if stats.PostsReactions > 0 {
				var reactionsExp, reactionsLevel int64
				switch {
				case stats.PostsReactions >= PostsReactionsLevel3:
					reactionsExp = XPPostsReactionsLevel3
					reactionsLevel = 3
				case stats.PostsReactions >= PostsReactionsLevel2:
					reactionsExp = XPPostsReactionsLevel2
					reactionsLevel = 2
				case stats.PostsReactions >= PostsReactionsLevel1:
					reactionsExp = XPPostsReactionsLevel1
					reactionsLevel = 1
				}

				if reactionsExp > 0 {
					if err := updateAchievement(AchievePostsReactions, reactionsExp, reactionsLevel, stats.PostsReactions, false); err != nil {
						log.Printf("Ошибка обновления postsReactions для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено реакции на посты для %s: %d XP (уровень %d)", userLogin, reactionsExp, reactionsLevel)
						totalUserExp += reactionsExp
					}
				}
			}

			// 3.13 Сумма потраченная на коллекции
			if stats.CollectionsSum > 0 {
				var sumExp, sumLevel int64
				switch {
				case stats.CollectionsSum >= CollectionsSumLevel3:
					sumExp = XPCollectionsSumLevel3
					sumLevel = 3
				case stats.CollectionsSum >= CollectionsSumLevel2:
					sumExp = XPCollectionsSumLevel2
					sumLevel = 2
				case stats.CollectionsSum >= CollectionsSumLevel1:
					sumExp = XPCollectionsSumLevel1
					sumLevel = 1
				}

				if sumExp > 0 {
					if err := updateAchievement(AchieveCollectionsSum, sumExp, sumLevel, stats.CollectionsSum, false); err != nil {
						log.Printf("Ошибка обновления collectionsSum для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено сумма на коллекции для %s: %d XP (уровень %d)", userLogin, sumExp, sumLevel)
						totalUserExp += sumExp
					}
				}
			}

			// 3.14 Сумма потраченная на доставки
			if stats.CollectionShipSum > 0 {
				var shipExp, shipLevel int64
				switch {
				case stats.CollectionShipSum >= CollectionsShipSumLevel3:
					shipExp = XPCollectionsShipSumLevel3
					shipLevel = 3
				case stats.CollectionShipSum >= CollectionsShipSumLevel2:
					shipExp = XPCollectionsShipSumLevel2
					shipLevel = 2
				case stats.CollectionShipSum >= CollectionsShipSumLevel1:
					shipExp = XPCollectionsShipSumLevel1
					shipLevel = 1
				}

				if shipExp > 0 {
					if err := updateAchievement(AchieveCollectionShipSum, shipExp, shipLevel, stats.CollectionShipSum, false); err != nil {
						log.Printf("Ошибка обновления collectionShipSum для %s: %v", userLogin, err)
					} else {
						log.Printf("Начислено сумма на доставки для %s: %d XP (уровень %d)", userLogin, shipExp, shipLevel)
						totalUserExp += shipExp
					}
				}
			}

			// 4. Обновляем общий опыт пользователя в таблице loot.users (ДОБАВЛЯЕМ к существующему)
			if totalUserExp > 0 {
				if err := tx.Exec(`
					UPDATE loot.users 
					SET exp = COALESCE(exp, 0) + ?
					WHERE login = ? AND deleted_at IS NULL
				`, totalUserExp, userLogin).Error; err != nil {
					log.Printf("Ошибка обновления общего опыта для %s: %v", userLogin, err)
				} else {
					log.Printf("Добавлен опыт для %s: %d XP", userLogin, totalUserExp)
				}
			}
		}

		log.Println("Миграция успешно завершена!")
		return nil
	})
}
