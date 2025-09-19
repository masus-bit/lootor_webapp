package services

import (
	"context"
	"github.com/google/uuid"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/eventsclient"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/pkg/utils"
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
		events = s.normalizeEvents(evs.Items, authUserLogin, dbUser.IsPremium)
		totalCount = evs.Total
	}
	return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
}

func (s *EventsService) GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset, authUserLogin string) (*models.EventsDataResponse, error) {
	dbUser, err := s.userRepo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}
	var collectionItemIDs []string
	if collectionId != "" {
		collectionItemIDs, err = s.collectionsRepo.GetCollectionItemsIds(collectionId)
		if err != nil {
			return nil, err
		}
	}
	evs, err := s.eventClient.GetFilteredEvents(context.Background(), &models.GetFilteredEventsRequest{
		UserLogin:         userLogin,
		CollectionId:      collectionId,
		CollectionItemId:  collectionItem,
		WishListItemId:    wlId,
		Limit:             limit,
		Offset:            offset,
		CollectionItemIds: collectionItemIDs,
	})
	if err != nil {
		return nil, err
	}
	events := s.normalizeEvents(evs.GetItems(), authUserLogin, dbUser.IsPremium)
	return &models.EventsDataResponse{Data: events, Total: evs.Total}, nil
}

func (s *EventsService) normalizeEvents(events []*microservices.EventsItem, authUserLogin string, isPremium bool) []models.Events {

	authUser, err := s.userRepo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil
	}

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
	var collectionsMap map[string]models.CollectionsResponse
	var itemsMap map[string]models.CollectionItemsResponse
	var wlMap map[string]models.WishListItemResponse
	var postsMap map[string]models.Posts

	wg.Add(5)

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
		counts, _ := s.collectionItemsRepo.GetCountCIByIDs(collectionIDsUUID)
		totalPrices, _ := s.collectionItemsRepo.GetSumsByCollectionIDs(collectionIDsUUID)
		shippingCosts, _ := s.collectionItemsRepo.GetShippingCostsByCollectionIDs(collectionIDsUUID)
		var subArray []string
		if authUserLogin != "" {
			subArray = authUser.CollectionSubscriptions
		}
		collectionsMap, _ = s.collectionsRepo.GetCollectionsByIdsMap(collectionIDs, counts, totalPrices, shippingCosts, subArray, authUserLogin)
	}()
	go func() {
		defer wg.Done()
		itemsMap, _ = s.collectionItemsRepo.GetCollectionItemsByIdsMap(itemIDs, authUserLogin)
	}()
	go func() {
		defer wg.Done()
		wlMap, _ = s.wishListRepo.GetWLByIdsMap(wlIDs)
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
					Id:             postFormatted.Data.Id,
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
			normalizedEvent.TargetCollection = &collection
		}

		if item, exists := itemsMap[event.TargetItemId]; exists {
			normalizedEvent.TargetItem = &item
		}

		if wlItem, exists := wlMap[event.TargetWishListItemId]; exists {
			normalizedEvent.TargetWishListItem = &wlItem
		}

		if post, exists := postsMap[event.TargetPostId]; exists {
			normalizedEvent.TargetPost = &post
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
