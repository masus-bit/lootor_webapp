package services

import (
	"context"
	"github.com/google/uuid"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/eventsclient"
	"lootor/internal/infrastructure/photosclient"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/utils"
	"strconv"
	"sync"
)

type EventsService struct {
	eventClient         *eventsclient.GRPCEventsClient
	userRepo            *repositories.UsersRepository
	collectionsRepo     *repositories.CollectionsRepository
	collectionItemsRepo *repositories.CiRepository
	postService         *postsclient.GRPCPostsClient
	wishListRepo        *repositories.WLRepository
	tagsClient          *tagsclient.GRPCTagsClient
	photosClient        *photosclient.GRPCPhotosClient
}

func NewEventsService(userRepo *repositories.UsersRepository, eventClient *eventsclient.GRPCEventsClient, collectionRepo *repositories.CollectionsRepository, collectionItemsRepo *repositories.CiRepository, postService *postsclient.GRPCPostsClient, wishListRepo *repositories.WLRepository, tagsClient *tagsclient.GRPCTagsClient, photosClient *photosclient.GRPCPhotosClient) *EventsService {
	return &EventsService{
		userRepo:            userRepo,
		eventClient:         eventClient,
		collectionsRepo:     collectionRepo,
		collectionItemsRepo: collectionItemsRepo,
		postService:         postService,
		wishListRepo:        wishListRepo,
		tagsClient:          tagsClient,
		photosClient:        photosClient,
	}
}

func (s *EventsService) AddEvent(userLogin, action, target, title string, params *models.EventsParams) error {
	ok, err := s.eventClient.AddEvent(context.Background(), &models.AddEventRequest{
		Action:          action,
		EventTargetType: target,
		TargetName:      title,
		InitiatorLogin:  userLogin,
		Params: models.EventsParamsStrings{
			TargetUserLogin:      params.TargetUserLogin,
			TargetCollectionID:   params.TargetCollectionID.String(),
			TargetItemID:         params.TargetItemID.String(),
			TargetWLID:           params.TargetWLID.String(),
			TargetPostID:         params.TargetPostID,
			TargetTagID:          params.TargetTagID.String(),
			TagRelatedEntityType: params.TagRelatedEntityType,
			TargetPhotoID:        params.TargetPhotoID,
		},
	})
	if err != nil {
		return err
	}
	if !ok.Success {
		return err
	}
	return nil
}

func (s *EventsService) GetEvents(authUserLogin, limit, offset string, eventTargetTypes, actions []string) (*models.EventsDataResponse, error) {
	dbUser, err := s.userRepo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}
	var events []models.Events
	var totalCount int64
	subscriptions := dbUser.Subscriptions
	if subscriptions != nil {
		evs, err := s.eventClient.GetEvents(context.Background(), &models.GetEventsRequest{
			Limit:            limit,
			Offset:           offset,
			Subscriptions:    subscriptions,
			Actions:          actions,
			EventTargetTypes: eventTargetTypes,
		})
		if err != nil {
			return nil, err
		}
		events = s.normalizeEvents(evs.GetItems(), authUserLogin, dbUser.IsPremium)
		totalCount = evs.Total
	}
	return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
}

func (s *EventsService) GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset, authUserLogin, tagId string) (*models.EventsDataResponse, error) {
	dbUser, _ := s.userRepo.GetUserByLogin(authUserLogin)
	var collectionItemIDs []string
	var err error
	if collectionId != "" {
		collectionItemIDs, err = s.collectionsRepo.GetCollectionItemsIds(collectionId)
		if err != nil {
			return nil, err
		}
	}
	evs, err := s.eventClient.GetFilteredEvents(context.Background(), &models.GetFilteredEventsRequest{
		UserLogin:         userLogin,
		CollectionID:      collectionId,
		CollectionItemID:  collectionItem,
		WishListItemID:    wlId,
		Limit:             limit,
		Offset:            offset,
		CollectionItemIDs: collectionItemIDs,
		TagID:             tagId,
	})
	if err != nil {
		return nil, err
	}
	var isPremium bool

	if dbUser != nil {
		isPremium = dbUser.IsPremium
	} else {
		isPremium = false
	}
	events := s.normalizeEvents(evs.GetItems(), authUserLogin, isPremium)
	return &models.EventsDataResponse{Data: events, Total: evs.Total}, nil
}

