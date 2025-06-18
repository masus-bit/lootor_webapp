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
	"strings"
	"unicode"
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

func slugify(input string) string {
	translitMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "a", 'Б': "b", 'В': "v", 'Г': "g", 'Д': "d", 'Е': "e", 'Ё': "yo",
		'Ж': "zh", 'З': "z", 'И': "i", 'Й': "y", 'К': "k", 'Л': "l", 'М': "m",
		'Н': "n", 'О': "o", 'П': "p", 'Р': "r", 'С': "s", 'Т': "t", 'У': "u",
		'Ф': "f", 'Х': "kh", 'Ц': "ts", 'Ч': "ch", 'Ш': "sh", 'Щ': "shch",
		'Ъ': "", 'Ы': "y", 'Ь': "", 'Э': "e", 'Ю': "yu", 'Я': "ya",
	}

	var result strings.Builder
	hasRussian := false

	// Проверяем, есть ли русские буквы в строке
	for _, char := range input {
		if unicode.Is(unicode.Cyrillic, char) {
			hasRussian = true
			break
		}
	}

	for _, char := range input {
		switch {
		case char == ' ':
			result.WriteString("_")
		case hasRussian && unicode.Is(unicode.Cyrillic, char):
			if val, ok := translitMap[char]; ok {
				result.WriteString(val)
			}
		case unicode.IsUpper(char):
			result.WriteRune(unicode.ToLower(char))
		default:
			result.WriteRune(char)
		}
	}

	return strings.ToLower(result.String())
}

func (s *CiService) getEntities(entities []string) []models.Entities {
	var resultEntities []models.Entities
	for _, entity := range entities {
		entityByTranslit, err := s.entityRepo.GetEntityByName(entity)
		if err != nil {
			fmt.Errorf("failed to get entity: %w", err)
		}
		if entityByTranslit == nil {
			newEntity := &models.Entities{Name: entity, Transliteration: slugify(entity)}
			entityByTranslit, err = s.entityRepo.CreateEntity(newEntity)
			if err != nil {
				fmt.Errorf("failed to add entity: %w", err)
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
	collectionItemResponse.IsOwner = true

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

func (s *CiService) Update(id string, dto *models.CollectionItemsRequestUpdate) (*models.CollectionItems, error) {
	exists, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	var resultEntities []models.Entities

	if len(dto.Entities) != 0 {
		resultEntities = s.getEntities(dto.Entities)
	} else {
		resultEntities = make([]models.Entities, 0)
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

	if dto.Platform != nil {
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

	if *dto.ItemType != "" {
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

	return &models.CollectionItemsDataSortedResponse{Data: sortedCollectionItems, Total: totalCount}, nil
}

func (s *CiService) GetByEntityAndType(entity string, itemType string, limit string, offset string, search string, orderBy string, order string, authUser string) ([]models.CollectionItems, int64, error) {
	dbEntity, err := s.entityRepo.GetEntityByTranslit(entity)
	if dbEntity == nil {
		return nil, 0, fmt.Errorf("entity not found")
	}
	collectionItems, totalCount, err := s.repo.GetByEntityAndType(entity, itemType, limit, offset, search, orderBy, order)
	if err != nil {
		return nil, totalCount, err
	}

	return collectionItems, totalCount, err
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
