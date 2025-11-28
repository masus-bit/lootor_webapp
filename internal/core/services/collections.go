package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/photosclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/s3"
	"lootor/internal/pkg/utils"
	"reflect"
	"slices"
	"strings"
	"time"
)

type CollectionService struct {
	repo                 *repositories.CollectionsRepository
	userRepo             *repositories.UsersRepository
	eventsService        *EventsService
	collectionItemRepo   *repositories.CiRepository
	s3Service            *s3.StorageService
	notificationsService *NotificationsService
	tagsClient           *tagsclient.GRPCTagsClient
	photosClient         *photosclient.GRPCPhotosClient
}

func NewCollectionService(repo *repositories.CollectionsRepository, userRepo *repositories.UsersRepository, eventsService *EventsService, collectionItemRepo *repositories.CiRepository, seService *s3.StorageService, notificationsService *NotificationsService, tagsClient *tagsclient.GRPCTagsClient, photosClient *photosclient.GRPCPhotosClient) *CollectionService {
	return &CollectionService{repo: repo, userRepo: userRepo, eventsService: eventsService, collectionItemRepo: collectionItemRepo, s3Service: seService, notificationsService: notificationsService, tagsClient: tagsClient, photosClient: photosClient}
}

func (s *CollectionService) Create(dto *models.CollectionCreateRequest) (*models.CollectionDataResponse, error) {
	existsCollection, err := s.repo.GetUniqueByName(dto.Name, dto.UserLogin)
	if err != nil {
		fmt.Println(err)
	}
	if existsCollection != nil {
		return nil, errors.New("collection с таким именем уже существует")
	}

	dtoUser, _ := s.userRepo.GetUserByLogin(dto.UserLogin)
	dbCollection := &models.Collections{
		Name:            dto.Name,
		Description:     dto.Description,
		BannerURL:       dto.BannerURL,
		IsPrivate:       dto.IsPrivate,
		Transliteration: dto.Transliteration,
		User:            dtoUser,
		Created:         time.Now().Format(time.RFC3339),
		ShippingTotal:   0,
		TotalPrice:      0,
		ShareString:     utils.GenerateRandomStringNoHex(8),
		UserLogin:       dto.UserLogin,
	}

	res, err := s.repo.CreateCollection(dbCollection)

	if err != nil {
		return nil, err
	}

	err = s.userRepo.IncrementExperience(dto.UserLogin, utils.CollectionExp)
	if err != nil {
		return nil, err
	}

	if !dto.IsPrivate {
		eventError := s.eventsService.AddEvent(dto.UserLogin, utils.EventActionCreate, utils.EventTargetCollection, dto.Name, &models.EventsParams{TargetCollectionID: res.ID})
		if eventError != nil {
			log.Default().Print(eventError)
		}

	}

	var response models.CollectionsResponse

	e := mapstructure.Decode(res, &response)
	response.ID = res.ID

	var tags *microservices.GetShortsResponse

	if len(dto.Tags) > 0 {
		tags, _ = s.tagsClient.AddTagsToEntity(context.Background(), &microservices.AddFewTagsToEntityRequest{
			EntityType: "collection",
			EntityId:   res.ID.String(),
			TagIds:     dto.Tags,
			Author:     dto.UserLogin,
			ShowSearch: !dto.IsPrivate,
		})
		if !dto.IsPrivate {
			for _, tag := range tags.GetTags() {
				tagUUID, _ := uuid.Parse(tag.GetId())
				eventError := s.eventsService.AddEvent(dto.UserLogin, utils.EventActionAddTag, utils.EventTargetTag, tag.Name, &models.EventsParams{TargetTagID: tagUUID, TargetCollectionID: res.ID, TagRelatedEntityType: utils.EventTargetCollection})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}
		}
	}

	var resultTags []models.ShortTags
	resultTags = make([]models.ShortTags, 0)
	if tags != nil {
		for _, tag := range tags.Tags {
			resultTags = append(resultTags, models.ShortTags{
				ID:   tag.Id,
				Name: tag.Name,
				Slug: tag.Slug,
			})
		}
		response.Tags = resultTags
	}

	if e != nil {
		return nil, e
	}

	er := mapstructure.Decode(dtoUser, &response.User)

	if er != nil {
		return nil, er
	}
	response.CanLike = true
	response.IsOwner = true
	return &models.CollectionDataResponse{Data: response}, nil
}

