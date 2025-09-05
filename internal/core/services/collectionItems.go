package services

import (
	"context"
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
	"time"
)

type CiService struct {
	repo                 *repositories.CiRepository
	eventRepo            *repositories.EventsRepository
	collectionRepo       *repositories.CollectionsRepository
	userRepo             *repositories.UsersRepository
	platformsRepo        *repositories.PlatformsRepository
	entityRepo           *repositories.EntitiesRepository
	s3Service            *s3.S3Service
	itemTypeRepo         *repositories.ItemTypesRepository
	notificationsService *NotificationsService
}

func NewCiService(repo *repositories.CiRepository, eventRepo *repositories.EventsRepository, collectionRepo *repositories.CollectionsRepository, userRepo *repositories.UsersRepository, platformsRepo *repositories.PlatformsRepository, entityRepo *repositories.EntitiesRepository, s3Service *s3.S3Service, itemTypeRepo *repositories.ItemTypesRepository, notificationsService *NotificationsService) *CiService {
	return &CiService{repo: repo, eventRepo: eventRepo, collectionRepo: collectionRepo, userRepo: userRepo, platformsRepo: platformsRepo, entityRepo: entityRepo, s3Service: s3Service, itemTypeRepo: itemTypeRepo, notificationsService: notificationsService}
}

func (s *CiService) getEntities(entities []string, userLogin string) []models.Entities {
	var resultEntities []models.Entities
	for _, entity := range entities {
		entityByTranslit, err := s.entityRepo.GetEntityByName(entity)
		if err != nil {
			fmt.Errorf("failed to get entity: %w", err)
		}
		if entityByTranslit == nil {
			newEntity := &models.Entities{Name: entity, Transliteration: utils.Slugify(entity), Author: userLogin}
			entityByTranslit, err = s.entityRepo.CreateEntity(newEntity)
			if err != nil {
				fmt.Errorf("failed to add entity: %w", err)
			}
			existsUser, _ := s.userRepo.GetUserByLogin(userLogin)
			err = s.userRepo.IncrementExperience(existsUser.Login, utils.EntityExp)
			err = s.userRepo.IncrementSocialScore(userLogin, 1)
			if err != nil {
				fmt.Errorf("can't add rating: %w", err)
			}

			if entityByTranslit == nil {
				fmt.Errorf("unexpected nil entity after adding")
			}
		}
		resultEntities = append(resultEntities, *entityByTranslit)
	}
	return resultEntities
}

