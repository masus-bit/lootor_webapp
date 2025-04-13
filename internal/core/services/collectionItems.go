package services

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
)

type CiService struct {
	repo           *repositories.CiRepository
	eventRepo      *repositories.EventsRepository
	collectionRepo *repositories.CollectionsRepository
	userRepo       *repositories.UsersRepository
	platformsRepo  *repositories.PlatformsRepository
	entityRepo     *repositories.EntitiesRepository
}

func NewCiService(repo *repositories.CiRepository, eventRepo *repositories.EventsRepository, collectionRepo *repositories.CollectionsRepository, userRepo *repositories.UsersRepository, platformsRepo *repositories.PlatformsRepository, entityRepo *repositories.EntitiesRepository) *CiService {
	return &CiService{repo: repo, eventRepo: eventRepo, collectionRepo: collectionRepo, userRepo: userRepo, platformsRepo: platformsRepo, entityRepo: entityRepo}
}

func (s *CiService) getEntities(entities []string) []models.Entities {
	var resultEntities []models.Entities
	for _, entity := range entities {
		dbEntity, _ := s.entityRepo.GetEntityByName(entity)
		resultEntities = append(resultEntities, *dbEntity)
	}
	return resultEntities
}

func (s *CiService) Create(dto models.CollectionItemsRequestCreate, authUserLogin string) (*models.CollectionItemsDataResponse, error) {

	platform, platformErr := s.platformsRepo.GetPlatformById(dto.Platform)
	if platformErr != nil {
		return nil, platformErr
	}

	entities := s.getEntities(dto.Entities)

	owner, ownerErr := s.userRepo.GetUserByLogin(authUserLogin)
	if ownerErr != nil {
		return nil, ownerErr
	}

	collection, _ := s.collectionRepo.GetByIdWithoutCollectionItems(dto.CollectionId)

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
		CopyNumber:    dto.CopyNumber,
		ShippingCost:  dto.ShippingCost,
		Entities:      entities,
		Platform:      *platform,
		Owner:         *owner,
		Collections:   collectionsSlice,
	}
	collectionItem, err := s.repo.CreateCI(dbCollectionItem)
	if err != nil {
		return nil, err
	}
	if !collection.IsPrivate {
		go func() {
			eventError := s.eventRepo.AddEvent(authUserLogin, utils.EventActionCreate, utils.EventTargetCollectionItem, dto.Name, nil, nil, &collectionItem.Id)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}()
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(collectionItem, &collectionItemResponse)
	collectionItemResponse.Collection = collection.Id
	if err != nil {
		return nil, err
	}

	return &models.CollectionItemsDataResponse{Data: collectionItemResponse}, nil
}

func (s *CiService) Delete(id string) (*dto.CommonResponse, error) {
	exists, _ := s.repo.GetCIByID(id)
	var result *dto.CommonResponse
	if exists != nil {
		if !exists.Collections[0].IsPrivate {
			go func() {
				eventError := s.eventRepo.AddEvent(exists.Owner.Login, utils.EventActionDelete, utils.EventTargetCollectionItem, exists.Name, nil, nil, &exists.Id)
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
		}
		err := s.repo.DeleteCI(id)
		if err != nil {
			return nil, err
		}
		result = &dto.CommonResponse{Data: dto.Resp{Success: true}}
	}
	return result, nil
}

func (s *CiService) Update(id string, dto models.CollectionItemsRequestCreate) (*models.CollectionItemsDataResponse, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	resultEntities := s.getEntities(dto.Entities)

	var dbCollectionItem models.CollectionItems
	err = mapstructure.Decode(dto, &dbCollectionItem)
	dbCollectionItem.Entities = resultEntities

	result, err := s.repo.UpdateCI(*exists, dbCollectionItem)
	if err != nil {
		return nil, err
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(result, &collectionItemResponse)

	return &models.CollectionItemsDataResponse{Data: collectionItemResponse}, nil
}

func (s *CiService) GetById(id string) (*models.CollectionItemsDataResponse, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	var collectionItemResponse models.CollectionItemsResponse
	err = mapstructure.Decode(exists, &collectionItemResponse)

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
			existsCollection := collection
			existsCollectionItems := collection.CollectionItems
			existsCollectionItems = append(existsCollectionItems, *collectionItem)
			collection.CollectionItems = existsCollectionItems
			_, errCol := s.collectionRepo.UpdateCollection(existsCollection, collection)
			if errCol != nil {
				return nil, err
			}
		}
		result = &dto.CommonResponse{Data: dto.Resp{Success: true}}
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
	_, errCol := s.collectionRepo.UpdateCollection(*existsSourceCollection, *sourceCollection)
	if errCol != nil {
		return nil, err
	}
	_, errCol = s.collectionRepo.UpdateCollection(*existsTargetCollection, *targetCollection)
	if errCol != nil {
		return nil, err
	}
	result = &dto.CommonResponse{Data: dto.Resp{Success: true}}

	return result, nil
}