func (s *CollectionService) Update(id string, dto *models.CollectionUpdateRequest) (*models.CollectionDataResponse, error) {
	exists, err := s.repo.GetCollectionByIdWithoutLimits(id)
	if err != nil {
		return nil, err
	}

	dst := reflect.ValueOf(exists).Elem()
	src := reflect.ValueOf(dto).Elem()

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		if !field.IsNil() {
			fieldName := src.Type().Field(i).Name
			if fieldName == "Tags" || fieldName == "UserLogin" || fieldName == "IsPrivate" {
				continue
			}
			dstField := dst.FieldByName(fieldName)
			if dstField.IsValid() {
				dstField.Set(field.Elem())
			}
		}
	}

	exists.IsPrivate = *dto.IsPrivate

	resultCollection, err := s.repo.UpdateCollectionFull(exists)
	if err != nil {
		return nil, err
	}

	if !*dto.IsPrivate {
		eventError := s.eventsService.AddEvent(
			*dto.UserLogin,
			utils.EventActionUpdate,
			utils.EventTargetCollection,
			*dto.Name,
			&models.EventsParams{TargetCollectionID: resultCollection.ID},
		)
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}

	var finalCollection models.CollectionsResponse
	if err = mapstructure.Decode(resultCollection, &finalCollection); err != nil {
		return nil, err
	}

	var tags *microservices.TagsCreateResponse
	tags, err = s.tagsClient.UpdateTagsOfEntity(context.Background(), &microservices.UpdateTagsOfEntityRequest{
		EntityType: "collection",
		EntityId:   id,
		TagIds:     dto.Tags,
		Author:     exists.UserLogin,
	})
	if err != nil {
		return nil, err
	}
	if len(dto.Tags) > 0 {
		if !*dto.IsPrivate {
			var entitiesIDs []string
			collectionItemsIDs, err := s.repo.GetCollectionItemsIds(id)
			if err != nil {
				return nil, err
			}
			entitiesIDs = append(entitiesIDs, collectionItemsIDs...)
			entitiesIDs = append(entitiesIDs, id)
			_, err = s.tagsClient.UpdateVisibleLinks(context.Background(), &microservices.UpdateVisibleLinksRequest{
				EntityIds: entitiesIDs,
				Visible:   true,
			})
			if err != nil {
				return nil, err
			}
			for _, tag := range tags.GetTags() {
				tagUUID, _ := uuid.Parse(tag.GetId())
				idUUID, _ := uuid.Parse(id)
				eventError := s.eventsService.AddEvent(exists.UserLogin, utils.EventActionAddTag, utils.EventTargetTag, tag.Name, &models.EventsParams{TargetTagID: tagUUID, TargetCollectionID: idUUID, TagRelatedEntityType: utils.EventTargetCollection})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}

		} else {
			var entitiesIDs []string
			collectionItemsIDs, err := s.repo.GetCollectionItemsIds(id)
			if err != nil {
				return nil, err
			}
			entitiesIDs = append(entitiesIDs, collectionItemsIDs...)
			entitiesIDs = append(entitiesIDs, id)
			_, err = s.tagsClient.UpdateVisibleLinks(context.Background(), &microservices.UpdateVisibleLinksRequest{
				EntityIds: entitiesIDs,
				Visible:   false,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	var resultTags []models.ShortTags
	resultTags = make([]models.ShortTags, 0)
	if tags != nil {
		for _, tag := range tags.GetTags() {
			resultTags = append(resultTags, models.ShortTags{
				ID:        tag.GetId(),
				Name:      tag.GetName(),
				Slug:      tag.GetSlug(),
				PrimaryID: tag.GetPrimaryId(),
				SeriesID:  tag.GetSeriesId(),
			})
		}
	}

	finalCollection.Tags = resultTags
	finalCollection.CanLike = true
	finalCollection.IsOwner = true

	return &models.CollectionDataResponse{Data: finalCollection}, nil
}

func (s *CollectionService) Delete(id string, ctx context.Context) (*dto.CommonResponse, error) {
	exists, _ := s.repo.GetCollectionByIdWithoutLimits(id)
	if !exists.IsPrivate {
		go func() {
			eventError := s.eventsService.AddEvent(exists.User.Login, utils.EventActionDelete, utils.EventTargetCollection, exists.Name, &models.EventsParams{TargetCollectionID: exists.ID})
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}()
		if exists.BannerURL != "" {
			_, err := s.s3Service.DeleteFiles(ctx, s3.DeleteFilesRequest{Keys: []string{exists.BannerURL}})
			if err != nil {
				return nil, err
			}
		}
		go func() {
			_ = s.notificationsService.DeleteAllNotificationsByTargetID(context.Background(), id)
		}()

	}

	go func() {
		_, _ = s.tagsClient.RemoveEntityTags(context.Background(), &microservices.RemoveEntityTagsRequest{
			EntityId: id,
		})
	}()

	exp := len(exists.CollectionItems) + utils.CollectionExp

	cis := exists.CollectionItems
	for _, ci := range cis {
		if ci.PurchaseDate != "" {
			exp += utils.CIPurchaseDateExp
		}
		if ci.PurchasePrice != 0 {
			exp += utils.CIPurchasePriceExp
		}
		if ci.Rating != 0 {
			exp += utils.CIRatingExp
		}
		if len(ci.Images) != 0 {
			exp += utils.PictureExp
		}
		if len(ci.CopyNumber) != 0 {
			exp += utils.CICopyNumberExp
		}
	}
	err := s.userRepo.DecrementExperience(exists.UserLogin, exp)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, s.repo.DeleteCollection(id)
}

func (s *CollectionService) GetByUserLogin(login string, authorizedUser string, orderBy string, order string, search string) (*models.AllCollectionsDataResponse, error) {
	var collections []models.Collections
	var total int64
	if authorizedUser == login {

		collections, total, _ = s.repo.GetByUserIdWithoutCollectionItems(login, search)
	} else {
		collections, total, _ = s.repo.GetByUserIdWithoutPrivates(login, search)
	}

	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	result := make([]models.CollectionsResponse, 0)

	collectionIds := make([]uuid.UUID, 0)
	collectionStringIds := make([]string, 0)
	for _, dbCollection := range collections {
		collectionIds = append(collectionIds, dbCollection.ID)
		collectionStringIds = append(collectionStringIds, dbCollection.ID.String())
	}

	tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{EntityIds: collectionStringIds})
	if err != nil {
		return nil, err
	}
	counts, _ := s.collectionItemRepo.GetCountCIByIDs(collectionIds)
	totalPrices, _ := s.collectionItemRepo.GetSumsByCollectionIDs(collectionIds)
	shippingCosts, _ := s.collectionItemRepo.GetShippingCostsByCollectionIDs(collectionIds)
	photosCounts, _ := s.photosClient.GetCountsByCollectionsIDsMap(context.Background(), &microservices.GetCountsByCollectionsIdsMapRequest{CollectionsIds: collectionStringIds})

	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		err = mapstructure.Decode(dbCollection, &temp)

		itemTags := tags.GetTags()[dbCollection.ID.String()]

		var resultTags []models.ShortTags = utils.NormalizeTagsShort(itemTags.GetTags())

		temp.Tags = resultTags

		temp.ShareString = utils.DefineShareString(authorizedUser, login, &dbCollection)
		temp.CollectionItemsCount = counts[dbCollection.ID]
		temp.PhotosCount = photosCounts.GetData()[dbCollection.ID.String()]
		temp.TotalPrice = totalPrices[dbCollection.ID]
		temp.ShippingTotal = shippingCosts[dbCollection.ID]
		temp.LikesCount = int64(len(dbCollection.Likes))
		temp.CanLike = utils.CanLike(dbCollection.Likes, authorizedUser, dbCollection.UserLogin)
		temp.IsOwner = authorizedUser == login

		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}

	sortedCollections := utils.GetCollectionOrderBy(orderBy, result, order)

	return &models.AllCollectionsDataResponse{Data: sortedCollections, Total: total, ProfileName: user.ProfileName}, nil
}

func (s *CollectionService) GetAll(authorizedUser, orderBy, order, search, limit, offset string) (*models.AllCollectionsDataResponse, error) {
	var collections []models.Collections
	collections, total, _ := s.repo.GetAllWithoutPrivates(search, limit, offset, orderBy, order)
	result := make([]models.CollectionsResponse, 0)

	collectionIds := make([]uuid.UUID, 0)
	collectionStringIds := make([]string, 0)
	for _, dbCollection := range collections {
		collectionIds = append(collectionIds, dbCollection.ID)
		collectionStringIds = append(collectionStringIds, dbCollection.ID.String())
	}

	tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{EntityIds: collectionStringIds})
	if err != nil {
		return nil, err
	}
	photosCounts, _ := s.photosClient.GetCountsByCollectionsIDsMap(context.Background(), &microservices.GetCountsByCollectionsIdsMapRequest{CollectionsIds: collectionStringIds})
	counts, _ := s.collectionItemRepo.GetCountCIByIDs(collectionIds)
	totalPrices, _ := s.collectionItemRepo.GetSumsByCollectionIDs(collectionIds)
	shippingCosts, _ := s.collectionItemRepo.GetShippingCostsByCollectionIDs(collectionIds)

	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		err = mapstructure.Decode(dbCollection, &temp)

		itemTags := tags.GetTags()[dbCollection.ID.String()]

		var resultTags []models.ShortTags = utils.NormalizeTagsShort(itemTags.GetTags())

		temp.Tags = resultTags

		temp.CollectionItemsCount = counts[dbCollection.ID]
		temp.TotalPrice = totalPrices[dbCollection.ID]
		temp.ShippingTotal = shippingCosts[dbCollection.ID]
		temp.LikesCount = int64(len(dbCollection.Likes))
		temp.CanLike = utils.CanLike(dbCollection.Likes, authorizedUser, dbCollection.UserLogin)
		temp.IsOwner = authorizedUser == temp.User.Login
		temp.PhotosCount = photosCounts.GetData()[dbCollection.ID.String()]

		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}

	return &models.AllCollectionsDataResponse{Data: result, Total: total}, nil
}

