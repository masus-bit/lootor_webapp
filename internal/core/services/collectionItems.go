package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/s3"
	"lootor/internal/pkg/utils"
	"reflect"
	"slices"
	"time"
)

type CiService struct {
	repo                 *repositories.CiRepository
	eventsService        *EventsService
	collectionRepo       *repositories.CollectionsRepository
	userRepo             *repositories.UsersRepository
	platformsRepo        *repositories.PlatformsRepository
	s3Service            *s3.StorageService
	itemTypeRepo         *repositories.ItemTypesRepository
	notificationsService *NotificationsService
	tagsClient           *tagsclient.GRPCTagsClient
}

func NewCiService(repo *repositories.CiRepository, eventsService *EventsService, collectionRepo *repositories.CollectionsRepository, userRepo *repositories.UsersRepository, platformsRepo *repositories.PlatformsRepository, s3Service *s3.StorageService, itemTypeRepo *repositories.ItemTypesRepository, notificationsService *NotificationsService, tagsClient *tagsclient.GRPCTagsClient) *CiService {
	return &CiService{repo: repo, eventsService: eventsService, collectionRepo: collectionRepo, userRepo: userRepo, platformsRepo: platformsRepo, s3Service: s3Service, itemTypeRepo: itemTypeRepo, notificationsService: notificationsService, tagsClient: tagsClient}
}