func (s *CiService) Create(dto *models.CollectionItemsRequestCreate, authUserLogin string) (*models.CollectionItemsDataResponse, error) {
	var platform *models.Platforms
	var platformID *uuid.UUID
	var itemType *models.ItemTypes
	var itemTypeID *uuid.UUID

	if dto.Platform != "" {
		foundPlatform, err := s.platformsRepo.GetPlatformById(dto.Platform)
		if err != nil {
			return nil, fmt.Errorf("error getting platform: %v", err)
		}
		platform = foundPlatform
		platformID = &foundPlatform.Id
	} else {
		platform = nil
		platformID = nil
	}

	if dto.ItemType != "" {
		foundItemType, err := s.itemTypeRepo.GetTypeById(dto.ItemType)
		if err != nil {
			return nil, fmt.Errorf("error getting item type: %v", err)
		}
		itemType = foundItemType
		itemTypeID = &foundItemType.Id
	} else {
		itemType = nil
		itemTypeID = nil
	}

	var entities []models.Entities

	if len(dto.Entities) == 0 {
		entities = make([]models.Entities, 0)
	} else {
		entities = s.getEntities(dto.Entities, authUserLogin)
	}
	owner, ownerErr := s.userRepo.GetUserByLogin(authUserLogin)
	if ownerErr != nil {
		return nil, ownerErr
	}

	collection, _ := s.collectionRepo.GetByIdWithoutCollectionItems(dto.Collection)
	collectionsSlice := make([]models.Collections, 0)
	collectionsSlice = append(collectionsSlice, *collection)
	dbCollectionItem := &models.CollectionItems{
		Name:          dto.Name,
		Description:   dto.Description,
		Images:        dto.Images,
		PurchaseDate:  dto.PurchaseDate,
		PurchasePrice: dto.PurchasePrice,
		Sealed:        dto.Sealed,
		Edition:       dto.Edition,
		ShippingCost:  dto.ShippingCost,
		Entities:      entities,
		Platform:      platform,
		ItemType:      itemType,
		Owner:         *owner,
		Collections:   collectionsSlice,
		PlatformID:    platformID,
		UserLogin:     owner.Login,
		ItemTypeID:    itemTypeID,
		Rating:        dto.Rating,
	}
	if dto.CopyNumber != nil {
		dbCollectionItem.CopyNumber = dto.CopyNumber
	}
	collectionItem, err := s.repo.CreateCI(dbCollectionItem)
	if err != nil {
		return nil, err
	}
	if !collection.IsPrivate {
		eventError := s.eventRepo.AddEvent(authUserLogin, utils.EventActionCreate, utils.EventTargetCollectionItem, dto.Name, nil, nil, &collectionItem.Id, nil)
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(collectionItem, &collectionItemResponse)
	collectionItemResponse.Collection = collection.Id
	if err != nil {
		return nil, err
	}
	collectionItemResponse.LikesCount = int64(len(collectionItem.Likes))
	collectionItemResponse.CanLike = true
	collectionItemResponse.IsOwner = true

	exp := utils.CIExp

	if len(dto.Entities) != 0 {
		exp += utils.EntityAttachExp
	}
	if dto.PurchaseDate != "" {
		exp += utils.CIPurchaseDateExp
	}
	if dto.PurchasePrice != 0 {
		exp += utils.CIPurchasePriceExp
	}
	if dto.Rating != 0 {
		exp += utils.CIRatingExp
	}
	if len(dto.Images) != 0 {
		exp += utils.PictureExp
	}
	if len(dto.CopyNumber) != 0 {
		exp += utils.CICopyNumberExp
	}

	err = s.userRepo.IncrementExperience(authUserLogin, exp)
	if err != nil {
		return nil, err
	}

	return &models.CollectionItemsDataResponse{Data: collectionItemResponse}, nil
}

func (s *CiService) Delete(id string, ctx context.Context) (*dto.CommonResponse, error) {
	exists, _ := s.repo.GetCIByID(id)
	var result *dto.CommonResponse
	if exists != nil {
		if !exists.Collections[0].IsPrivate {
			eventError := s.eventRepo.AddEvent(exists.Owner.Login, utils.EventActionDelete, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id, nil)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}
		err := s.repo.DeleteCI(id)
		if err != nil {
			return nil, err
		}
		if len(exists.Images) > 0 {
			_, err = s.s3Service.DeleteFiles(ctx, s3.DeleteFilesRequest{Keys: exists.Images})
			if err != nil {
				return nil, err
			}
		}
		result = &dto.CommonResponse{Data: dto.Resp{Success: true}}

		exp := utils.CIExp

		if len(exists.Entities) != 0 {
			exp += utils.EntityAttachExp
		}
		if exists.PurchaseDate != "" {
			exp += utils.CIPurchaseDateExp
		}
		if exists.PurchasePrice != 0 {
			exp += utils.CIPurchasePriceExp
		}
		if exists.Rating != 0 {
			exp += utils.CIRatingExp
		}
		if len(exists.Images) != 0 {
			exp += utils.PictureExp
		}
		if len(exists.CopyNumber) != 0 {
			exp += utils.CICopyNumberExp
		}

		err = s.userRepo.DecrementExperience(exists.UserLogin, exp)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *CiService) Update(id string, dto *models.CollectionItemsRequestUpdate) (*models.CollectionItems, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	var resultEntities []models.Entities

	if len(dto.Entities) != 0 {
		resultEntities = s.getEntities(dto.Entities, exists.Owner.Login)
	} else {
		resultEntities = make([]models.Entities, 0)
	}

	err = s.updateCIExp(exists, dto)
	if err != nil {
		return nil, err
	}

	dst := reflect.ValueOf(exists).Elem()
	src := reflect.ValueOf(dto).Elem()

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		fieldName := src.Type().Field(i).Name

		if fieldName == "Entities" || fieldName == "Platform" || fieldName == "ItemType" {
			continue
		}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		dstField := dst.FieldByName(fieldName)
		if dstField.IsValid() && dstField.CanSet() {
			if fieldName == "PlatformID" || fieldName == "ItemTypeID" {
				if str, ok := field.Interface().(string); ok && str != "" {
					uuidVal, err := uuid.Parse(str)
					if err != nil {
						return nil, fmt.Errorf("invalid UUID format for field %s: %v", fieldName, err)
					}
					dstField.Set(reflect.ValueOf(uuidVal))
				}
			} else {
				dstField.Set(field)
			}
		}
	}

	//var dbCollectionItem models.CollectionItems
	//err = mapstructure.Decode(dto, &dbCollectionItem)
	//dbCollectionItem.Entities = resultEntities

	exists.Entities = resultEntities

	var platform *models.Platforms
	var platformID *uuid.UUID
	var itemType *models.ItemTypes
	var itemTypeID *uuid.UUID

	if dto.Platform != nil && *dto.Platform != "" {
		foundPlatform, err := s.platformsRepo.GetPlatformById(*dto.Platform)
		if err != nil {
			return nil, fmt.Errorf("error getting platform: %v", err)
		}
		platform = foundPlatform
		platformID = &foundPlatform.Id
	} else {
		platform = nil
		platformID = nil
	}

	if dto.ItemType != nil && *dto.ItemType != "" {
		foundItemType, err := s.itemTypeRepo.GetTypeById(*dto.ItemType)
		if err != nil {
			return nil, fmt.Errorf("error getting item type: %v", err)
		}
		itemType = foundItemType
		itemTypeID = &foundItemType.Id
	} else {
		itemType = nil
		itemTypeID = nil
	}

	exists.Platform = platform
	exists.PlatformID = platformID
	exists.ItemTypeID = itemTypeID
	exists.ItemType = itemType
	result, err := s.repo.UpdateCIFull(exists)
	if err != nil {
		return nil, err
	}

	if !exists.Collections[0].IsPrivate {
		eventError := s.eventRepo.AddEvent(exists.Owner.Login, utils.EventActionUpdate, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id, nil)
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}
	return result, nil
}

func (s *CiService) GetById(id string, authUser string) (*models.CollectionItems, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	return exists, nil
}

func (s *CiService) CopyOrMove(id string, targetCollectionIds []string, sourceCollectionId string) (*dto.CommonResponse, error) {
	var result *dto.CommonResponse
	if sourceCollectionId == "" {
		collections := make([]models.Collections, 0)
		for _, collectionId := range targetCollectionIds {
			collection, _ := s.collectionRepo.GetCollectionByIdWithoutLimits(collectionId)
			collections = append(collections, *collection)
		}
		collectionItem, err := s.repo.GetCIByID(id)
		if err != nil {
			return nil, err
		}

		for _, collection := range collections {
			collectionItem.Collections = append(collectionItem.Collections, collection)
			_, erroring := s.repo.UpdateCI(collectionItem, collectionItem)
			if erroring != nil {
				return nil, erroring
			}
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}

	sourceCollection, err := s.collectionRepo.GetCollectionByIdWithoutLimits(sourceCollectionId)
	if err != nil {
		return nil, err
	}
	targetCollection, err := s.collectionRepo.GetCollectionByIdWithoutLimits(targetCollectionIds[0])
	if err != nil {
		return nil, err
	}
	collectionItem, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}
	existsSourceCollection := sourceCollection
	existsTargetCollection := targetCollection
	sourceCollectionItems := utils.RemoveByValueStruct(sourceCollection.CollectionItems, collectionItem.Id)
	targetCollectionItems := append(targetCollection.CollectionItems, *collectionItem)
	sourceCollection.CollectionItems = sourceCollectionItems
	targetCollection.CollectionItems = targetCollectionItems
	_, errCol := s.collectionRepo.UpdateCollection(existsSourceCollection, sourceCollection)
	if errCol != nil {
		return nil, err
	}
	_, errCol = s.collectionRepo.UpdateCollection(existsTargetCollection, targetCollection)
	if errCol != nil {
		return nil, err
	}
	ok, err := s.collectionRepo.DeleteRelation(sourceCollectionId, id)
	if !ok || err != nil {
		return nil, err
	}
	result = &dto.CommonResponse{Data: dto.Resp{Success: true}}

	return result, nil
}

func (s *CiService) Like(id string, userLogin string) (*dto.CommonResponse, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}
	isUserLikes := slices.Contains(exists.Likes, userLogin)
	if !isUserLikes {
		exists.Likes = append(exists.Likes, userLogin)
		if !exists.Collections[0].IsPrivate {
			go func() {
				eventError := s.eventRepo.AddEvent(userLogin, utils.EventActionLike, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id, nil)
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
		}
		err = s.userRepo.IncrementExperience(exists.UserLogin, utils.CISelfLikeExp)
		if err != nil {
			return nil, err
		}
		target := &dto.TargetItem{
			Id:              id,
			Name:            exists.Name,
			Transliteration: exists.Transliteration,
			TargetType:      "collectionItem",
		}
		go func() {
			err = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
				Login:       exists.UserLogin,
				TargetId:    id,
				SenderLogin: userLogin,
				Type:        utils.NotificationTypeCollectionItem,
				Action:      utils.NotificationActionLike,
				Date:        time.Now().Format(time.RFC3339),
				OwnerLogin:  exists.UserLogin,
			}, target)
		}()
	} else {
		exists.Likes = utils.RemoveByValue(exists.Likes, userLogin)
		err = s.userRepo.DecrementExperience(exists.UserLogin, utils.CISelfLikeExp)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.repo.UpdateCI(exists, exists)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CiService) GetByEntity(entity string, authUserLogin string, limit string) (*models.CollectionItemsDataSortedResponse, error) {
	dbEntity, err := s.entityRepo.GetEntityByTranslit(entity)
	if dbEntity == nil {
		return nil, fmt.Errorf("entity not found")
	}
	collectionItems, totalCount, err := s.repo.GetCollectionItemsByEntity(entity, limit)
	if err != nil {
		return nil, err
	}

	var sortedCollectionItems models.CollectionItemsSortedResponse

	sortedCollectionItems.CollectibleFigures, _ = processCI(collectionItems.CollectibleFigures, authUserLogin)
	sortedCollectionItems.Books, _ = processCI(collectionItems.Books, authUserLogin)
	sortedCollectionItems.BoardGames, _ = processCI(collectionItems.BoardGames, authUserLogin)
	sortedCollectionItems.Comics, _ = processCI(collectionItems.Comics, authUserLogin)
	sortedCollectionItems.GamingHardware, _ = processCI(collectionItems.GamingHardware, authUserLogin)
	sortedCollectionItems.Vinyl, _ = processCI(collectionItems.Vinyl, authUserLogin)
	sortedCollectionItems.VideoGames, _ = processCI(collectionItems.VideoGames, authUserLogin)
	sortedCollectionItems.Steelbooks, _ = processCI(collectionItems.Steelbooks, authUserLogin)
	sortedCollectionItems.CollectibleCards, _ = processCI(collectionItems.CollectibleCards, authUserLogin)

	return &models.CollectionItemsDataSortedResponse{Data: sortedCollectionItems, Total: totalCount, Entity: *dbEntity}, nil
}

