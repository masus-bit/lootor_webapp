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
	photosService  *PhotosService
	ciRepo         *repositories.CiRepository
	collectionRepo *repositories.CollectionsRepository
	eventsService  *EventsService
	itemTypeRepo   *repositories.ItemTypesRepository
	es             *elasticsearch.ElasticService
}

func NewTagsService(tagsClient *tagsclient.GRPCTagsClient, userRepo *repositories.UsersRepository, postsService *PostsService, ciRepo *repositories.CiRepository, collectionRepo *repositories.CollectionsRepository, eventsService *EventsService, itemTypeRepo *repositories.ItemTypesRepository, es *elasticsearch.ElasticService, photosService *PhotosService) *TagsService {
	return &TagsService{
		tagsClient:     tagsClient,
		userRepo:       userRepo,
		postsService:   postsService,
		ciRepo:         ciRepo,
		collectionRepo: collectionRepo,
		eventsService:  eventsService,
		itemTypeRepo:   itemTypeRepo,
		es:             es,
		photosService:  photosService,
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
		exp := int(utils.TagExp * float64(resp.GetDeletedCount()))
		_ = s.userRepo.IncrementExperience(authUser, exp)
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
			"id":        tag.GetId(),
			"name":      tag.GetName(),
			"slug":      tag.GetSlug(),
			"primaryId": tag.GetPrimaryId(),
			"seriesId":  tag.GetSeriesId(),
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

func (s *TagsService) MergeSeries(req *models.MergeTagsRequest) (*dto.CommonResponse, error) {
	_, err := s.tagsClient.MergeSeries(context.Background(), &microservices.MergeTagsRequest{
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
	resultTag.TotalPhotos = tag.GetTotalPhotos()
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
	author := authUser
	if role == "admin" {
		author = req.Author
	}
	canActivate, err := s.canActivate(authUser, role, req.EntityID, req.EntityType)
	if err != nil {
		return nil, err
	}
	if !canActivate {
		return nil, fmt.Errorf("вы не можете добавить тег к сущности, созданной другим пользователем")
	}

	_, err = s.tagsClient.AddTagsToEntity(context.Background(), &microservices.AddFewTagsToEntityRequest{
		EntityId:   req.EntityID,
		TagIds:     req.TagIDs,
		EntityType: req.EntityType,
		Author:     author,
	})
	_ = s.userRepo.IncrementExperience(author, utils.TagAttachExp)

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
		ID:          req.ID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении тега: %v", err)
	}
	return &models.TagDataResponse{Data: *s.convertProtoToModel(tag, "", false, "", nil)}, nil
}

func (s *TagsService) GetAllTagsForElastic() ([]models.ShortTags, error) {
	tags, err := s.tagsClient.FindAllTags(context.Background(), &microservices.FindAllTagsRequest{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []models.ShortTags
	for _, tag := range tags.GetTags() {
		result = append(result, models.ShortTags{
			ID:        tag.GetId(),
			Name:      tag.GetName(),
			Slug:      tag.GetSlug(),
			PrimaryID: tag.GetPrimaryId(),
			SeriesID:  tag.GetSeriesId(),
		})
	}
	return result, nil
}

func (s *TagsService) GetAllTags(limit, offset string) (*models.TagsDataResponse, error) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)
	tags, err := s.tagsClient.GetAllTags(context.Background(), &microservices.GetAllTagsRequest{Limit: int64(intLimit), Offset: int64(intOffset)})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []models.Tags
	for _, tag := range tags.GetTags() {
		normalizedTag := s.convertProtoToModel(tag, "", false, "", nil)
		result = append(result, *normalizedTag)
	}
	return &models.TagsDataResponse{Data: result, Total: tags.GetTotal()}, nil
}

func (s *TagsService) GetTagsByEntityId(entityId string) ([]models.ShortTags, error) {
	tags, err := s.tagsClient.GetTagsByEntityId(context.Background(), &microservices.GetTagsByEntityIdRequest{EntityId: entityId})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []models.ShortTags
	for _, tag := range tags.GetTags() {
		result = append(result, models.ShortTags{
			ID:        tag.GetId(),
			Name:      tag.GetName(),
			Slug:      tag.GetSlug(),
			PrimaryID: tag.GetPrimaryId(),
			SeriesID:  tag.GetSeriesId(),
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
	tag, err := s.tagsClient.GetTag(context.Background(), &microservices.GetTagRequest{ID: id})
	if err != nil {
		return nil, err
	}
	subsTags := user.TagsSubscriptions
	if slices.Contains(subsTags, id) {
		utils.RemoveByValue(subsTags, id)
		go func() {
			tagUUID, _ := uuid.Parse(tag.GetId())
			err = s.eventsService.AddEvent(userLogin, utils.EventActionUnsubscribe, utils.EventTargetTag, tag.GetName(), &models.EventsParams{TargetTagID: tagUUID})
			if err != nil {
				fmt.Println(err)
			}
		}()
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

func (s *TagsService) SearchSmartTags(query string, limit string) (*models.SearchResponse, error) {
	intLimit, _ := strconv.Atoi(limit)
	resp, err := s.tagsClient.SearchGameTitles(context.Background(), &microservices.SearchGameTitlesRequest{
		Query: query,
		Limit: int32(intLimit),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var results []models.Results
	for _, result := range resp.GetResults() {
		var seriesEntries []models.Tags
		for _, seriesEntry := range result.GetSeriesEntries() {
			seriesEntries = append(seriesEntries, *s.convertProtoToModel(seriesEntry, "", false, "", nil))
		}
		results = append(results, models.Results{
			Tag:   s.convertProtoToModel(result.GetTag(), "", false, "", nil),
			Score: result.GetScore(),
			ParsedData: &models.ParsedTitle{Original: result.GetParsedData().GetEdition(),
				SeriesCore: result.GetParsedData().GetSeriesCore(),
				GamePart:   result.GetParsedData().GetGamePart(),
				Edition:    result.GetParsedData().GetEdition(),
				Platform:   result.GetParsedData().GetPlatform(),
				Year:       result.GetParsedData().GetYear(),
				Confidence: result.GetParsedData().GetConfidence(),
			},
			SeriesTag:     s.convertProtoToModel(result.GetSeriesTag(), "", false, "", nil),
			MatchType:     result.GetMatchType(),
			Confidence:    result.GetConfidence(),
			SeriesEntries: seriesEntries,
		})
	}
	parsedData := &models.ParsedTitle{
		Original:   resp.GetParsedTitle().GetOriginal(),
		SeriesCore: resp.GetParsedTitle().GetSeriesCore(),
		GamePart:   resp.GetParsedTitle().GetGamePart(),
		Edition:    resp.GetParsedTitle().GetEdition(),
		Platform:   resp.GetParsedTitle().GetPlatform(),
		Year:       resp.GetParsedTitle().GetYear(),
		Confidence: resp.GetParsedTitle().GetConfidence(),
	}

	response := &models.SearchResponse{
		Data: &models.SearchResult{
			Results:      results,
			SearchTerm:   resp.GetSearchTerm(),
			ParsedData:   parsedData,
			TotalCount:   resp.GetTotalCount(),
			SearchTimeNs: resp.GetSearchTimeNs(),
			DidLearn:     resp.GetDidLearn(),
		},
	}
	return response, nil
}

func (s *TagsService) RecordUserChoice(req *models.UserChoiceRequest) (*dto.CommonResponse, error) {
	_, err := s.tagsClient.RecordUserChoice(context.Background(), &microservices.RecordUserChoiceRequest{
		SearchQuery:   req.SearchQuery,
		SelectedTagId: req.SelectedTagID,
		SessionId:     req.SessionID,
		UserId:        req.UserID,
		WasCorrect:    req.WasCorrect,
		FeedbackScore: req.FeedbackScore,
	})
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{
		Data: dto.Resp{Success: true},
	}, nil
}

func (s *TagsService) GetSuggestions(query, limit string) (*models.SearchSuggestionResponse, error) {
	intLimit, _ := strconv.Atoi(limit)
	resp, err := s.tagsClient.GetSearchSuggestions(context.Background(), &microservices.GetSearchSuggestionsRequest{
		Query: query,
		Limit: int32(intLimit),
	})
	if err != nil {
		return nil, err
	}
	var resultSuggestions []models.SearchSuggestions
	for _, suggestion := range resp.GetSuggestions() {
		resultSuggestions = append(resultSuggestions, models.SearchSuggestions{
			Title:      suggestion.GetTitle(),
			Type:       suggestion.GetType(),
			Score:      suggestion.GetScore(),
			Confidence: suggestion.GetConfidence(),
			TagID:      suggestion.GetTagId(),
		})
	}
	return &models.SearchSuggestionResponse{Data: resultSuggestions}, nil
}

func (s *TagsService) DeleteTags(req *models.DeleteTags) (*dto.CommonResponse, error) {
	_, err := s.tagsClient.DeleteTags(context.Background(), &microservices.GetTagsByIDsRequest{
		Ids: req.IDs,
	})
	if err != nil {
		return nil, err
	}
	for _, tagId := range req.IDs {
		if err := s.es.DeleteDocument(context.Background(), "tags", tagId); err != nil {
			log.Printf("Failed to delete tag from index: %v", err)
		}
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *TagsService) convertProtoToModel(tag *microservices.TagItem, userAuthLogin string, isPremium bool, filter string, tagsMap map[string]*microservices.GetShortsResponse) *models.Tags {
	var primaryTag *models.Tags
	var seriesTag *models.ShortTags
	var synonyms []models.Tags
	var seriesEntries []models.Tags
	var entities models.Entities
	var collectionItemProps *models.CollectionItemsProps
	collectionItemProps = nil
	author := &models.Users{}
	if tag != nil {
		author, _ = s.userRepo.GetUserByLogin(tag.GetAuthor())
		if tag.GetPrimaryId() != "" {
			primaryTag = s.convertProtoToModel(tag.Primary, userAuthLogin, isPremium, filter, tagsMap)
		} else {
			primaryTag = nil
		}

		if tag.GetSeriesId() != "" {
			seriesTag = &models.ShortTags{
				ID:        tag.GetSeriesId(),
				Name:      tag.GetSeries().GetName(),
				Slug:      tag.GetSeries().GetSlug(),
				PrimaryID: tag.GetSeries().GetPrimaryId(),
				SeriesID:  tag.GetSeries().GetSeriesId(),
			}
		} else {
			seriesTag = nil
		}

		if tag.GetSynonyms() != nil || len(tag.GetSynonyms()) > 0 {
			for _, synonym := range tag.GetSynonyms() {
				synonyms = append(synonyms, *s.convertProtoToModel(synonym, userAuthLogin, isPremium, filter, tagsMap))
			}
		} else {
			synonyms = []models.Tags{}
		}
		if tag.GetSeriesEntries() != nil || len(tag.GetSeriesEntries()) > 0 {
			for _, entry := range tag.GetSeriesEntries() {
				seriesEntries = append(seriesEntries, *s.convertProtoToModel(entry, userAuthLogin, isPremium, filter, tagsMap))
			}
		} else {
			seriesEntries = []models.Tags{}
		}
		if tag.GetEntities() != nil || len(tag.GetEntities()) > 0 {
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

		var authorTag *models.SubUsers
		authorTag = &models.SubUsers{}
		if author != nil {
			authorTag = &models.SubUsers{
				Login:       author.Login,
				AvatarUrl:   author.AvatarUrl,
				ProfileName: author.ProfileName,
				IsPremium:   author.IsPremium,
			}
		}

		return &models.Tags{
			ID:                   tag.ID,
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
			Author:               *authorTag,
			IsSeries:             tag.IsSeries,
			SeriesID:             tag.SeriesId,
			SeriesEntries:        seriesEntries,
			Series:               seriesTag,
		}
	}
	return &models.Tags{}
}

func (s *TagsService) getEntitiesByType(entityType, authUserLogin string, entityIDs []string, isPremium bool, filter string, tagsMap map[string]*microservices.GetShortsResponse) (*models.Entities, *models.CollectionItemsProps, error) {
	var entities models.Entities
	var collectionItemsProps *models.CollectionItemsProps = &models.CollectionItemsProps{}
	var filterId string
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
		counts, err := s.ciRepo.GetCICountsByItemType(entityIDs)
		if err != nil {
			return nil, nil, err
		}

		typeMap := make(map[string]*int64)
		for _, itemType := range itemTypes {
			if filter == itemType.Slug {
				filterId = itemType.ID.String()
			}

			if itemType.ID.String() != "" {
				typeMap[itemType.ID.String()] = nil
			}
		}

		for key, item := range counts {
			if _, exists := typeMap[key]; exists {
				for _, it := range itemTypes {
					if it.ID.String() != "" && it.ID.String() == key {
						switch it.Slug {
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
		}

		collectionItems, err := s.ciRepo.GetCollectionItemsByIDs(entityIDs, authUserLogin, filterId, tagsMap)
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
	} else if entityType == "photo" {
		photos, err := s.photosService.FindPhotosByIds(context.Background(), entityIDs, authUserLogin)
		if err != nil {
			return nil, nil, err
		}
		entities.Photos = photos.Data
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
	if entityType == "photo" {
		photo, _ := s.photosService.FindOneByID(context.Background(), entityID, authUser)
		authorLogin = photo.Data.Author.Login
	}

	if authorLogin != authUser {
		if role == "user" {
			return false, nil
		}
		return true, nil
	}

	return true, nil
}
