package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/utils"
	"slices"
	"strconv"
	"sync"
)

type TagsService struct {
	tagsClient     *tagsclient.GRPCTagsClient
	userRepo       *repositories.UsersRepository
	postsService   *PostsService
	ciRepo         *repositories.CiRepository
	collectionRepo *repositories.CollectionsRepository
	eventsService  *EventsService
	itemTypeRepo   *repositories.ItemTypesRepository
	es             *elasticsearch.ElasticService
}

func NewTagsService(tagsClient *tagsclient.GRPCTagsClient, userRepo *repositories.UsersRepository, postsService *PostsService, ciRepo *repositories.CiRepository, collectionRepo *repositories.CollectionsRepository, eventsService *EventsService, itemTypeRepo *repositories.ItemTypesRepository, es *elasticsearch.ElasticService) *TagsService {
	return &TagsService{
		tagsClient:     tagsClient,
		userRepo:       userRepo,
		postsService:   postsService,
		ciRepo:         ciRepo,
		collectionRepo: collectionRepo,
		eventsService:  eventsService,
		itemTypeRepo:   itemTypeRepo,
		es:             es,
	}
}

func (s *TagsService) CreateTag(req *models.TagCreateRequest, authUser string) (*models.TagIDsResponse, error) {
	resp, err := s.tagsClient.CreateTag(context.Background(), req)
	if err != nil {
		return nil, err
	}

	tags := resp.GetTags()
	if len(tags) == 0 {
		return &models.TagIDsResponse{Data: []string{}}, nil
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.bulkIndexTags(tags); err != nil {
			log.Printf("Failed to index tags in bulk: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		exp := int(utils.TagExp * float64(len(tags)))
		// Выполняем последовательно, но параллельно с индексацией
		_ = s.userRepo.IncrementExperience(authUser, exp)
		_ = s.userRepo.IncrementSocialScore(authUser, len(tags))
	}()

	tagIDs := make([]string, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.GetId()
	}

	return &models.TagIDsResponse{Data: tagIDs}, nil
}

func (s *TagsService) bulkIndexTags(tags []*microservices.TagCreateResponse) error {
	docs := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		docs[i] = map[string]interface{}{
			"id":   tag.GetId(),
			"name": tag.GetName(),
			"slug": tag.GetSlug(),
		}
	}

	return s.es.BulkIndexDocuments(context.Background(), "tags", docs)
}

func (s *TagsService) SearchTags(name string) (*models.TagsDataResponse, error) {
	searchTerm := fmt.Sprintf("%%%s%%", name)

	tags, err := s.tagsClient.SearchTags(context.Background(), &models.TagsSearchRequest{Name: searchTerm})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}

	var result []models.Tags
	for _, tag := range tags.GetTags() {
		result = append(result, *s.convertProtoToModel(tag, "", false, "", nil))
	}

	return &models.TagsDataResponse{Data: result}, nil
}

func (s *TagsService) MergeTags(req *models.MergeTagsRequest) (*dto.CommonResponse, error) {
	_, err := s.tagsClient.MergeTags(context.Background(), &microservices.MergeTagsRequest{
		FromTagIds: req.FromTagIDs,
		ToTagId:    req.ToTagID,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при слиянии тегов: %v", err)
	}
	return &dto.CommonResponse{
		Data: dto.Resp{Success: true},
	}, nil
}

func (s *TagsService) FindAllEntitiesByTag(tagID, entityType, limit, offset, authUserLogin, ciFilter string) (*models.TagDataResponse, error) {
	authUser, _ := s.userRepo.GetUserByLogin(authUserLogin)
	var isPremium bool
	if authUser != nil {
		isPremium = authUser.IsPremium
	} else {
		isPremium = false
	}
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)
	tag, err := s.tagsClient.FindAllEntitiesByTag(context.Background(), &microservices.GetEntitiesByTagRequest{TagId: tagID, EntityType: entityType, Limit: int64(intLimit), Offset: int64(intOffset)})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	author, _ := s.userRepo.GetUserByLogin(tag.GetTag().GetAuthor())

	tagsMap := tag.TagsMap
	resultTag := s.convertProtoToModel(tag.GetTag(), authUserLogin, isPremium, ciFilter, tagsMap)
	resultTag.TotalPosts = tag.GetTotalPosts()
	resultTag.TotalCollections = tag.GetTotalCollections()
	resultTag.TotalCollectionItems = tag.GetTotalCollectionItems()
	if author != nil {
		resultTag.Author = models.SubUsers{
			Login:       author.Login,
			AvatarUrl:   author.AvatarUrl,
			ProfileName: author.ProfileName,
			IsPremium:   author.IsPremium,
		}
	}

	return &models.TagDataResponse{Data: *resultTag, Total: tag.GetTotalEntities()}, nil
}

