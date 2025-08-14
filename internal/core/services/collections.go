package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
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
	tagsRepo             *repositories.TagsRepository
	userRepo             *repositories.UsersRepository
	eventRepo            *repositories.EventsRepository
	collectionItemRepo   *repositories.CiRepository
	s3Service            *s3.S3Service
	notificationsService *NotificationsService
}

func NewCollectionService(repo *repositories.CollectionsRepository, tagsRepo *repositories.TagsRepository, userRepo *repositories.UsersRepository, eventRepo *repositories.EventsRepository, collectionItemRepo *repositories.CiRepository, seService *s3.S3Service, notificationsService *NotificationsService) *CollectionService {
	return &CollectionService{repo: repo, tagsRepo: tagsRepo, userRepo: userRepo, eventRepo: eventRepo, collectionItemRepo: collectionItemRepo, s3Service: seService, notificationsService: notificationsService}
}

func (s *CollectionService) processTags(tags []string) ([]models.Tags, error) {
	var resultTags []models.Tags
	for _, tag := range tags {
		tagByName, err := s.tagsRepo.GetTagByName(tag)
		if err != nil {
			fmt.Errorf("failed to get tag: %w", err)
		}

		if tagByName == nil {
			newTag := &models.Tags{Name: tag}
			tagByName, err = s.tagsRepo.AddTag(newTag)
			if err != nil {
				fmt.Errorf("failed to add tag: %w", err)
			}
			if tagByName == nil {
				fmt.Errorf("unexpected nil tag after adding")
			}
		}

		resultTags = append(resultTags, *tagByName)
	}
	return resultTags, nil
}
func (s *CollectionService) Create(dto *models.CollectionCreateRequest) (*models.CollectionDataResponse, error) {
	existsCollection, err := s.repo.GetUniqueByName(dto.Name, dto.UserLogin)
	if err != nil {
		fmt.Println(err)
	}
	if existsCollection != nil {
		return nil, errors.New("collection с таким именем уже существует")
	}

	processTags, _ := s.processTags(dto.Tags)
	dtoUser, _ := s.userRepo.GetUserByLogin(dto.UserLogin)
	dbCollection := &models.Collections{
		Name:            dto.Name,
		Description:     dto.Description,
		BannerUrl:       dto.BannerUrl,
		IsPrivate:       dto.IsPrivate,
		Transliteration: dto.Transliteration,
		Tags:            processTags,
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
		eventError := s.eventRepo.AddEvent(dto.UserLogin, utils.EventActionCreate, utils.EventTargetCollection, dto.Name, nil, &res.Id, nil, nil)
		if eventError != nil {
			log.Default().Print(eventError)
		}

	}

	var response models.CollectionsResponse

	e := mapstructure.Decode(res, &response)
	response.Id = res.Id

	if e != nil {
		return nil, e
	}

	er := mapstructure.Decode(dtoUser, &response.User)

	if er != nil {
		return nil, er
	}
	response.CanLike = true
	response.CanSubscribe = true
	response.IsOwner = true
	return &models.CollectionDataResponse{Data: response}, nil
}

func (s *CollectionService) Update(id string, dto *models.CollectionUpdateRequest) (*models.CollectionDataResponse, error) {
	exists, err := s.repo.GetCollectionByIdWithoutLimits(id)
	if err != nil {
		return nil, err
	}

	processTags, _ := s.processTags(dto.Tags)

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

	exists.Tags = processTags
	exists.IsPrivate = *dto.IsPrivate

	resultCollection, err := s.repo.UpdateCollectionFull(exists)
	if err != nil {
		return nil, err
	}

	if !*dto.IsPrivate {
		eventError := s.eventRepo.AddEvent(
			*dto.UserLogin,
			utils.EventActionUpdate,
			utils.EventTargetCollection,
			*dto.Name,
			nil,
			&resultCollection.Id,
			nil,
			nil,
		)
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}

	var finalCollection models.CollectionsResponse
	if err := mapstructure.Decode(resultCollection, &finalCollection); err != nil {
		return nil, err
	}
	finalCollection.CanLike = true
	finalCollection.CanSubscribe = true
	finalCollection.IsOwner = true

	return &models.CollectionDataResponse{Data: finalCollection}, nil
}