func (s *EventsService) normalizeEvents(events []*microservices.EventsItem, authUserLogin string, isPremium bool) []models.Events {

	var normalizedEvents []models.Events
	var userLogins []string
	var collectionIDs []string
	var itemIDs []string
	var wlIDs []string
	var postIDs []string
	var tagIDs []string
	var photoIDs []string

	for _, event := range events {
		if event.InitiatorLogin != "" {
			userLogins = append(userLogins, event.InitiatorLogin)
		}
		if event.TargetUserLogin != "" {
			userLogins = append(userLogins, event.TargetUserLogin)
		}
		if event.TargetCollectionId != "" {
			collectionIDs = append(collectionIDs, event.TargetCollectionId)
		}
		if event.TargetItemId != "" {
			itemIDs = append(itemIDs, event.TargetItemId)
		}
		if event.TargetWishListItemId != "" {
			wlIDs = append(wlIDs, event.TargetWishListItemId)
		}
		if event.TargetPostId != "" {
			postIDs = append(postIDs, event.TargetPostId)
		}
		if event.TargetTagId != "" {
			tagIDs = append(tagIDs, event.TargetTagId)
			switch event.TagRelatedEntityType {
			case "collection":
				collectionIDs = append(collectionIDs, event.TargetCollectionId)
			case "item":
				itemIDs = append(itemIDs, event.TargetItemId)
			case "post":
				postIDs = append(postIDs, event.TargetPostId)
			case "photo":
				photoIDs = append(photoIDs, event.TargetPhotoId)

			}
		}
		if event.TargetPhotoId != "" {
			photoIDs = append(photoIDs, event.TargetPhotoId)
		}
	}

	var wg sync.WaitGroup

	var usersMap map[string]models.SubUsers
	var collectionsMap map[string]models.CollectionsResponse
	var itemsMap map[string]models.CollectionItemsResponse
	var wlMap map[string]models.WishListItemResponse
	var postsMap map[string]models.Posts
	var tagsMap map[string]models.ShortTags
	var photosMap map[string]models.Photos

	var collectionsTagsMap, itemsTagsMap, postsTagsMap, photosTagsMap map[string][]models.ShortTags

	collectionsTagsMap = make(map[string][]models.ShortTags)
	itemsTagsMap = make(map[string][]models.ShortTags)
	postsTagsMap = make(map[string][]models.ShortTags)
	photosTagsMap = make(map[string][]models.ShortTags)

	wg.Add(7)

	go func() {
		defer wg.Done()
		usersMap, _ = s.userRepo.GetForSubsMap(userLogins)
	}()
	go func() {
		defer wg.Done()
		collectionIDsUUID := make([]uuid.UUID, len(collectionIDs))
		for i, id := range collectionIDs {
			collectionIDsUUID[i], _ = uuid.Parse(id)
		}
		tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
			EntityIds: collectionIDs,
		})
		if err != nil {
			return
		}
		for key, tag := range tags.GetTags() {
			collectionsTagsMap[key] = toShortTags(tag.GetTags())
		}
		counts, _ := s.collectionItemsRepo.GetCountCIByIDs(collectionIDsUUID)
		totalPrices, _ := s.collectionItemsRepo.GetSumsByCollectionIDs(collectionIDsUUID)
		shippingCosts, _ := s.collectionItemsRepo.GetShippingCostsByCollectionIDs(collectionIDsUUID)
		collectionsMap, _ = s.collectionsRepo.GetCollectionsByIdsMap(collectionIDs, counts, totalPrices, shippingCosts, authUserLogin)
	}()
	go func() {
		defer wg.Done()
		itemsMap, _ = s.collectionItemsRepo.GetCollectionItemsByIdsMap(itemIDs, authUserLogin)
		tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
			EntityIds: itemIDs,
		})
		if err != nil {
			return
		}
		for key, tag := range tags.GetTags() {
			itemsTagsMap[key] = toShortTags(tag.GetTags())
		}
	}()
	go func() {
		defer wg.Done()
		wlMap, _ = s.wishListRepo.GetWLByIDsMap(wlIDs)
	}()
	go func() {
		defer wg.Done()
		photos, _ := s.photosClient.GetPhotosByIDsMap(context.Background(), &microservices.GetByIdsRequest{
			Ids: photoIDs,
		})
		var authors []string
		var collectionsPhotoIds []string
		var collectionsPhotoMap map[string]models.CollectionShort
		for _, ph := range photos.GetData() {
			authors = append(authors, ph.GetAuthor())
			collectionsPhotoIds = append(collectionsPhotoIds, ph.GetCollectionId())
		}
		collectionsPhotoMap, _ = s.collectionsRepo.GetCollectionsByIdsMapForShort(collectionsPhotoIds)
		authorsShort, _ := s.userRepo.GetForSubsMap(authors)
		photosMap = make(map[string]models.Photos)
		if photos != nil && photos.GetData() != nil {
			for key, photo := range photos.GetData() {
				createdAtAsTime, err := parseDate(photo.GetCreatedAt())
				if err != nil {
					log.Printf("Error parsing date: %v", err)
				}
				photosMap[key] = models.Photos{
					ID:            photo.GetId(),
					Author:        authorsShort[photo.GetAuthor()],
					CollectionID:  photo.GetCollectionId(),
					Path:          photo.GetPath(),
					LikesCount:    int64(len(photo.GetLikes())),
					CommentsCount: photo.GetCommentsCount(),
					CreatedAt:     createdAtAsTime,
					Collection:    collectionsPhotoMap[photo.GetCollectionId()],
				}
			}
		}
		tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
			EntityIds: itemIDs,
		})
		if err != nil {
			return
		}
		for key, tag := range tags.GetTags() {
			photosTagsMap[key] = toShortTags(tag.GetTags())
		}
	}()
	go func() {
		defer wg.Done()
		resp, _ := s.postService.GetPostsByIds(context.Background(), postIDs, authUserLogin, isPremium)
		postsMap = make(map[string]models.Posts)
		if resp != nil && resp.Data != nil {
			for key, p := range resp.Data {
				var reactUsers []string

				for _, r := range p.Reactions {
					reactUsers = append(reactUsers, r.UserLogin)
				}
				reactsLen := len(p.GetReactions())
				var users []models.SubUsers
				if reactsLen != 0 {
					users, _ = s.userRepo.GetForSubs(reactUsers)
				}
				author, _ := s.userRepo.GetUserByLogin(p.Author)
				postFormatted, err := utils.FormatPost(p, users, author, reactsLen)
				if err != nil {
					continue
				}

				postsMap[key] = models.Posts{
					ID:             postFormatted.Data.ID,
					Title:          postFormatted.Data.Title,
					Author:         postFormatted.Data.Author,
					Translit:       postFormatted.Data.Translit,
					IsDraft:        postFormatted.Data.IsDraft,
					Content:        postFormatted.Data.Content,
					Views:          postFormatted.Data.Views,
					Date:           postFormatted.Data.Date,
					CommentsCount:  postFormatted.Data.CommentsCount,
					HeartCount:     postFormatted.Data.HeartCount,
					FireCount:      postFormatted.Data.FireCount,
					GlassesCount:   postFormatted.Data.GlassesCount,
					LaughCount:     postFormatted.Data.LaughCount,
					TearsCount:     postFormatted.Data.TearsCount,
					PokerFaceCount: postFormatted.Data.PokerFaceCount,
					EyesCount:      postFormatted.Data.EyesCount,
					AngryCount:     postFormatted.Data.AngryCount,
					ShitCount:      postFormatted.Data.ShitCount,
					ClownCount:     postFormatted.Data.ClownCount,
					TotalReactions: postFormatted.Data.TotalReactions,
					Reacted:        postFormatted.Data.Reacted,
				}
			}
		}
		go func() {
			defer wg.Done()
			tagsRaw, err := s.tagsClient.GetTagsByIdsMap(context.Background(), &microservices.GetTagsByIDsRequest{
				Ids: tagIDs,
			})
			if err != nil {
				return
			}
			if tagsRaw != nil {
				for key, tag := range tagsRaw.GetTags() {
					tagsMap[key] = models.ShortTags{
						ID:        tag.GetId(),
						Name:      tag.GetName(),
						Slug:      tag.GetSlug(),
						PrimaryID: tag.GetPrimaryId(),
						SeriesID:  tag.GetSeriesId(),
					}
				}
			}

		}()
		tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
			EntityIds: postIDs,
		})
		if err != nil {
			return
		}
		for key, tag := range tags.GetTags() {
			postsTagsMap[key] = toShortTags(tag.GetTags())
		}
	}()

	wg.Wait()

	if usersMap == nil {
		usersMap = make(map[string]models.SubUsers)
	}
	if collectionsMap == nil {
		collectionsMap = make(map[string]models.CollectionsResponse)
	}
	if itemsMap == nil {
		itemsMap = make(map[string]models.CollectionItemsResponse)
	}
	if wlMap == nil {
		wlMap = make(map[string]models.WishListItemResponse)
	}
	if postsMap == nil {
		postsMap = make(map[string]models.Posts)
	}
	if tagsMap == nil {
		tagsMap = make(map[string]models.ShortTags)
	}
	if photosMap == nil {
		photosMap = make(map[string]models.Photos)
	}

	for _, event := range events {
		normalizedEvent := models.Events{
			ID:                   event.Id,
			Date:                 event.Date,
			Action:               event.Action,
			EventTargetType:      event.EventTargetType,
			TargetName:           event.TargetName,
			InitiatorLogin:       event.InitiatorLogin,
			TargetUserLogin:      event.TargetUserLogin,
			TargetCollectionID:   event.TargetCollectionId,
			TargetItemID:         event.TargetItemId,
			TargetWishListItemID: event.TargetWishListItemId,
			TargetPostID:         event.TargetPostId,
			TargetTagID:          event.TargetTagId,
			TagRelatedEntityType: event.TagRelatedEntityType,
			TargetPhotoID:        event.TargetPhotoId,
		}

		if user, exists := usersMap[event.TargetUserLogin]; exists {
			normalizedEvent.TargetUser = &models.SubUsers{
				Login:       user.Login,
				AvatarURL:   user.AvatarURL,
				ProfileName: user.ProfileName,
				IsPremium:   user.IsPremium,
			}
		}

		if collection, exists := collectionsMap[event.TargetCollectionId]; exists {
			collection.Tags = collectionsTagsMap[collection.ID.String()]
			normalizedEvent.TargetCollection = &collection
		}

		if item, exists := itemsMap[event.TargetItemId]; exists {
			item.Tags = itemsTagsMap[item.ID.String()]
			normalizedEvent.TargetItem = &item
		}

		if wlItem, exists := wlMap[event.TargetWishListItemId]; exists {
			normalizedEvent.TargetWishListItem = &wlItem
		}

		if post, exists := postsMap[event.TargetPostId]; exists {
			post.Tags = postsTagsMap[strconv.FormatUint(post.ID, 10)]
			normalizedEvent.TargetPost = &post
		}

		if photo, exists := photosMap[event.TargetPhotoId]; exists {
			photo.Tags = photosTagsMap[strconv.FormatUint(photo.ID, 10)]
			normalizedEvent.TargetPhoto = &photo
		}

		if initiator, exists := usersMap[event.InitiatorLogin]; exists {
			normalizedEvent.Initiator = &models.SubUsers{
				Login:       initiator.Login,
				AvatarURL:   initiator.AvatarURL,
				ProfileName: initiator.ProfileName,
				IsPremium:   initiator.IsPremium,
			}
		}

		if tag, exists := tagsMap[event.TargetTagId]; exists {
			normalizedEvent.TargetTag = &tag
			switch normalizedEvent.TagRelatedEntityType {
			case "collection":
				collectionRes := collectionsMap[normalizedEvent.TargetCollectionID]
				normalizedEvent.TargetCollection = &collectionRes
			case "collectionItem":
				itemRes := itemsMap[normalizedEvent.TargetItemID]
				normalizedEvent.TargetItem = &itemRes
			case "post":
				postRes := postsMap[normalizedEvent.TargetPostID]
				normalizedEvent.TargetPost = &postRes
			case "photos":
				photoRes := photosMap[normalizedEvent.TargetPhotoID]
				normalizedEvent.TargetPhoto = &photoRes
			}
		}

		normalizedEvents = append(normalizedEvents, normalizedEvent)
	}
	return normalizedEvents
}

func toShortTags(tags []*microservices.TagItemShort) []models.ShortTags {
	result := make([]models.ShortTags, len(tags))
	for i, t := range tags {
		result[i] = models.ShortTags{ID: t.GetId(), Name: t.GetName(), Slug: t.GetSlug(), PrimaryID: t.GetPrimaryId(), SeriesID: t.GetSeriesId()}
	}
	return result
}
