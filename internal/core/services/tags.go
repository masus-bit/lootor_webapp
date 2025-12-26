package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/achievementsclient"
	"lootor/internal/infrastructure/tagsclient"
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
	achService     *achievementsclient.GRPCAchievementsClient
}

func NewTagsService(tagsClient *tagsclient.GRPCTagsClient, userRepo *repositories.UsersRepository, postsService *PostsService, ciRepo *repositories.CiRepository, collectionRepo *repositories.CollectionsRepository, eventsService *EventsService, itemTypeRepo *repositories.ItemTypesRepository, es *elasticsearch.ElasticService, photosService *PhotosService, achService *achievementsclient.GRPCAchievementsClient) *TagsService {
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
		achService:     achService,
	}
}

func (s *TagsService) CreateTag(req *dto.TagCreateRequest, authUser string) (*dto.TagIDsResponse, error) {
	resp, err := s.tagsClient.CreateTag(context.Background(), req)
	if err != nil {
		return nil, err
	}

	tags := resp.GetTags()
	if len(tags) == 0 {
		return &dto.TagIDsResponse{Data: []string{}}, nil
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
		tagsCount, err := s.tagsClient.GetTagsTotalByUserLogin(context.Background(), authUser)
		if err != nil {
			fmt.Println(err)
		}
		xp, level := utils.GetAchievementTagsAddData(tagsCount.GetTotalTags())
		err = utils.AddAchievement(s.achService, utils.AchieveTagsCreated, authUser, level, xp, tagsCount.GetTotalTags())
		_ = s.userRepo.IncrementExperience(authUser, exp+int(xp))
	}()

	tagIDs := make([]string, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.GetId()
	}

	return &dto.TagIDsResponse{Data: tagIDs}, nil
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

func (s *TagsService) SearchTags(name string) (*dto.TagsDataResponse, error) {
	searchTerm := fmt.Sprintf("%%%s%%", name)

	tags, err := s.tagsClient.SearchTags(context.Background(), &dto.TagsSearchRequest{Name: searchTerm})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}

	var result []dto.Tags
	for _, tag := range tags.GetTags() {
		result = append(result, *s.convertProtoToModel(tag, "", false, "", nil))
	}

	return &dto.TagsDataResponse{Data: result}, nil
}

func (s *TagsService) MergeTags(req *dto.MergeTagsRequest) (*dto.CommonResponse, error) {
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

func (s *TagsService) MergeSeries(req *dto.MergeTagsRequest) (*dto.CommonResponse, error) {
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

func (s *TagsService) FindAllEntitiesByTag(tagID, entityType, limit, offset, authUserLogin, ciFilter string) (*dto.TagDataResponse, error) {
	authUser, _ := s.userRepo.GetUserByLogin(authUserLogin)
	var isPremium bool
	if authUser != nil {
		isPremium = authUser.IsPremium
	} else {
		isPremium = false
	}
	const maxInt32 = 1<<31 - 1

	var finalLimit int64

	intLimit64, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		finalLimit = 0
	}
	finalLimit = intLimit64
	if finalLimit > maxInt32 {
		return nil, fmt.Errorf("limit too large: %d", intLimit64)
	}

	intLimit := int32(finalLimit)
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
		resultTag.Author = dto.SubUsers{
			Login:       author.Login,
			AvatarURL:   author.AvatarURL,
			ProfileName: author.ProfileName,
			IsPremium:   author.IsPremium,
		}
	}

	return &dto.TagDataResponse{Data: *resultTag, Total: tag.GetTotalEntities()}, nil
}