func (s *CollectionService) Delete(id string, ctx context.Context) (*dto.CommonResponse, error) {
	exists, _ := s.repo.GetCollectionByIdWithoutLimits(id)
	if !exists.IsPrivate {
		go func() {
			eventError := s.eventRepo.AddEvent(exists.User.Login, utils.EventActionDelete, utils.EventTargetCollection, exists.Name, nil, &exists.Id, nil, nil)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}()
		if exists.BannerUrl != "" {
			_, err := s.s3Service.DeleteFiles(ctx, s3.DeleteFilesRequest{Keys: []string{exists.BannerUrl}})
			if err != nil {
				return nil, err
			}
		}
	}

	exp := len(exists.CollectionItems) + utils.CollectionExp

	cis := exists.CollectionItems
	for _, ci := range cis {
		if len(ci.Entities) != 0 {
			exp += utils.EntityAttachExp
		}
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

	var authUser *models.Users
	var subArray []string
	if authorizedUser != "" {
		authUser, _ = s.userRepo.GetUserByLogin(authorizedUser)
		subArray = authUser.CollectionSubscriptions
	}
	result := make([]models.CollectionsResponse, 0)

	collectionIds := make([]uuid.UUID, 0)
	for _, dbCollection := range collections {
		collectionIds = append(collectionIds, dbCollection.Id)
	}

	counts, _ := s.collectionItemRepo.GetCountCIByIDs(collectionIds)
	totalPrices, _ := s.collectionItemRepo.GetSumsByCollectionIDs(collectionIds)
	shippingCosts, _ := s.collectionItemRepo.GetShippingCostsByCollectionIDs(collectionIds)

	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		err := mapstructure.Decode(dbCollection, &temp)
		temp.ShareString = utils.DefineShareString(authorizedUser, login, &dbCollection)
		temp.CollectionItemsCount = counts[dbCollection.Id]
		temp.TotalPrice = totalPrices[dbCollection.Id]
		temp.ShippingTotal = shippingCosts[dbCollection.Id]
		temp.CanSubscribe = !slices.Contains(subArray, dbCollection.Id.String())
		temp.LikesCount = int64(len(dbCollection.Likes))
		temp.CanLike = true
		temp.IsOwner = authorizedUser == login
		if authorizedUser != "" {
			temp.CanLike = !slices.Contains(dbCollection.Likes, authorizedUser)
		}

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
	var authUser *models.Users
	var subArray []string
	if authorizedUser != "" {
		authUser, _ = s.userRepo.GetUserByLogin(authorizedUser)
		subArray = authUser.CollectionSubscriptions
	}
	result := make([]models.CollectionsResponse, 0)

	collectionIds := make([]uuid.UUID, 0)
	for _, dbCollection := range collections {
		collectionIds = append(collectionIds, dbCollection.Id)
	}

	counts, _ := s.collectionItemRepo.GetCountCIByIDs(collectionIds)
	totalPrices, _ := s.collectionItemRepo.GetSumsByCollectionIDs(collectionIds)
	shippingCosts, _ := s.collectionItemRepo.GetShippingCostsByCollectionIDs(collectionIds)

	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		err := mapstructure.Decode(dbCollection, &temp)
		temp.CollectionItemsCount = counts[dbCollection.Id]
		temp.TotalPrice = totalPrices[dbCollection.Id]
		temp.ShippingTotal = shippingCosts[dbCollection.Id]
		temp.CanSubscribe = !slices.Contains(subArray, dbCollection.Id.String())
		temp.LikesCount = int64(len(dbCollection.Likes))
		temp.CanLike = true
		temp.IsOwner = authorizedUser == temp.User.Login
		if authorizedUser != "" {
			temp.CanLike = !slices.Contains(dbCollection.Likes, authorizedUser)
		}

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
		subArray = authUser.CollectionSubscriptions
	}
	var finalCollection models.CollectionsResponse
	err := mapstructure.Decode(dbCollection, &finalCollection)
	if err != nil {
		return nil, err
	}

	collectionItems := make([]models.CollectionItemsResponse, 0)

	for _, item := range dbCollection.CollectionItems {
		var temp models.CollectionItemsResponse
		itemErr := mapstructure.Decode(item, &temp)
		if itemErr != nil {
			return nil, itemErr
		}
		temp.Collection = finalCollection.Id
		temp.LikesCount = int64(len(item.Likes))
		temp.CanLike = true
		temp.CollectionTransliteration = finalCollection.Transliteration
		if authorizerUser != "" {
			if authorizerUser == item.Owner.Login {
				temp.CanLike = true
				temp.IsOwner = true
			} else {
				temp.CanLike = !slices.Contains(item.Likes, authorizerUser)
				temp.IsOwner = false
			}
		}
		collectionItems = append(collectionItems, temp)
	}

	finalCollection.CollectionItems = collectionItems
	finalCollection.ShareString = utils.DefineShareString(authorizerUser, userLogin, dbCollection)
	finalCollection.CollectionItemsCount, _ = s.collectionItemRepo.GetCountCI(dbCollection.Id)
	finalCollection.TotalPrice, _ = s.collectionItemRepo.Sum(dbCollection.Id)
	finalCollection.ShippingTotal, _ = s.collectionItemRepo.SumShippingCost(dbCollection.Id)
	finalCollection.CanSubscribe = !slices.Contains(subArray, dbCollection.Id.String())
	finalCollection.LikesCount = int64(len(dbCollection.Likes))
	finalCollection.CanLike = true
	finalCollection.IsOwner = authorizerUser == userLogin
	if authUser != nil {
		subArray = authUser.Subscriptions
		var lowerCasedUsers []string
		for _, userInArray := range subArray {
			lowerCasedUsers = append(lowerCasedUsers, strings.ToLower(userInArray))
		}
		canSubscribe := !slices.Contains(lowerCasedUsers, strings.ToLower(userLogin))
		finalCollection.User.CanSubscribe = canSubscribe
		finalCollection.CanLike = !slices.Contains(dbCollection.Likes, authUser.Login)

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
			eventError := s.eventRepo.AddEvent(userLogin, utils.EventActionLike, utils.EventTargetCollection, exists.Name, nil, &exists.Id, nil, nil)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}

		err = s.userRepo.IncrementExperience(exists.UserLogin, utils.CollectionSelfLikeExp)
		if err != nil {
			return nil, err
		}
		go func() {
			err = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
				Login:       exists.UserLogin,
				TargetId:    id,
				SenderLogin: userLogin,
				Type:        utils.NotificationTypeCollection,
				Action:      utils.NotificationActionLike,
				Date:        time.Now().Format(time.RFC3339),
			})
		}()
	} else {
		exists.Likes = utils.RemoveByValue(exists.Likes, userLogin)
		if !exists.IsPrivate {
			go func() {
				eventError := s.eventRepo.AddEvent(userLogin, utils.EventActionLike, utils.EventTargetCollection, exists.Name, nil, &exists.Id, nil, nil)
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
		}
		err = s.userRepo.DecrementExperience(exists.UserLogin, utils.CollectionSelfLikeExp)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.repo.UpdateCollection(exists, exists)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CollectionService) Subscribe(targetId string, userLogin string, isSubscribe bool) (*dto.CommonResponse, error) {
	subscriber, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, err
	}
	dbCollection, err := s.repo.GetCollectionByIdWithoutLimits(targetId)
	if err != nil {
		return nil, err
	}
	if isSubscribe {
		subscriber.CollectionSubscriptions = append(subscriber.CollectionSubscriptions, targetId)
		dbCollection.SubscribersCount = dbCollection.SubscribersCount + 1
		eventError := s.eventRepo.AddEvent(userLogin, utils.EventActionSubscribe, utils.EventTargetCollection, dbCollection.Name, nil, &dbCollection.Id, nil, nil)
		if eventError != nil {
			log.Default().Print(eventError)
		}
		err = s.userRepo.IncrementExperience(dbCollection.UserLogin, utils.CollectionSelfSubExp)
		if err != nil {
			return nil, err
		}
		go func() {
			err = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
				Login:       dbCollection.UserLogin,
				TargetId:    targetId,
				SenderLogin: userLogin,
				Type:        utils.NotificationTypeCollection,
				Action:      utils.NotificationActionSubscribe,
				Date:        time.Now().Format(time.RFC3339),
			})
		}()
	} else {
		subscriber.CollectionSubscriptions = utils.RemoveByValue(subscriber.CollectionSubscriptions, targetId)
		dbCollection.SubscribersCount = dbCollection.SubscribersCount - 1
		err = s.userRepo.DecrementExperience(dbCollection.UserLogin, utils.CollectionSelfSubExp)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.userRepo.UpdateUser(subscriber, *subscriber)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.UpdateCollection(dbCollection, dbCollection)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil

}

