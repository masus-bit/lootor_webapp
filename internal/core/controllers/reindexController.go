package controllers

import (
	"context"
	"fmt"
	"lootor/internal/core/dto"
	"lootor/internal/core/repositories"
	"lootor/internal/core/services"
	"lootor/internal/pkg/elasticsearch"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type ReindexController struct {
	es             *elasticsearch.ElasticService
	userRepo       *repositories.UsersRepository
	collectionRepo *repositories.CollectionsRepository
	itemRepo       *repositories.CiRepository
	tagService     *services.TagsService
}

func NewReindexController(
	es *elasticsearch.ElasticService,
	userRepo *repositories.UsersRepository,
	collectionRepo *repositories.CollectionsRepository,
	itemRepo *repositories.CiRepository,
	tagService *services.TagsService,
) *ReindexController {
	return &ReindexController{
		es:             es,
		userRepo:       userRepo,
		collectionRepo: collectionRepo,
		itemRepo:       itemRepo,
		tagService:     tagService,
	}
}

func (c *ReindexController) Reindex(ctx echo.Context) error {
	if err := c.reindexAll(ctx.Request().Context()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Reindexing completed successfully"})
}

func (c *ReindexController) ReindexInternal(ctx context.Context) error {
	return c.reindexAll(ctx)
}

func (c *ReindexController) reindexAll(ctx context.Context) error {
	const batchSize = 500

	// Канал для ошибок (размер = кол-во горутин + 1 для тегов)
	errChan := make(chan error, 4)
	var wg sync.WaitGroup

	// 1. Запускаем индексацию тегов в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.tagService.TriggerReindex(); err != nil {
			errChan <- fmt.Errorf("failed to reindex tags: %w", err)
		}
	}()

	// 2. Запускаем индексацию остальных сущностей
	dataProviders := map[string]func() ([]map[string]interface{}, error){
		"users":            c.getUserData,
		"collections":      c.getCollectionData,
		"collection_items": c.getCollectionItemData,
	}

	for entityName, provider := range dataProviders {
		wg.Add(1)
		go func(name string, dataProvider func() ([]map[string]interface{}, error)) {
			defer wg.Done()

			// Проверяем отмену контекста
			select {
			case <-ctx.Done():
				errChan <- fmt.Errorf("%s: %w", name, ctx.Err())
				return
			default:
			}

			data, err := dataProvider()
			if err != nil {
				errChan <- fmt.Errorf("failed to get %s data: %w", name, err)
				return
			}

			if len(data) == 0 {
				return
			}

			// Индексируем пачками
			for i := 0; i < len(data); i += batchSize {
				// Проверяем отмену контекста перед каждой пачкой
				select {
				case <-ctx.Done():
					errChan <- fmt.Errorf("%s: %w", name, ctx.Err())
					return
				default:
				}

				end := i + batchSize
				if end > len(data) {
					end = len(data)
				}

				batch := data[i:end]

				if err := c.es.BulkIndexDocuments(ctx, name, batch); err != nil {
					errChan <- fmt.Errorf("failed to index %s batch %d-%d: %w", name, i, end, err)
					return
				}

				// Небольшая задержка между батчами, чтобы не перегружать ES
				time.Sleep(100 * time.Millisecond)
			}

		}(entityName, provider)
	}

	// Ожидаем завершения всех горутин
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Собираем ошибки
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("reindexing completed with errors: %v", errors)
	}

	return nil
}

func (c *ReindexController) getUserData() ([]map[string]interface{}, error) {
	users, err := c.userRepo.FindAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	result := make([]map[string]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"id":          user.Login,
			"login":       user.Login,
			"userName":    user.UserName,
			"email":       user.Email,
			"vkId":        user.VkID,
			"tgId":        user.TelegramID,
			"avatarUrl":   user.AvatarURL,
			"isPremium":   user.IsPremium,
			"profileName": user.ProfileName,
		}
	}
	return result, nil
}

func (c *ReindexController) getCollectionData() ([]map[string]interface{}, error) {
	collections, err := c.collectionRepo.FindAllCollections()
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	result := make([]map[string]interface{}, 0)
	for _, collection := range collections {
		if !collection.IsPrivate {
			result = append(
				result, map[string]interface{}{
					"id":          collection.ID.String(),
					"name":        collection.Name,
					"description": collection.Description,
					"isPrivate":   collection.IsPrivate,
					"bannerUrl":   collection.BannerURL,
					"owner": dto.SubUsers{
						Login:         collection.User.Login,
						AvatarURL:     collection.User.AvatarURL,
						ProfileName:   collection.User.ProfileName,
						IsPremium:     collection.User.IsPremium,
						DonateTotal:   "0",
						BackgroundUrl: collection.User.BackgroundURL,
					},
				},
			)
		}
	}
	return result, nil
}

func (c *ReindexController) getCollectionItemData() ([]map[string]interface{}, error) {
	items, err := c.itemRepo.FindAllCI()
	if err != nil {
		return nil, fmt.Errorf("failed to get collection items: %w", err)
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":          item.ID.String(),
			"name":        item.Name,
			"description": item.Description,
			"images":      item.Images,
			"owner": dto.SubUsers{
				Login:         item.Owner.Login,
				AvatarURL:     item.Owner.AvatarURL,
				ProfileName:   item.Owner.ProfileName,
				IsPremium:     item.Owner.IsPremium,
				DonateTotal:   "0",
				BackgroundUrl: item.Owner.BackgroundURL,
			},
		}
	}
	return result, nil
}
