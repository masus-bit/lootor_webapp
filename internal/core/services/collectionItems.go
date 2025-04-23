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
	"slices"
)

type CiService struct {
	repo           *repositories.CiRepository
	eventRepo      *repositories.EventsRepository
	collectionRepo *repositories.CollectionsRepository
	userRepo       *repositories.UsersRepository
	platformsRepo  *repositories.PlatformsRepository
	entityRepo     *repositories.EntitiesRepository
	s3Service      *s3.S3Service
	itemTypeRepo   *repositories.ItemTypesRepository
}

func NewCiService(repo *repositories.CiRepository, eventRepo *repositories.EventsRepository, collectionRepo *repositories.CollectionsRepository, userRepo *repositories.UsersRepository, platformsRepo *repositories.PlatformsRepository, entityRepo *repositories.EntitiesRepository, s3Service *s3.S3Service, itemTypeRepo *repositories.ItemTypesRepository) *CiService {
	return &CiService{repo: repo, eventRepo: eventRepo, collectionRepo: collectionRepo, userRepo: userRepo, platformsRepo: platformsRepo, entityRepo: entityRepo, s3Service: s3Service, itemTypeRepo: itemTypeRepo}
}

func (s *CiService) getEntities(entities []string) []models.Entities {
	var resultEntities []models.Entities
	for _, entity := range entities {
		dbEntity, _ := s.entityRepo.GetEntityByName(entity)
		resultEntities = append(resultEntities, *dbEntity)
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

	entities := s.getEntities(dto.Entities)
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
			_, err := s.s3Service.DeleteFiles(ctx, s3.DeleteFilesRequest{Keys: exists.Images})
			if err != nil {
				return nil, err
			}
		}
		result = &dto.CommonResponse{Data: dto.Resp{Success: true}}
	}
	return result, nil
}

func (s *CiService) Update(id string, dto *models.CollectionItemsRequestCreate) (*models.CollectionItemsDataResponse, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	resultEntities := s.getEntities(dto.Entities)

	var dbCollectionItem models.CollectionItems
	err = mapstructure.Decode(dto, &dbCollectionItem)
	dbCollectionItem.Entities = resultEntities

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

	dbCollectionItem.Platform = platform
	dbCollectionItem.PlatformID = platformID
	dbCollectionItem.ItemTypeID = itemTypeID
	dbCollectionItem.ItemType = itemType
	result, err := s.repo.UpdateCI(exists, &dbCollectionItem)
	if err != nil {
		return nil, err
	}

	if !exists.Collections[0].IsPrivate {
		eventError := s.eventRepo.AddEvent(exists.Owner.Login, utils.EventActionUpdate, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id, nil)
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(result, &collectionItemResponse)

	colId, _ := uuid.Parse(dto.Collection)

	collectionItemResponse.Collection = colId
	collectionItemResponse.Owner = exists.Owner
	collectionItemResponse.CanLike = true
	collectionItemResponse.LikesCount = int64(len(exists.Likes))
	return &models.CollectionItemsDataResponse{Data: collectionItemResponse}, nil
}

func (s *CiService) GetById(id string, authUser string) (*models.CollectionItemsDataResponse, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(exists, &collectionItemResponse)
	collectionItemResponse.Collection = exists.Collections[0].Id
	collectionItemResponse.Owner = exists.Owner
	collectionItemResponse.LikesCount = int64(len(exists.Likes))
	collectionItemResponse.CanLike = true
	if authUser != "" {
		if authUser == exists.Owner.Login {
			collectionItemResponse.CanLike = true
		} else {
			collectionItemResponse.CanLike = !slices.Contains(exists.Likes, authUser)
		}

	}

	return &models.CollectionItemsDataResponse{Data: collectionItemResponse}, nil
}

func (s *CiService) CopyOrMove(id string, targetCollectionIds []string, sourceCollectionId string) (*dto.CommonResponse, error) {
	var result *dto.CommonResponse
	if sourceCollectionId == "" {
		collections := make([]models.Collections, 0)
		for _, collectionId := range targetCollectionIds {
			collection, _ := s.collectionRepo.GetCollectionById(collectionId)
			collections = append(collections, *collection)
		}
		collectionItem, err := s.repo.GetCIByID(id)
		if err != nil {
			return nil, err
		}

		for _, collection := range collections {
			collectionItem.Collections = append(collectionItem.Collections, collection)
			//existsCollection := collection
			//existsCollectionItems := collection.CollectionItems
			//existsCollectionItems = append(existsCollectionItems, *collectionItem)
			//collection.CollectionItems = existsCollectionItems
			//_, errCol := s.collectionRepo.UpdateCollection(&existsCollection, &collection)
			//if errCol != nil {
			//	return nil, err
			//}
			_, erroring := s.repo.UpdateCI(collectionItem, collectionItem)
			if erroring != nil {
				return nil, erroring
			}
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}

	sourceCollection, err := s.collectionRepo.GetCollectionById(sourceCollectionId)
	if err != nil {
		return nil, err
	}
	targetCollection, err := s.collectionRepo.GetCollectionById(targetCollectionIds[0])
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
	existsSourceCollection.CollectionItems = sourceCollectionItems
	existsTargetCollection.CollectionItems = targetCollectionItems
	_, errCol := s.collectionRepo.UpdateCollection(existsSourceCollection, sourceCollection)
	if errCol != nil {
		return nil, err
	}
	_, errCol = s.collectionRepo.UpdateCollection(existsTargetCollection, targetCollection)
	if errCol != nil {
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
			eventError := s.eventRepo.AddEvent(exists.Owner.Login, utils.EventActionLike, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id, nil)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}
	} else {
		exists.Likes = utils.RemoveByValue(exists.Likes, userLogin)
		if !exists.Collections[0].IsPrivate {
			go func() {
				eventError := s.eventRepo.AddEvent(exists.Owner.Login, utils.EventActionLike, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id, nil)
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
		}
	}

	_, err = s.repo.UpdateCI(exists, exists)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}