func (s *CollectionService) GetByTag(tag string, authUserLogin string, limit string, offset string, orderBy string, order string, search string) (*models.AllCollectionsDataByTag, error) {
	collections, totalCount, err := s.repo.GetCollectionByTag(tag, limit, offset, search)
	if err != nil {
		return nil, err
	}
	dbTag, err := s.tagsRepo.GetTagByName(tag)
	if err != nil {
		return nil, err
	}
	var authUser *models.Users
	var subArray []string
	if authUserLogin != "" {
		authUser, _ = s.userRepo.GetUserByLogin(authUserLogin)
		subArray = authUser.CollectionSubscriptions
	}
	result := make([]models.CollectionsResponse, 0)

	collectionIds := make([]uuid.UUID, 0)
	for _, dbCollection := range collections {
		collectionIds = append(collectionIds, dbCollection.Id)
	}

	counts, _ := s.collectionItemRepo.GetCountCIByIDs(collectionIds)
	totalPrices, _ := s.collectionItemRepo.GetSumsByCollectionIDs(collectionIds)
	shippingCosts, _ := s.collectionItemRepo.GetShippingCostsByCollectionIDs(collectionIds)

	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		errMap := mapstructure.Decode(dbCollection, &temp)
		temp.ShareString = utils.DefineShareString(authUserLogin, dbCollection.User.Login, &dbCollection)
		temp.CollectionItemsCount = counts[dbCollection.Id]
		temp.TotalPrice = totalPrices[dbCollection.Id]
		temp.ShippingTotal = shippingCosts[dbCollection.Id]
		temp.CanSubscribe = !slices.Contains(subArray, dbCollection.Id.String())
		temp.IsOwner = authUserLogin == dbCollection.User.Login
		temp.CanLike = !slices.Contains(dbCollection.Likes, authUserLogin)
		if errMap != nil {
			return nil, errMap
		}
		result = append(result, temp)
	}

	sortedCollections := utils.GetCollectionOrderBy(orderBy, result, order)

	return &models.AllCollectionsDataByTag{Data: sortedCollections, Total: totalCount, Tag: *dbTag}, nil
}