func (s *CiService) Create(dto *models.CollectionItemsRequestCreate, authUserLogin string) (*models.CollectionItemsDataResponse, error) {
	var platform *models.Platforms
	var platformID *uuid.UUID
	var itemType *models.ItemTypes
	var itemTypeID *uuid.UUID

	if dto.Platform != "" {
		foundPlatform, err := s.platformsRepo.GetPlatformByID(dto.Platform)
		if err != nil {
			return nil, fmt.Errorf("error getting platform: %v", err)
		}
		platform = foundPlatform
		platformID = &foundPlatform.ID
	} else {
		platform = nil
		platformID = nil
	}

	if dto.ItemType != "" {
		foundItemType, err := s.itemTypeRepo.GetTypeByID(dto.ItemType)
		if err != nil {
			return nil, fmt.Errorf("error getting item type: %v", err)
		}
		itemType = foundItemType
		itemTypeID = &foundItemType.ID
	} else {
		itemType = nil
		itemTypeID = nil
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
		eventError := s.eventsService.AddEvent(authUserLogin, utils.EventActionCreate, utils.EventTargetCollectionItem, dto.Name, &models.EventsParams{TargetItemID: collectionItem.ID})
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(collectionItem, &collectionItemResponse)
	collectionItemResponse.Collection = collection.ID
	if err != nil {
		return nil, err
	}
	collectionItemResponse.LikesCount = int64(len(collectionItem.Likes))
	collectionItemResponse.CanLike = true
	collectionItemResponse.IsOwner = true

	exp := utils.CIExp

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

	var tags *microservices.GetShortsResponse

	if len(dto.Tags) > 0 {
		tags, _ = s.tagsClient.AddTagsToEntity(context.Background(), &microservices.AddFewTagsToEntityRequest{
			EntityType: "collectionItem",
			EntityId:   collectionItemResponse.ID.String(),
			TagIds:     dto.Tags,
			Author:     authUserLogin,
			ShowSearch: !collection.IsPrivate,
		})

		if !collection.IsPrivate {
			for _, tag := range tags.GetTags() {
				tagUUID, _ := uuid.Parse(tag.GetId())
				eventError := s.eventsService.AddEvent(collectionItemResponse.Owner.Login, utils.EventActionAddTag, utils.EventTargetTag, tag.Name, &models.EventsParams{TargetTagID: tagUUID, TargetItemID: collectionItemResponse.ID, TagRelatedEntityType: utils.EventTargetCollectionItem})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}
		}
	}

	var resultTags []models.ShortTags
	resultTags = make([]models.ShortTags, 0)
	if tags != nil {
		resultTags = utils.NormalizeTagsShort(tags.GetTags())
		collectionItemResponse.Tags = resultTags
	}
	return &models.CollectionItemsDataResponse{Data: collectionItemResponse}, nil
}

func (s *CiService) Delete(id string, ctx context.Context) (*dto.CommonResponse, error) {
	exists, _ := s.repo.GetCIByID(id)
	var result *dto.CommonResponse
	if exists != nil {
		if !exists.Collections[0].IsPrivate {
			eventError := s.eventsService.AddEvent(exists.Owner.Login, utils.EventActionDelete, utils.EventTargetCollectionItem, exists.Name, &models.EventsParams{TargetItemID: exists.ID})
			if eventError != nil {
				log.Default().Print(eventError)
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

	err = s.updateCIExp(exists, dto)
	if err != nil {
		return nil, err
	}

	dst := reflect.ValueOf(exists).Elem()
	src := reflect.ValueOf(dto).Elem()

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		fieldName := src.Type().Field(i).Name

		if fieldName == "Platform" || fieldName == "ItemType" {
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
	tags, err := s.tagsClient.UpdateTagsOfEntity(context.Background(), &microservices.UpdateTagsOfEntityRequest{
		EntityType: "collectionItem",
		EntityId:   id,
		TagIds:     dto.Tags,
		Author:     exists.Owner.Login,
	})
	if err != nil {
		return nil, err
	}
	if len(dto.Tags) > 0 {
		if !exists.Collections[0].IsPrivate {
			for _, tag := range tags.GetTags() {
				tagUUID, _ := uuid.Parse(tag.GetId())
				eventError := s.eventsService.AddEvent(exists.UserLogin, utils.EventActionAddTag, utils.EventTargetTag, tag.Name, &models.EventsParams{TargetTagID: tagUUID, TargetItemID: exists.ID, TagRelatedEntityType: utils.EventTargetCollectionItem})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}
		}
	}

	var platform *models.Platforms
	var platformID *uuid.UUID
	var itemType *models.ItemTypes
	var itemTypeID *uuid.UUID

	if dto.Platform != nil && *dto.Platform != "" {
		foundPlatform, err := s.platformsRepo.GetPlatformByID(*dto.Platform)
		if err != nil {
			return nil, fmt.Errorf("error getting platform: %v", err)
		}
		platform = foundPlatform
		platformID = &foundPlatform.ID
	} else {
		platform = nil
		platformID = nil
	}

	if dto.ItemType != nil && *dto.ItemType != "" {
		foundItemType, err := s.itemTypeRepo.GetTypeByID(*dto.ItemType)
		if err != nil {
			return nil, fmt.Errorf("error getting item type: %v", err)
		}
		itemType = foundItemType
		itemTypeID = &foundItemType.ID
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
		eventError := s.eventsService.AddEvent(exists.Owner.Login, utils.EventActionUpdate, utils.EventTargetCollectionItem, exists.Name, &models.EventsParams{TargetItemID: exists.ID})
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}
	return result, nil
}

func (s *CiService) GetByID(id string, authUser string) (*models.CollectionItems, error) {
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
	sourceCollectionItems := utils.RemoveByValueStruct(sourceCollection.CollectionItems, collectionItem.ID)
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
				eventError := s.eventsService.AddEvent(userLogin, utils.EventActionLike, utils.EventTargetCollectionItem, exists.Name, &models.EventsParams{TargetItemID: exists.ID})
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
			ID:               id,
			Name:             exists.Name,
			Transliteration:  exists.Transliteration,
			TargetType:       "collectionItem",
			TargetParentName: exists.Collections[0].Name,
		}
		go func() {
			err = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
				Login:       exists.UserLogin,
				TargetID:    id,
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
		go func() {
			_ = s.notificationsService.DeleteNotification(context.Background(), id, userLogin)
			eventError := s.eventsService.AddEvent(userLogin, utils.EventActionDislike, utils.EventTargetCollectionItem, exists.Name, &models.EventsParams{TargetItemID: exists.ID})
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}()
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

func (s *CiService) GetAll(limit, offset, search string) (*models.CollectionItemsDataPoor, error) {
	collectionItems, err := s.repo.GetAll(limit, offset, search)
	if err != nil {
		return nil, err
	}

	var collectionItemIds []string

	for _, item := range collectionItems {
		collectionItemIds = append(collectionItemIds, item.ID.String())
	}

	tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{EntityIds: collectionItemIds})
	if err != nil {
		return nil, err
	}

	var collectionItemsAll []models.CollectionItemsResponse

	for _, item := range collectionItems {
		var temp models.CollectionItemsResponse
		err = mapstructure.Decode(item, &temp)
		if err != nil {
			return nil, err
		}
		itemTags := tags.GetTags()[item.ID.String()]

		var resultTags []models.ShortTags = utils.NormalizeTagsShort(itemTags.GetTags())

		if len(item.Collections) > 0 {
			temp.Collection = item.Collections[0].ID
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
		temp.Tags = resultTags
		temp.Owner = item.Owner
		collectionItemsAll = append(collectionItemsAll, temp)
	}

	return &models.CollectionItemsDataPoor{Data: collectionItemsAll}, nil
}

func (s *CiService) updateCIExp(exists *models.CollectionItems, dto *models.CollectionItemsRequestUpdate) error {
	var exp int

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