func (s *CollectionService) GetOne(authorizerUser string, id string, transliteration string, userLogin string, shareString string, ciLimit string, ciOffset string, orderBy string, order string, search string) (*models.CollectionDataResponse, error) {
	var dbCollection *models.Collections

	if id != "" {
		dbCollection, _ = s.repo.GetCollectionById(id, ciLimit, ciOffset, orderBy, order, search)
	} else if userLogin != "" {
		dbCollection, _ = s.repo.GetOneByTransliteration(userLogin, transliteration, ciLimit, ciOffset, orderBy, order, search)
	} else if shareString != "" {
		dbCollection, _ = s.repo.GetByShareString(shareString, ciLimit, ciOffset, orderBy, order, search)
	}

	if dbCollection == nil {
		return nil, errors.New("collection not found")
	}

	var authUser *models.Users
	var subArray []string

	if authorizerUser != "" {
		authUser, _ = s.userRepo.GetUserByLogin(authorizerUser)
	}
	var finalCollection models.CollectionsResponse
	err := mapstructure.Decode(dbCollection, &finalCollection)
	if err != nil {
		return nil, err
	}

	collectionItems := make([]models.CollectionItemsResponse, 0)

	collectionItemIds := make([]string, 0)
	for _, item := range dbCollection.CollectionItems {
		collectionItemIds = append(collectionItemIds, item.ID.String())
	}

	ciTags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{EntityIds: collectionItemIds})
	if err != nil {
		return nil, err
	}
	photosCount, err := s.photosClient.GetCountByCollection(context.Background(), &microservices.CountRequestPhoto{CollectionId: dbCollection.ID.String()})
	if err != nil {
		return nil, err
	}

	for _, item := range dbCollection.CollectionItems {
		var temp models.CollectionItemsResponse
		itemErr := mapstructure.Decode(item, &temp)
		if itemErr != nil {
			return nil, itemErr
		}

		itemTags := ciTags.GetTags()[item.ID.String()]

		var resultTags []models.ShortTags = utils.NormalizeTagsShort(itemTags.GetTags())
		temp.Tags = resultTags

		temp.Collection = finalCollection.ID
		temp.LikesCount = int64(len(item.Likes))
		temp.CanLike = utils.CanLike(item.Likes, authorizerUser, item.Owner.Login)
		temp.CollectionTransliteration = finalCollection.Transliteration
		if authorizerUser != "" {
			if authorizerUser == item.Owner.Login {
				temp.IsOwner = true
			} else {
				temp.IsOwner = false
			}
		}
		collectionItems = append(collectionItems, temp)
	}

	protoTags, err := s.tagsClient.GetTagsByEntityId(context.Background(), &microservices.GetTagsByEntityIdRequest{EntityId: dbCollection.ID.String()})
	if err != nil {
		return nil, err
	}

	var resultTags []models.ShortTags = utils.NormalizeTagsShort(protoTags.GetTags())

	finalCollection.Tags = resultTags

	finalCollection.CollectionItems = collectionItems
	finalCollection.ShareString = utils.DefineShareString(authorizerUser, userLogin, dbCollection)
	finalCollection.CollectionItemsCount, _ = s.collectionItemRepo.GetCountCI(dbCollection.ID)
	finalCollection.TotalPrice, _ = s.collectionItemRepo.Sum(dbCollection.ID)
	finalCollection.ShippingTotal, _ = s.collectionItemRepo.SumShippingCost(dbCollection.ID)
	finalCollection.LikesCount = int64(len(dbCollection.Likes))
	finalCollection.CanLike = utils.CanLike(dbCollection.Likes, authorizerUser, dbCollection.UserLogin)
	finalCollection.IsOwner = authorizerUser == userLogin
	finalCollection.PhotosCount = photosCount.GetCount()
	if authUser != nil {
		subArray = authUser.Subscriptions
		subsExtended, _ := s.userRepo.GetForSubs(subArray)
		subscribersArray := authUser.SubscribersLogins
		subscribersExtended, _ := s.userRepo.GetForSubs(subscribersArray)
		finalCollection.User.SubscriptionsExtended = subsExtended
		finalCollection.User.SubscribersExtended = subscribersExtended

		var lowerCasedUsers []string
		for _, userInArray := range subArray {
			lowerCasedUsers = append(lowerCasedUsers, strings.ToLower(userInArray))
		}
		canSubscribe := !slices.Contains(lowerCasedUsers, strings.ToLower(userLogin))
		finalCollection.User.CanSubscribe = canSubscribe

	}
	return &models.CollectionDataResponse{Data: finalCollection}, nil

}