func (s *TagsService) AddTagToEntity(req *models.AddTagToEntityRequest, authUser, role string) (*dto.CommonResponse, error) {
	canActivate, err := s.canActivate(authUser, role, req.EntityID, req.EntityType)
	if err != nil {
		return nil, err
	}
	if !canActivate {
		return nil, fmt.Errorf("вы не можете добавить тег к сущности, созданной другим пользователем")
	}

	_, err = s.tagsClient.AddTagToEntity(context.Background(), &microservices.AddTagsToEntityRequest{
		EntityId:   req.EntityID,
		TagId:      req.TagID,
		EntityType: req.EntityType,
	})
	_ = s.userRepo.IncrementExperience(authUser, utils.TagAttachExp)

	if err != nil {
		return nil, fmt.Errorf("ошибка при добавлении тега: %v", err)
	}
	return &dto.CommonResponse{
		Data: dto.Resp{Success: true},
	}, nil
}

func (s *TagsService) RemoveTagsFromEntity(req *models.RemoveTagsRequest, authUser, role string) (*dto.CommonResponse, error) {
	canActivate, err := s.canActivate(authUser, role, req.EntityID, req.EntityType)
	if err != nil {
		return nil, err
	}
	if !canActivate {
		return nil, fmt.Errorf("вы не можете удалять тег сущности, созданной другим пользователем")
	}

	_, err = s.tagsClient.RemoveTagsFromEntity(context.Background(), &microservices.RemoveTagsRequest{
		EntityId:   req.EntityID,
		TagIds:     req.TagIDs,
		EntityType: req.EntityType,
	})
	_ = s.userRepo.DecrementExperience(authUser, int(utils.TagAttachExp*float64(len(req.TagIDs))))

	if err != nil {
		return nil, fmt.Errorf("ошибка при удалении тегov: %v", err)
	}
	return &dto.CommonResponse{
		Data: dto.Resp{Success: true},
	}, nil
}

