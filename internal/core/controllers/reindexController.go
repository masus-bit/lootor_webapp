package controllers

import (
	"context"
	"fmt"
	"lootor/internal/core/repositories"
	"lootor/internal/core/services"
	"lootor/internal/pkg/elasticsearch"
	"net/http"

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
		"tags":             c.getTagData,
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
		"tags":             c.getTagData,
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
			"vkId":        user.VkId,
			"tgId":        user.TelegramId,
			"avatarUrl":   user.AvatarUrl,
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
				"id":          collection.Id.String(),
				"name":        collection.Name,
				"description": collection.Description,
				"isPrivate":   collection.IsPrivate,
				"bannerUrl":   collection.BannerUrl,
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
			"id":          item.Id.String(),
			"name":        item.Name,
			"description": item.Description,
			"images":      item.Images,
		}
	}
	return result, nil
}

func (c *ReindexController) getTagData() ([]map[string]interface{}, error) {
	tags, err := c.tagService.GetAllTagsForElastic()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	result := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = map[string]interface{}{
			"id":        tag.ID,
			"name":      tag.Name,
			"slug":      tag.Slug,
			"primaryId": tag.PrimaryID,
			"seriesId":  tag.SeriesID,
		}
	}
	return result, nil
}