func (s *TagsService) AddTagToEntity(req *dto.AddTagToEntityRequest, authUser, role string) (*dto.CommonResponse, error) {
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

func (s *TagsService) RemoveTagsFromEntity(req *dto.RemoveTagsRequest, authUser, role string) (*dto.CommonResponse, error) {
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

func (s *TagsService) UpdateTag(req *dto.TagUpdateRequest) (*dto.TagDataResponse, error) {
	addFields, err := utils.GormJSONToProtoStruct(req.AdditionalFields)
	if err != nil {
		return nil, err
	}
	tag, err := s.tagsClient.UpdateTag(context.Background(), &microservices.UpdateTagRequest{
		Id:               req.ID,
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		AdditionalFields: addFields,
		OnModeration:     req.OnModeration,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении тега: %v", err)
	}
	return &dto.TagDataResponse{Data: *s.convertProtoToModel(tag, "", false, "", nil)}, nil
}

func (s *TagsService) GetAllTagsForElastic() ([]dto.ShortTags, error) {
	tags, err := s.tagsClient.FindAllTags(context.Background(), &microservices.FindAllTagsRequest{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []dto.ShortTags
	for _, tag := range tags.GetTags() {
		result = append(result, dto.ShortTags{
			ID:        tag.GetId(),
			Name:      tag.GetName(),
			Slug:      tag.GetSlug(),
			PrimaryID: tag.GetPrimaryId(),
			SeriesID:  tag.GetSeriesId(),
		})
	}
	return result, nil
}

func (s *TagsService) GetAllTags(limit, offset, authUserLogin string) (*dto.TagsDataResponse, error) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)
	tags, err := s.tagsClient.GetAllTags(context.Background(), &microservices.GetAllTagsRequest{Limit: int64(intLimit), Offset: int64(intOffset)})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []dto.Tags
	for _, tag := range tags.GetTags() {
		normalizedTag := s.convertProtoToModel(tag, authUserLogin, false, "", nil)
		result = append(result, *normalizedTag)
	}
	return &dto.TagsDataResponse{Data: result, Total: tags.GetTotal()}, nil
}

func (s *TagsService) GetTagsByEntityId(entityId string) ([]dto.ShortTags, error) {
	tags, err := s.tagsClient.GetTagsByEntityId(context.Background(), &microservices.GetTagsByEntityIdRequest{EntityId: entityId})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var result []dto.ShortTags
	for _, tag := range tags.GetTags() {
		result = append(result, dto.ShortTags{
			ID:        tag.GetId(),
			Name:      tag.GetName(),
			Slug:      tag.GetSlug(),
			PrimaryID: tag.GetPrimaryId(),
			SeriesID:  tag.GetSeriesId(),
		})
	}
	return result, nil
}

func (s *TagsService) GetTagBySlug(slug string) (*dto.Tags, error) {
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
		subsTags = utils.RemoveByValue(subsTags, id)
		go func() {
			tagUUID, _ := uuid.Parse(tag.GetId())
			err = s.eventsService.AddEvent(userLogin, utils.EventActionUnsubscribe, utils.EventTargetTag, tag.GetName(), &dto.EventsParams{TargetTagID: tagUUID})
			if err != nil {
				fmt.Println(err)
			}
		}()
	} else {
		subsTags = append(subsTags, id)
		go func() {
			tagUUID, _ := uuid.Parse(tag.GetId())
			err = s.eventsService.AddEvent(userLogin, utils.EventActionSubscribe, utils.EventTargetTag, tag.GetName(), &dto.EventsParams{TargetTagID: tagUUID})
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	updUser := user
	updUser.TagsSubscriptions = subsTags
	_, err = s.userRepo.UpdateUser(user, *updUser)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{
		Data: dto.Resp{Success: true},
	}, nil
}

func (s *TagsService) SearchSmartTags(query string, limit string) (*dto.SearchResponse, error) {
	const maxInt32 = 1<<31 - 1

	var finalLimit int64

	intLimit64, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		finalLimit = 0
	}
	finalLimit = intLimit64
	if finalLimit > maxInt32 {
		return nil, fmt.Errorf("limit too large: %d", intLimit64)
	}

	intLimit := int32(finalLimit)
	resp, err := s.tagsClient.SearchGameTitles(context.Background(), &microservices.SearchGameTitlesRequest{
		Query: query,
		Limit: intLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске тегов: %v", err)
	}
	var results []dto.Results
	for _, result := range resp.GetResults() {
		var seriesEntries []dto.Tags
		for _, seriesEntry := range result.GetSeriesEntries() {
			seriesEntries = append(seriesEntries, *s.convertProtoToModel(seriesEntry, "", false, "", nil))
		}
		results = append(results, dto.Results{
			Tag:   s.convertProtoToModel(result.GetTag(), "", false, "", nil),
			Score: result.GetScore(),
			ParsedData: &dto.ParsedTitle{Original: result.GetParsedData().GetEdition(),
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
	parsedData := &dto.ParsedTitle{
		Original:   resp.GetParsedTitle().GetOriginal(),
		SeriesCore: resp.GetParsedTitle().GetSeriesCore(),
		GamePart:   resp.GetParsedTitle().GetGamePart(),
		Edition:    resp.GetParsedTitle().GetEdition(),
		Platform:   resp.GetParsedTitle().GetPlatform(),
		Year:       resp.GetParsedTitle().GetYear(),
		Confidence: resp.GetParsedTitle().GetConfidence(),
	}

	response := &dto.SearchResponse{
		Data: &dto.SearchResult{
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

func (s *TagsService) RecordUserChoice(req *dto.UserChoiceRequest) (*dto.CommonResponse, error) {
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

func (s *TagsService) GetSuggestions(query, limit string) (*dto.SearchSuggestionResponse, error) {
	const maxInt32 = 1<<31 - 1

	var finalLimit int64

	intLimit64, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		finalLimit = 0
	}
	finalLimit = intLimit64
	if finalLimit > maxInt32 {
		return nil, fmt.Errorf("limit too large: %d", intLimit64)
	}

	intLimit := int32(finalLimit)
	resp, err := s.tagsClient.GetSearchSuggestions(context.Background(), &microservices.GetSearchSuggestionsRequest{
		Query: query,
		Limit: intLimit,
	})
	if err != nil {
		return nil, err
	}
	var resultSuggestions []dto.SearchSuggestions
	for _, suggestion := range resp.GetSuggestions() {
		resultSuggestions = append(resultSuggestions, dto.SearchSuggestions{
			Title:      suggestion.GetTitle(),
			Type:       suggestion.GetType(),
			Score:      suggestion.GetScore(),
			Confidence: suggestion.GetConfidence(),
			TagID:      suggestion.GetTagId(),
		})
	}
	return &dto.SearchSuggestionResponse{Data: resultSuggestions}, nil
}

func (s *TagsService) DeleteTags(req *dto.DeleteTags) (*dto.CommonResponse, error) {
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

func (s *TagsService) PublicUpdate(req *dto.PublicUpdateTag) (*dto.TagDataResponse, error) {
	addFields, err := utils.GormJSONToProtoStruct(req.AdditionalFields)
	if err != nil {
		return nil, err
	}
	resp, err := s.tagsClient.PublicUpdateTag(context.Background(), &microservices.PublicUpdateTagRequest{
		Id:               req.ID,
		AdditionalFields: addFields,
	})
	if err != nil {
		return nil, err
	}
	tagResult := s.convertProtoToModel(resp.GetTag(), "", false, "", nil)

	return &dto.TagDataResponse{Data: *tagResult}, nil
}

func (s *TagsService) MoveTagLinks(req *dto.MoveTagLinks) (*dto.CommonResponse, error) {
	_, err := s.tagsClient.MoveTagLinks(context.Background(), &microservices.MoveTagLinksRequest{
		TargetId: req.TargetID,
		SourceId: req.SourceID,
	})
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *TagsService) GetOnModerationTags(req *dto.OnModerationTagsRequest) (*dto.TagsDataResponse, error) {
	intLimit, _ := strconv.Atoi(req.Limit)
	intOffset, _ := strconv.Atoi(req.Offset)
	resp, err := s.tagsClient.GetOnModerationTags(context.Background(), &microservices.OnModerationTagsRequest{
		Limit:  int64(intLimit),
		Offset: int64(intOffset),
	})
	if err != nil {
		return nil, err
	}
	var result []dto.Tags

	for _, tag := range resp.GetTags() {
		normalizedTag := s.convertProtoToModel(tag, "", false, "", nil)
		result = append(result, *normalizedTag)
	}
	return &dto.TagsDataResponse{Data: result, Total: resp.GetTotal()}, nil

}

func (s *TagsService) convertProtoToModel(tag *microservices.TagItem, userAuthLogin string, isPremium bool, filter string, tagsMap map[string]*microservices.GetShortsResponse) *dto.Tags {
	var primaryTag *dto.Tags
	var seriesTag *dto.ShortTags
	var synonyms []dto.Tags
	var seriesEntries []dto.Tags
	var entities dto.Entities
	var collectionItemProps *dto.CollectionItemsProps
	var canSubscribe bool
	canSubscribe = false

	collectionItemProps = nil
	author := &models.Users{}
	if tag != nil {
		author, _ = s.userRepo.GetUserByLogin(tag.GetAuthor())
		var user *models.Users
		var subsTags []string
		if userAuthLogin != "" {
			user, _ = s.userRepo.GetUserByLogin(userAuthLogin)
			subsTags = user.TagsSubscriptions
		}

		canSubscribe = !slices.Contains(subsTags, tag.Id)

		if tag.GetPrimaryId() != "" {
			primaryTag = s.convertProtoToModel(tag.Primary, userAuthLogin, isPremium, filter, tagsMap)
		} else {
			primaryTag = nil
		}

		if tag.GetSeriesId() != "" {
			seriesTag = &dto.ShortTags{
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
			synonyms = []dto.Tags{}
		}
		if tag.GetSeriesEntries() != nil || len(tag.GetSeriesEntries()) > 0 {
			for _, entry := range tag.GetSeriesEntries() {
				seriesEntries = append(seriesEntries, *s.convertProtoToModel(entry, userAuthLogin, isPremium, filter, tagsMap))
			}
		} else {
			seriesEntries = []dto.Tags{}
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
			entities = dto.Entities{}
		}

		var authorTag *dto.SubUsers
		authorTag = &dto.SubUsers{}
		if author != nil {
			authorTag = &dto.SubUsers{
				Login:       author.Login,
				AvatarURL:   author.AvatarURL,
				ProfileName: author.ProfileName,
				IsPremium:   author.IsPremium,
			}
		}
		addFields := utils.NormalizeContent(tag.AdditionalFields)
		return &dto.Tags{
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
			Author:               *authorTag,
			IsSeries:             tag.IsSeries,
			SeriesID:             tag.SeriesId,
			SeriesEntries:        seriesEntries,
			Series:               seriesTag,
			AdditionalFields:     addFields,
			TotalCollections:     tag.TotalCollections,
			TotalPosts:           tag.TotalPosts,
			TotalCollectionItems: tag.TotalCollectionItems,
			TotalPhotos:          tag.TotalPhotos,
			CanSubscribe:         canSubscribe,
		}
	}
	return &dto.Tags{}
}

func (s *TagsService) getEntitiesByType(entityType, authUserLogin string, entityIDs []string, isPremium bool, filter string, tagsMap map[string]*microservices.GetShortsResponse) (*dto.Entities, *dto.CollectionItemsProps, error) {
	var entities dto.Entities
	var collectionItemsProps *dto.CollectionItemsProps = &dto.CollectionItemsProps{}
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