func (s *CiService) GetByEntityAndType(entity string, itemType string, limit string, offset string, search string, orderBy string, order string, authUser string) ([]models.CollectionItems, *models.Entities, int64, error) {
	dbEntity, err := s.entityRepo.GetEntityByTranslit(entity)
	if dbEntity == nil {
		return nil, nil, 0, fmt.Errorf("entity not found")
	}
	collectionItems, totalCount, err := s.repo.GetByEntityAndType(entity, itemType, limit, offset, search, orderBy, order)
	if err != nil {
		return nil, nil, totalCount, err
	}

	return collectionItems, dbEntity, totalCount, err
}

func (s *CiService) GetAll(limit, offset, search string) (*models.CollectionItemsDataPoor, error) {
	collectionItems, err := s.repo.GetAll(limit, offset, search)
	if err != nil {
		return nil, err
	}

	var collectionItemsAll []models.CollectionItemsResponse

	for _, item := range collectionItems {
		var temp models.CollectionItemsResponse
		err := mapstructure.Decode(item, &temp)
		if err != nil {
			return nil, err
		}

		if len(item.Collections) > 0 {
			temp.Collection = item.Collections[0].Id
			temp.CollectionTransliteration = item.Collections[0].Transliteration
		} else {
			temp.Collection = uuid.Nil
			temp.CollectionTransliteration = ""
		}

		if item.Likes != nil {
			temp.LikesCount = int64(len(item.Likes))
		} else {
			temp.LikesCount = 0
		}

		temp.Owner = item.Owner
		collectionItemsAll = append(collectionItemsAll, temp)
	}

	return &models.CollectionItemsDataPoor{Data: collectionItemsAll}, nil
}