func (s *TagsService) UpdateTag(req *models.TagUpdateRequest) (*models.TagDataResponse, error) {
	tag, err := s.tagsClient.UpdateTag(context.Background(), &microservices.UpdateTagRequest{
		Id:          req.ID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении тега: %v", err)
	}
	return &models.TagDataResponse{Data: *s.convertProtoToModel(tag, "", false, "", nil)}, nil
}

func (s *TagsService) GetAllTags() ([]models.ShortTags, error) {
	tags, err := s.tagsClient.GetAllTags(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []models.ShortTags
	for _, tag := range tags.GetTags() {
		result = append(result, models.ShortTags{
			ID:   tag.GetId(),
			Name: tag.GetName(),
			Slug: tag.GetSlug(),
		})
	}
	return result, nil
}

func (s *TagsService) GetTagsByEntityId(entityId string) ([]models.ShortTags, error) {
	tags, err := s.tagsClient.GetTagsByEntityId(context.Background(), &microservices.GetTagsByEntityIdRequest{EntityId: entityId})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []models.ShortTags
	for _, tag := range tags.GetTags() {
		result = append(result, models.ShortTags{
			ID:   tag.GetId(),
			Name: tag.GetName(),
			Slug: tag.GetSlug(),
		})
	}
	return result, nil
}

func (s *TagsService) GetTagBySlug(slug string) (*models.Tags, error) {
	tag, err := s.tagsClient.GetTagBySlug(context.Background(), slug)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тега: %v", err)
	}
	return s.convertProtoToModel(tag, "", false, "", nil), nil
}

func (s *TagsService) Subscribe(id, userLogin string) (*dto.CommonResponse, error) {
	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, err
	}
	tag, err := s.tagsClient.GetTag(context.Background(), &microservices.GetTagRequest{Id: id})
	if err != nil {
		return nil, err
	}
	subsTags := user.TagsSubscriptions
	if slices.Contains(subsTags, id) {
		utils.RemoveByValue(subsTags, id)
	} else {
		subsTags = append(subsTags, id)
		go func() {
			tagUUID, _ := uuid.Parse(tag.GetId())
			err = s.eventsService.AddEvent(userLogin, utils.EventActionSubscribe, utils.EventTargetTag, tag.GetName(), &models.EventsParams{TargetTagID: tagUUID})
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	user.TagsSubscriptions = subsTags
	_, err = s.userRepo.UpdateUser(user, *user)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{
		Data: dto.Resp{Success: true},
	}, nil
}

func (s *TagsService) convertProtoToModel(tag *microservices.TagItem, userAuthLogin string, isPremium bool, filter string, tagsMap map[string]*microservices.GetShortsResponse) *models.Tags {
	var primaryTag *models.Tags
	var synonyms []models.Tags
	var entities models.Entities
	var collectionItemProps *models.CollectionItemsProps
	collectionItemProps = nil

	if tag.PrimaryId != "" {
		primaryTag = s.convertProtoToModel(tag.Primary, userAuthLogin, isPremium, filter, tagsMap)
	} else {
		primaryTag = nil
	}
	if tag.Synonyms != nil || len(tag.Synonyms) > 0 {
		for _, synonym := range tag.Synonyms {
			synonyms = append(synonyms, *s.convertProtoToModel(synonym, userAuthLogin, isPremium, filter, tagsMap))
		}
	} else {
		synonyms = []models.Tags{}
	}
	if tag.Entities != nil || len(tag.Entities) > 0 {
		entitiesType := tag.Entities[0].EntityType
		var entityIDs []string
		for _, entity := range tag.Entities {
			entityIDs = append(entityIDs, entity.EntityId)
		}
		resEntities, ciProps, err := s.getEntitiesByType(entitiesType, userAuthLogin, entityIDs, isPremium, filter, tagsMap)
		if err != nil {
			return nil
		}
		if ciProps != nil {
			collectionItemProps = ciProps
		}
		entities = *resEntities
	} else {
		entities = models.Entities{}
	}
	return &models.Tags{
		ID:                   tag.Id,
		Name:                 tag.Name,
		Slug:                 tag.Slug,
		Description:          tag.Description,
		CreatedAt:            tag.CreatedAt,
		PrimaryID:            tag.PrimaryId,
		Primary:              primaryTag,
		Entities:             entities,
		Synonyms:             synonyms,
		IsPrimary:            tag.IsPrimary,
		CollectionItemsProps: collectionItemProps,
	}
}

func (s *TagsService) getEntitiesByType(entityType, authUserLogin string, entityIDs []string, isPremium bool, filter string, tagsMap map[string]*microservices.GetShortsResponse) (*models.Entities, *models.CollectionItemsProps, error) {
	var entities models.Entities
	var collectionItemsProps *models.CollectionItemsProps
	collectionItemsProps = nil
	if entityType == "post" {
		posts, err := s.postsService.GetPostsByIDs(context.Background(), entityIDs, authUserLogin, isPremium)
		if err != nil {
			return nil, nil, err
		}
		entities.Posts = posts.Data
	} else if entityType == "collectionItem" {

		itemTypes, err := s.itemTypeRepo.FindAllTypes()
		if err != nil {
			return nil, nil, err
		}
		counts, err := s.ciRepo.GetCICountsByItemType()
		if err != nil {
			return nil, nil, err
		}

		for key, item := range counts {
			for _, itemType := range itemTypes {
				if key == itemType.Id.String() {
					switch itemType.Slug {
					case "books":
						collectionItemsProps.Books = item
					case "videoGames":
						collectionItemsProps.VideoGames = item
					case "boardGames":
						collectionItemsProps.BoardGames = item
					case "comics":
						collectionItemsProps.Comics = item
					case "gamingHardware":
						collectionItemsProps.GamingHardware = item
					case "vinyl":
						collectionItemsProps.Vinyl = item
					case "steelbooks":
						collectionItemsProps.Steelbooks = item
					case "collectibleCards":
						collectionItemsProps.CollectibleCards = item
					case "collectibleFigures":
						collectionItemsProps.CollectibleFigures = item
					}
					break
				}
			}
		}

		collectionItems, err := s.ciRepo.GetCollectionItemsByIDs(entityIDs, authUserLogin, filter, tagsMap)
		if err != nil {
			return nil, nil, err
		}
		entities.CollectionItems = collectionItems
	} else if entityType == "collection" {
		collectionIDsUUID := make([]uuid.UUID, len(entityIDs))
		for i, id := range entityIDs {
			collectionIDsUUID[i], _ = uuid.Parse(id)
		}
		counts, _ := s.ciRepo.GetCountCIByIDs(collectionIDsUUID)
		totalPrices, _ := s.ciRepo.GetSumsByCollectionIDs(collectionIDsUUID)
		shippingCosts, _ := s.ciRepo.GetShippingCostsByCollectionIDs(collectionIDsUUID)
		collections, err := s.collectionRepo.GetCollectionsByIds(entityIDs, counts, totalPrices, shippingCosts, authUserLogin, tagsMap)
		if err != nil {
			return nil, nil, err
		}
		entities.Collections = collections
	} else {
		return nil, nil, fmt.Errorf("unknown entity type: %s", entityType)
	}

	return &entities, collectionItemsProps, nil
}

func (s *TagsService) canActivate(authUser, role, entityID, entityType string) (bool, error) {
	var authorLogin string
	if entityType == "post" {
		postId, _ := strconv.ParseUint(entityID, 10, 64)
		postResp, _ := s.postsService.GetPostById(context.Background(), postId, authUser)
		authorLogin = postResp.Data.Author.Login
	}
	if entityType == "collection" {
		collection, _ := s.collectionRepo.GetByIdWithoutCollectionItems(entityID)
		authorLogin = collection.UserLogin
	}
	if entityType == "collectionItem" {
		collectionItem, _ := s.ciRepo.GetCIByID(entityID)
		authorLogin = collectionItem.UserLogin
	}

	if authorLogin != authUser {
		if role == "user" {
			return false, nil
		}
		return true, nil
	}

	return true, nil
}
