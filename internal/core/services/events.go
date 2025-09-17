package services

import (
	"context"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/eventsclient"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/types"
	"sync"
)

type EventsService struct {
	eventClient         *eventsclient.GRPCEventsClient
	userRepo            *repositories.UsersRepository
	collectionsRepo     *repositories.CollectionsRepository
	collectionItemsRepo *repositories.CiRepository
	postService         *postsclient.GRPCPostsClient
	wishListRepo        *repositories.WLRepository
}

func NewEventsService(userRepo *repositories.UsersRepository, eventClient *eventsclient.GRPCEventsClient, collectionRepo *repositories.CollectionsRepository, collectionItemsRepo *repositories.CiRepository, postService *postsclient.GRPCPostsClient, wishListRepo *repositories.WLRepository) *EventsService {
	return &EventsService{
		userRepo:            userRepo,
		eventClient:         eventClient,
		collectionsRepo:     collectionRepo,
		collectionItemsRepo: collectionItemsRepo,
		postService:         postService,
		wishListRepo:        wishListRepo,
	}
}

func (s *EventsService) AddEvent(userLogin, action, target, title string, params *models.EventsParams) error {
	ok, err := s.eventClient.AddEvent(context.Background(), &models.AddEventRequest{
		Action:          action,
		EventTargetType: target,
		TargetName:      title,
		InitiatorLogin:  userLogin,
		Params: models.EventsParamsStrings{
			TargetUserLogin:    params.TargetUserLogin,
			TargetCollectionID: params.TargetCollectionID.String(),
			TargetItemID:       params.TargetItemID.String(),
			TargetWLID:         params.TargetWLID.String(),
			TargetPostID:       params.TargetPostID,
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

func (s *EventsService) GetEvents(authUserLogin, limit, offset string) (*models.EventsDataResponse, error) {
	dbUser, err := s.userRepo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}
	var events []models.Events
	var totalCount int64
	subscriptions := dbUser.Subscriptions
	if subscriptions != nil {
		evs, err := s.eventClient.GetEvents(context.Background(), &models.GetEventsRequest{
			Limit:         limit,
			Offset:        offset,
			Subscriptions: subscriptions,
		})
		if err != nil {
			return nil, err
		}
		events = s.normalizeEvents(evs.Items)
		totalCount = evs.Total
	}
	return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
}

func (s *EventsService) GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset string) (*models.EventsDataResponse, error) {
	evs, err := s.eventClient.GetFilteredEvents(context.Background(), &models.GetFilteredEventsRequest{
		UserLogin:        userLogin,
		CollectionId:     collectionId,
		CollectionItemId: collectionItem,
		WishListItemId:   wlId,
		Limit:            limit,
		Offset:           offset,
	})
	if err != nil {
		return nil, err
	}
	events := s.normalizeEvents(evs.GetItems())
	return &models.EventsDataResponse{Data: events, Total: evs.Total}, nil
}

func (s *EventsService) normalizeEvents(events []*microservices.EventsItem) []models.Events {

	var normalizedEvents []models.Events
	var userLogins []string
	var collectionIDs []string
	var itemIDs []string
	var wlIDs []string
	var postIDs []string

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
	}

	var wg sync.WaitGroup

	var usersMap map[string]models.SubUsers
	var collectionsMap map[string]models.Collections
	var itemsMap map[string]models.CollectionItems
	var wlMap map[string]models.WishListItems
	var postsMap map[string]dto.ShortPost

	wg.Add(5)

	go func() {
		defer wg.Done()
		usersMap, _ = s.userRepo.GetForSubsMap(userLogins)
	}()
	go func() {
		defer wg.Done()
		collectionsMap, _ = s.collectionsRepo.GetCollectionsByIdsMap(collectionIDs)
	}()
	go func() {
		defer wg.Done()
		itemsMap, _ = s.collectionItemsRepo.GetCollectionItemsByIdsMap(itemIDs)
	}()
	go func() {
		defer wg.Done()
		wlMap, _ = s.wishListRepo.GetWLByIdsMap(wlIDs)
	}()
	go func() {
		defer wg.Done()
		resp, _ := s.postService.GetPostsByIds(context.Background(), postIDs)
		postsMap = make(map[string]dto.ShortPost)
		if resp != nil && resp.Data != nil {
			for key, p := range resp.Data {
				postsMap[key] = dto.ShortPost{
					Id:       p.Id,
					Author:   p.Author,
					Title:    p.Title,
					Translit: p.Translit,
				}
			}
		}
	}()

	wg.Wait()

	if usersMap == nil {
		usersMap = make(map[string]models.SubUsers)
	}
	if collectionsMap == nil {
		collectionsMap = make(map[string]models.Collections)
	}
	if itemsMap == nil {
		itemsMap = make(map[string]models.CollectionItems)
	}
	if wlMap == nil {
		wlMap = make(map[string]models.WishListItems)
	}
	if postsMap == nil {
		postsMap = make(map[string]dto.ShortPost)
	}

	for _, event := range events {
		normalizedEvent := models.Events{
			Id:                   event.Id,
			Date:                 event.Date,
			Action:               event.Action,
			EventTargetType:      event.EventTargetType,
			TargetName:           event.TargetName,
			InitiatorLogin:       event.InitiatorLogin,
			TargetUserLogin:      event.TargetUserLogin,
			TargetCollectionID:   event.TargetCollectionId,
			TargetItemID:         event.TargetItemId,
			TargetWishListItemID: event.TargetWishListItemId,
			TargetPostId:         event.TargetPostId,
		}

		if user, exists := usersMap[event.TargetUserLogin]; exists {
			normalizedEvent.TargetUser = &models.SubUsers{
				Login:       user.Login,
				AvatarUrl:   user.AvatarUrl,
				ProfileName: user.ProfileName,
				IsPremium:   user.IsPremium,
			}
		}

		if collection, exists := collectionsMap[event.TargetCollectionId]; exists {
			normalizedEvent.TargetCollection = &types.CommonShortTypeCollection{
				CommonShortType: &types.CommonShortType{
					ID:   collection.Id.String(),
					Name: collection.Name,
				},
				Owner:           collection.UserLogin,
				Transliteration: collection.Transliteration,
			}
		}

		if item, exists := itemsMap[event.TargetItemId]; exists {
			targetItem := &types.CommonShortTypeItem{
				CommonShortType: &types.CommonShortType{
					ID:   item.Id.String(),
					Name: item.Name,
				},
				Owner: item.Owner.Login,
			}

			if len(item.Collections) > 0 {
				targetItem.Collection = item.Collections[0].Id.String()
				targetItem.CollectionName = item.Collections[0].Name
			} else {
				targetItem.Collection = ""
				targetItem.CollectionName = ""
			}

			normalizedEvent.TargetItem = targetItem
		}

		if wlItem, exists := wlMap[event.TargetWishListItemId]; exists {
			normalizedEvent.TargetWishListItem = &types.CommonShortType{
				ID:   wlItem.Id.String(),
				Name: wlItem.ItemName,
			}
		}

		if post, exists := postsMap[event.TargetPostId]; exists {
			normalizedEvent.TargetPost = &types.CommonShortTypePost{
				CommonShortType: &types.CommonShortType{
					ID:   post.Id,
					Name: post.Title,
				},
				Author:          post.Author,
				Transliteration: post.Translit,
			}
		}

		if initiator, exists := usersMap[event.InitiatorLogin]; exists {
			normalizedEvent.Initiator = &models.SubUsers{
				Login:       initiator.Login,
				AvatarUrl:   initiator.AvatarUrl,
				ProfileName: initiator.ProfileName,
				IsPremium:   initiator.IsPremium,
			}
		}

		normalizedEvents = append(normalizedEvents, normalizedEvent)
	}
	return normalizedEvents
}
