package controllers

import (
	"context"
	"fmt"
	"lootor/internal/core/repositories"
	"lootor/internal/core/services"
	"lootor/internal/pkg/elasticsearch"
	"net/http"
	"sync"

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
	dataProviders := map[string]func() ([]map[string]interface{}, error){
		"users":            c.getUserData,
		"collections":      c.getCollectionData,
		"collection_items": c.getCollectionItemData,
		//"tags":             c.getTagData,
	}

	var errTags error

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := c.tagService.TriggerReindex()
		if err != nil {
			errTags = err
		}
	}()
	wg.Wait()

	if errTags != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if err := c.es.ReindexAll(ctx.Request().Context(), dataProviders); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Reindexing completed successfully"})
}

func (c *ReindexController) ReindexInternal(ctx context.Context) error {
	dataProviders := map[string]func() ([]map[string]interface{}, error){
		"users":            c.getUserData,
		"collections":      c.getCollectionData,
		"collection_items": c.getCollectionItemData,
		//"tags":             c.getTagData,
	}

	var errTags error

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := c.tagService.TriggerReindex()
		if err != nil {
			errTags = err
		}
	}()
	wg.Wait()

	if errTags != nil {
		return errTags
	}

	return c.es.ReindexAll(ctx, dataProviders)
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

	result := make([]map[string]interface{}, len(collections))
	for i, collection := range collections {
		if !collection.IsPrivate {
			result[i] = map[string]interface{}{
				"id":          collection.ID.String(),
				"name":        collection.Name,
				"description": collection.Description,
				"isPrivate":   collection.IsPrivate,
				"bannerUrl":   collection.BannerURL,
			}
		}
		continue
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
		}
	}
	return result, nil
}

//func (c *ReindexController) getTagData() ([]map[string]interface{}, error) {
//	tags, err := c.tagService.GetAllTagsForElastic()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get tags: %w", err)
//	}
//
//	result := make([]map[string]interface{}, len(tags))
//	for i, tag := range tags {
//		result[i] = map[string]interface{}{
//			"id":                   tag.ID,
//			"name":                 tag.Name,
//			"slug":                 tag.Slug,
//			"primaryId":            tag.PrimaryID,
//			"seriesId":             tag.SeriesID,
//			"totalPosts":           tag.TotalPosts,
//			"totalCollectionItems": tag.TotalCollectionItems,
//			"totalPhotos":          tag.TotalPhotos,
//			"totalCollections":     tag.TotalCollections,
//		}
//	}
//	return result, nil
//}
