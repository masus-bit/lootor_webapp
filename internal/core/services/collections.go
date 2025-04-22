package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/s3"
	"lootor/internal/pkg/utils"
	"slices"
	"time"
)

type CollectionService struct {
	repo               *repositories.CollectionsRepository
	tagsRepo           *repositories.TagsRepository
	userRepo           *repositories.UsersRepository
	eventRepo          *repositories.EventsRepository
	collectionItemRepo *repositories.CiRepository
	s3Service          *s3.S3Service
}

func NewCollectionService(repo *repositories.CollectionsRepository, tagsRepo *repositories.TagsRepository, userRepo *repositories.UsersRepository, eventRepo *repositories.EventsRepository, collectionItemRepo *repositories.CiRepository, seService *s3.S3Service) *CollectionService {
	return &CollectionService{repo: repo, tagsRepo: tagsRepo, userRepo: userRepo, eventRepo: eventRepo, collectionItemRepo: collectionItemRepo, s3Service: seService}
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
	if !dto.IsPrivate {
		eventError := s.eventRepo.AddEvent(dto.UserLogin, utils.EventActionCreate, utils.EventTargetCollection, dto.Name, nil, &dbCollection.Id, nil)
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
	return &models.CollectionDataResponse{Data: response}, nil
}

func (s *CollectionService) Update(id string, dto *models.CollectionCreateRequest) (*models.CollectionDataResponse, error) {
	exists, err := s.repo.GetCollectionById(id)
	if err != nil {
		return nil, err
	}
	processTags, _ := s.processTags(dto.Tags)

	var dbCollection models.Collections
	err = mapstructure.Decode(dto, &dbCollection)
	dbCollection.Tags = processTags
	resultCollection, err := s.repo.UpdateCollection(exists, &dbCollection)
	resultCollection.Tags = processTags
	if err != nil {
		return nil, err
	}

	if !dto.IsPrivate {
		eventError := s.eventRepo.AddEvent(dto.UserLogin, utils.EventActionCreate, utils.EventTargetCollection, dto.Name, nil, &dbCollection.Id, nil)
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}
	var finalCollection models.CollectionsResponse
	err = mapstructure.Decode(resultCollection, &finalCollection)
	if err != nil {
		return nil, err
	}
	finalCollection.CanLike = true
	finalCollection.CanSubscribe = true

	return &models.CollectionDataResponse{Data: finalCollection}, nil
}

func (s *CollectionService) Delete(id string, ctx context.Context) (*dto.CommonResponse, error) {
	exists, _ := s.repo.GetCollectionById(id)
	if !exists.IsPrivate {
		go func() {
			eventError := s.eventRepo.AddEvent(exists.User.Login, utils.EventActionDelete, utils.EventTargetCollection, exists.Name, nil, &exists.Id, nil)
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
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, s.repo.DeleteCollection(id)
}

func (s *CollectionService) GetByUserLogin(login string, authorizedUser string) (*models.AllCollectionsDataResponse, error) {
	var collections []models.Collections
	if authorizedUser == login {
		collections, _ = s.repo.GetByUserIdWithoutCollectionItems(login)
	} else {
		collections, _ = s.repo.GetByUserIdWithoutPrivates(login)
	}
	var authUser *models.Users
	var subArray []string
	if authorizedUser != "" {
		authUser, _ = s.userRepo.GetUserByLogin(authorizedUser)
		subArray = authUser.CollectionSubscriptions
	}
	result := make([]models.CollectionsResponse, 0)
	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		err := mapstructure.Decode(dbCollection, &temp)
		temp.ShareString = utils.DefineShareString(authorizedUser, login, &dbCollection)
		temp.CollectionItemsCount, _ = s.collectionItemRepo.GetCountCI(dbCollection.Id)
		temp.TotalPrice, _ = s.collectionItemRepo.Sum(dbCollection.Id)
		temp.ShippingTotal, _ = s.collectionItemRepo.SumShippingCost(dbCollection.Id)
		temp.CanSubscribe = !slices.Contains(subArray, dbCollection.Id.String())
		temp.LikesCount = int64(len(dbCollection.Likes))
		temp.CanLike = true
		if authorizedUser != "" {
			temp.CanLike = !slices.Contains(dbCollection.Likes, authorizedUser)
		}

		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}

	return &models.AllCollectionsDataResponse{Data: result}, nil
}

func (s *CollectionService) GetOne(authorizerUser string, id string, transliteration string, userLogin string, shareString string) (*models.CollectionDataResponse, error) {
	var dbCollection *models.Collections

	if id != "" {
		dbCollection, _ = s.repo.GetCollectionById(id)
	} else if userLogin != "" {
		dbCollection, _ = s.repo.GetOneByTransliteration(userLogin, transliteration)
	} else if shareString != "" {
		dbCollection, _ = s.repo.GetByShareString(shareString)
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
		if authorizerUser != "" {
			if authorizerUser == item.Owner.Login {
				temp.CanLike = true
			} else {
				temp.CanLike = !slices.Contains(item.Likes, authorizerUser)
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
	if authUser != nil {

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
			eventError := s.eventRepo.AddEvent(exists.User.Login, utils.EventActionCreate, utils.EventTargetCollection, exists.Name, nil, &exists.Id, nil)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}
	} else {
		exists.Likes = utils.RemoveByValue(exists.Likes, userLogin)
		if !exists.IsPrivate {
			go func() {
				eventError := s.eventRepo.AddEvent(exists.User.Login, utils.EventActionDelete, utils.EventTargetCollection, exists.Name, nil, &exists.Id, nil)
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}()
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
	dbCollection, err := s.repo.GetCollectionById(targetId)
	if err != nil {
		return nil, err
	}
	if isSubscribe {
		subscriber.CollectionSubscriptions = append(subscriber.CollectionSubscriptions, targetId)
		dbCollection.SubscribersCount = dbCollection.SubscribersCount + 1
	} else {
		subscriber.CollectionSubscriptions = utils.RemoveByValue(subscriber.CollectionSubscriptions, targetId)
		dbCollection.SubscribersCount = dbCollection.SubscribersCount - 1
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

func (s *CollectionService) GetByTag(tag string, authUserLogin string) (*models.AllCollectionsDataResponse, error) {
	collections, err := s.repo.GetCollectionByTag(tag)
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

	for _, dbCollection := range collections {
		var temp models.CollectionsResponse
		errMap := mapstructure.Decode(dbCollection, &temp)
		temp.ShareString = utils.DefineShareString(authUserLogin, dbCollection.User.Login, &dbCollection)
		temp.CollectionItemsCount, _ = s.collectionItemRepo.GetCountCI(dbCollection.Id)
		temp.TotalPrice, _ = s.collectionItemRepo.Sum(dbCollection.Id)
		temp.ShippingTotal, _ = s.collectionItemRepo.SumShippingCost(dbCollection.Id)
		temp.CanSubscribe = !slices.Contains(subArray, dbCollection.Id.String())
		if errMap != nil {
			return nil, errMap
		}
		result = append(result, temp)
	}

	return &models.AllCollectionsDataResponse{Data: result}, nil
}