func (s *CollectionService) Like(id string, userLogin string) (*dto.CommonResponse, error) {
	exists, err := s.repo.GetByIdWithoutCollectionItems(id)
	if err != nil {
		return nil, err
	}
	isUserLikes := slices.Contains(exists.Likes, userLogin)
	if !isUserLikes {
		exists.Likes = append(exists.Likes, userLogin)
		if !exists.IsPrivate {
			go func() {
				eventError := s.eventsService.AddEvent(userLogin, utils.EventActionLike, utils.EventTargetCollection, exists.Name, &models.EventsParams{TargetCollectionID: exists.ID})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
		}

		err = s.userRepo.IncrementExperience(exists.UserLogin, utils.CollectionSelfLikeExp)
		if err != nil {
			return nil, err
		}
		target := &dto.TargetItem{
			ID:              id,
			Name:            exists.Name,
			Transliteration: exists.Transliteration,
			TargetType:      "collection",
		}
		go func() {
			err = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
				Login:       exists.UserLogin,
				TargetID:    id,
				SenderLogin: userLogin,
				Type:        utils.NotificationTypeCollection,
				Action:      utils.NotificationActionLike,
				Date:        time.Now().Format(time.RFC3339),
				OwnerLogin:  exists.UserLogin,
			}, target)
		}()
	} else {
		exists.Likes = utils.RemoveByValue(exists.Likes, userLogin)
		if !exists.IsPrivate {
			go func() {
				eventError := s.eventsService.AddEvent(userLogin, utils.EventActionDislike, utils.EventTargetCollection, exists.Name, &models.EventsParams{TargetCollectionID: exists.ID})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
		}
		err = s.userRepo.DecrementExperience(exists.UserLogin, utils.CollectionSelfLikeExp)
		if err != nil {
			return nil, err
		}
		go func() {
			_ = s.notificationsService.DeleteNotification(context.Background(), id, userLogin)
		}()
	}

	_, err = s.repo.UpdateCollection(exists, exists)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

// FIXME отключено
//func (s *CollectionService) Subscribe(targetId string, userLogin string, isSubscribe bool) (*dto.CommonResponse, error) {
//	var subscriber *models.Users
//	var dbCollection *models.Collections
//	var err1, err2 error
//	var wg sync.WaitGroup
//
//	wg.Add(2)
//	go func() {
//		defer wg.Done()
//		subscriber, err1 = s.userRepo.GetUserByLogin(userLogin)
//	}()
//	go func() {
//		defer wg.Done()
//		dbCollection, err2 = s.repo.GetCollectionByIdWithoutLimits(targetId)
//	}()
//	wg.Wait()
//	if err1 != nil {
//		return nil, err1
//	}
//	if err2 != nil {
//		return nil, err2
//	}
//	if isSubscribe {
//		subscriber.CollectionSubscriptions = append(subscriber.CollectionSubscriptions, targetId)
//		dbCollection.SubscribersCount = dbCollection.SubscribersCount + 1
//		go func() {
//			eventError := s.eventsService.AddEvent(userLogin, utils.EventActionSubscribe, utils.EventTargetCollection, dbCollection.Name, &models.EventsParams{TargetCollectionID: dbCollection.ID})
//			if eventError != nil {
//				log.Default().Print(eventError)
//			}
//		}()
//		err := s.userRepo.IncrementExperience(dbCollection.UserLogin, utils.CollectionSelfSubExp)
//		if err != nil {
//			return nil, err
//		}
//		target := &dto.TargetItem{
//			ID:              targetId,
//			Name:            dbCollection.Name,
//			Transliteration: dbCollection.Transliteration,
//			TargetType:      "collection",
//		}
//		go func() {
//			_ = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
//				Login:       dbCollection.UserLogin,
//				TargetID:    targetId,
//				SenderLogin: userLogin,
//				Type:        utils.NotificationTypeCollection,
//				Action:      utils.NotificationActionSubscribe,
//				Date:        time.Now().Format(time.RFC3339),
//				OwnerLogin:  dbCollection.UserLogin,
//			}, target)
//		}()
//	} else {
//		subscriber.CollectionSubscriptions = utils.RemoveByValue(subscriber.CollectionSubscriptions, targetId)
//		dbCollection.SubscribersCount = dbCollection.SubscribersCount - 1
//		err := s.userRepo.DecrementExperience(dbCollection.UserLogin, utils.CollectionSelfSubExp)
//		if err != nil {
//			return nil, err
//		}
//		go func() {
//			_ = s.notificationsService.DeleteNotification(context.Background(), targetId, userLogin)
//		}()
//	}
//	wg.Add(2)
//	go func() {
//		defer wg.Done()
//		_, err1 = s.userRepo.UpdateUser(subscriber, *subscriber)
//	}()
//	go func() {
//		defer wg.Done()
//		_, err2 = s.repo.UpdateCollection(dbCollection, dbCollection)
//	}()
//	wg.Wait()
//	if err1 != nil {
//		return nil, err1
//	}
//	if err2 != nil {
//		return nil, err2
//	}
//	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
//
//}