func processCI(slice []models.CollectionItems, authUser string) ([]models.CollectionItemsResponse, error) {
	var resultCollectionItems []models.CollectionItemsResponse
	for _, item := range slice {
		var temp models.CollectionItemsResponse
		err := mapstructure.Decode(item, &temp)
		if err != nil {
			return nil, err
		}
		temp.Collection = item.Collections[0].Id
		temp.Owner = item.Owner
		temp.LikesCount = int64(len(item.Likes))
		temp.CanLike = true
		temp.IsOwner = authUser == item.Owner.Login
		temp.CollectionTransliteration = item.Collections[0].Transliteration
		if authUser != "" {
			if authUser == item.Owner.Login {
				temp.CanLike = true
			} else {
				temp.CanLike = !slices.Contains(item.Likes, authUser)
			}

		}
		resultCollectionItems = append(resultCollectionItems, temp)
	}
	return resultCollectionItems, nil
}

func (s *CiService) updateCIExp(exists *models.CollectionItems, dto *models.CollectionItemsRequestUpdate) error {
	var exp int
	if dto.Entities != nil {
		if len(exists.Entities) != 0 && len(dto.Entities) == 0 {
			exp -= utils.EntityAttachExp
		} else if len(exists.Entities) == 0 && len(dto.Entities) != 0 {
			exp += utils.EntityAttachExp
		}
	}

	if exists.PurchaseDate != "" && dto.PurchaseDate == nil {
		exp -= utils.CIPurchaseDateExp
	} else if exists.PurchaseDate == "" && dto.PurchaseDate != nil {
		exp += utils.CIPurchaseDateExp
	}

	if exists.PurchasePrice != 0 && dto.PurchasePrice == nil {
		exp -= utils.CIPurchasePriceExp
	} else if exists.PurchasePrice == 0 && dto.PurchasePrice != nil {
		exp += utils.CIPurchasePriceExp
	}

	if exists.Rating != 0 && dto.Rating == nil {
		exp -= utils.CIRatingExp
	} else if exists.Rating == 0 && dto.Rating != nil {
		exp += utils.CIRatingExp
	}

	if len(exists.Images) != 0 && len(dto.Images) == 0 {
		exp -= utils.PictureExp
	} else if len(exists.Images) == 0 && len(dto.Images) != 0 {
		exp += utils.PictureExp
	}
	if len(exists.CopyNumber) != 0 && len(dto.CopyNumber) == 0 {
		exp -= utils.CICopyNumberExp
	} else if len(exists.CopyNumber) == 0 && len(dto.CopyNumber) != 0 {
		exp += utils.CICopyNumberExp
	}

	err := s.userRepo.IncrementExperience(exists.UserLogin, exp)
	if err != nil {
		return err
	}

	return nil
}
