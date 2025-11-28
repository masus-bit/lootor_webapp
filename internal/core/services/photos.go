package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/photosclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/s3"
	"lootor/internal/pkg/utils"
	"strconv"
	"time"
)

type PhotosService struct {
	photosClient         *photosclient.GRPCPhotosClient
	userRepo             *repositories.UsersRepository
	tagsClient           *tagsclient.GRPCTagsClient
	collectionRepo       *repositories.CollectionsRepository
	eventsService        *EventsService
	s3Service            *s3.StorageService
	notificationsService *NotificationsService
}

func NewPhotosService(photosClient *photosclient.GRPCPhotosClient, userRepo *repositories.UsersRepository, tagsClient *tagsclient.GRPCTagsClient, collectionRepo *repositories.CollectionsRepository, eventsService *EventsService, s3Service *s3.StorageService, notificationsService *NotificationsService) *PhotosService {

	return &PhotosService{
		photosClient:         photosClient,
		userRepo:             userRepo,
		tagsClient:           tagsClient,
		collectionRepo:       collectionRepo,
		eventsService:        eventsService,
		s3Service:            s3Service,
		notificationsService: notificationsService,
	}
}

func (s *PhotosService) IncrementCommentsCount(ctx context.Context, id string) (*dto.CommonResponse, error) {
	resp, err := s.photosClient.IncrementCommentsCount(ctx, s.stringToUint64(id))
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PhotosService) DecrementCommentsCount(ctx context.Context, id string) (*dto.CommonResponse, error) {
	resp, err := s.photosClient.DecrementCommentsCount(ctx, s.stringToUint64(id))
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PhotosService) CreatePhoto(ctx context.Context, req *models.PhotoCreateRequest, authUserLogin string) (*models.PhotosDataResponse, error) {
	resp, err := s.photosClient.CreatePhotos(ctx, &microservices.CreatePhotosRequest{
		CollectionID:  req.CollectionID,
		Paths:         req.Paths,
		Author:        req.Author,
		AuthUserLogin: authUserLogin,
	})
	if err != nil {
		return nil, err
	}
	dbCollection, err := s.collectionRepo.GetByIdWithoutCollectionItems(req.CollectionID)
	if err != nil {
		return nil, err
	}
	if !dbCollection.IsPrivate {
		for _, ph := range resp.GetData() {
			stringId := strconv.Itoa(int(ph.GetId()))
			uuidColID, _ := uuid.Parse(req.CollectionID)
			eventError := s.eventsService.AddEvent(authUserLogin, utils.EventActionCreate, utils.EventTargetPhoto, ph.Path, &models.EventsParams{TargetPhotoID: stringId, TargetCollectionID: uuidColID})
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}
	}
	result := s.convertProtoToModels(resp.GetData(), authUserLogin)
	return &models.PhotosDataResponse{Data: result}, nil
}

func (s *PhotosService) GetByUser(ctx context.Context, userLogin, authUserLogin, limit, offset string) (*models.PhotosDataResponse, error) {
	resp, err := s.photosClient.FindByUser(ctx, &microservices.FindByUserRequest{
		AuthUserLogin: authUserLogin,
		Login:         userLogin,
		Offset:        offset,
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}
	result := s.convertProtoToModels(resp.GetData(), authUserLogin)
	return &models.PhotosDataResponse{Data: result}, nil
}

func (s *PhotosService) DeletePhoto(ctx context.Context, id string) (*dto.CommonResponse, error) {
	resp, err := s.photosClient.DeletePhoto(ctx, &microservices.DeletePhotoRequest{
		ID: s.stringToUint64(id),
	})
	if err != nil {
		return nil, err
	}
	paths := []string{resp.GetPath()}
	ok, err := s.s3Service.DeleteFiles(ctx, s3.DeleteFilesRequest{Keys: paths})
	if err != nil {
		return nil, err
	}
	if !ok.Success {
		return nil, errors.New("failed to delete photo")
	}
	go func() {
		_, _ = s.tagsClient.RemoveEntityTags(context.Background(), &microservices.RemoveEntityTagsRequest{
			EntityId: id,
		})
		eventError := s.eventsService.AddEvent("", utils.EventActionDelete, utils.EventTargetPhoto, "", &models.EventsParams{TargetPhotoID: id})
		if eventError != nil {
			log.Default().Print(eventError)
		}
	}()
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PhotosService) FindOneByID(ctx context.Context, id, authUserLogin string) (*models.PhotoDataResponse, error) {
	resp, err := s.photosClient.FindOneByID(ctx, &microservices.FindOneByIdRequest{
		ID: s.stringToUint64(id),
	})
	if err != nil {
		return nil, err
	}
	photo := s.convertProtoToModel(resp.GetData(), authUserLogin)
	return &models.PhotoDataResponse{Data: *photo}, nil
}

func (s *PhotosService) FindAllByCollectionID(ctx context.Context, collectionId, limit, offset, authUserLogin string) (*models.PhotosDataResponse, error) {
	resp, err := s.photosClient.FindAllByCollectionID(ctx, &microservices.FindByCollectionIdRequest{
		CollectionID:  collectionId,
		Limit:         limit,
		Offset:        offset,
		AuthUserLogin: authUserLogin,
	})
	if err != nil {
		return nil, err
	}
	photos := s.convertProtoToModels(resp.GetData(), authUserLogin)
	return &models.PhotosDataResponse{Data: photos, Total: resp.GetTotal()}, nil
}

func (s *PhotosService) LikePhoto(ctx context.Context, id string, authUserLogin string) (*dto.CommonResponse, error) {
	resp, err := s.photosClient.LikePhoto(ctx, &microservices.LikePhotoRequest{
		ID:             s.stringToUint64(id),
		InitiatorLogin: authUserLogin,
	})
	if err != nil {
		return nil, err
	}
	if resp.GetIsLike() {
		dbCollection, err := s.collectionRepo.GetByIdWithoutCollectionItems(resp.GetCollectionId())
		if err != nil {
			return nil, err
		}
		target := &dto.TargetItem{
			ID:               id,
			Name:             resp.GetPath(),
			Transliteration:  dbCollection.Transliteration,
			TargetType:       "photo",
			TargetParentName: dbCollection.Name,
		}
		go func() {
			err = s.notificationsService.SendNotification(context.Background(), &dto.NotificationsRequest{
				Login:       dbCollection.UserLogin,
				TargetID:    id,
				SenderLogin: authUserLogin,
				Type:        utils.NotificationTypePhoto,
				Action:      utils.NotificationActionLike,
				Date:        time.Now().Format(time.RFC3339),
				OwnerLogin:  dbCollection.UserLogin,
			}, target)
		}()
	} else {
		go func() {
			_ = s.notificationsService.DeleteNotification(context.Background(), id, authUserLogin)
		}()
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PhotosService) FindPhotosByIds(ctx context.Context, ids []string, authUserLogin string) (*models.PhotosDataResponse, error) {
	resp, err := s.photosClient.GetPhotosByIDs(ctx, &microservices.GetByIdsRequest{
		IDs:           ids,
		AuthUserLogin: authUserLogin,
	})
	if err != nil {
		return nil, err
	}
	photos := s.convertProtoToModels(resp.GetData(), authUserLogin)
	return &models.PhotosDataResponse{Data: photos}, nil
}

func (s *PhotosService) GetCountByCollection(ctx context.Context, collectionId string) (int64, error) {
	resp, err := s.photosClient.GetCountByCollection(ctx, &microservices.CountRequestPhoto{
		CollectionID: collectionId,
	})
	if err != nil {
		return 0, err
	}
	return resp.GetCount(), nil
}

func (s *PhotosService) UpdatePhoto(ctx context.Context, req *models.PhotoUpdateRequest, authUserLogin string) (*models.PhotoDataResponse, error) {
	resp, err := s.photosClient.UpdatePhoto(ctx, &microservices.UpdatePhotoRequest{
		ID:          s.stringToUint64(req.ID),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	dbCollection, err := s.collectionRepo.GetByIdWithoutCollectionItems(resp.GetData().GetCollectionId())
	if err != nil {
		return nil, err
	}

	if len(req.Tags) > 0 {
		stringedUint := strconv.Itoa(int(resp.GetData().GetId()))
		_, _ = s.tagsClient.AddTagsToEntity(context.Background(), &microservices.AddFewTagsToEntityRequest{
			EntityType: "photo",
			EntityId:   stringedUint,
			TagIds:     req.Tags,
			Author:     authUserLogin,
			ShowSearch: !dbCollection.IsPrivate,
		})
	}

	photo := s.convertProtoToModel(resp.GetData(), authUserLogin)

	if len(req.Tags) > 0 {
		if !dbCollection.IsPrivate {
			for _, tag := range photo.Tags {
				tagUUID, _ := uuid.Parse(tag.ID)
				stringPhotoID := strconv.Itoa(int(resp.GetData().GetId()))
				eventError := s.eventsService.AddEvent(authUserLogin, utils.EventActionAddTag, utils.EventTargetTag, tag.Name, &models.EventsParams{TargetTagID: tagUUID, TargetPhotoID: stringPhotoID, TagRelatedEntityType: utils.EventTargetPhoto})
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}
		}
	}

	return &models.PhotoDataResponse{Data: *photo}, nil
}

func (s *PhotosService) convertProtoToModels(photos []*microservices.PhotoItem, authUserLogin string) []models.Photos {
	var resultPhotos []models.Photos
	var userLogins []string
	var collectionIds []string
	var photoIds []string
	for i := range photos {
		userLogins = append(userLogins, photos[i].GetAuthor())
		collectionIds = append(collectionIds, photos[i].GetCollectionId())
		id := strconv.Itoa(int(photos[i].GetId()))
		photoIds = append(photoIds, id)
	}
	collectionsMap, err := s.collectionRepo.GetCollectionsByIdsMapForShort(collectionIds)
	if err != nil {
		return nil
	}
	usersMap, err := s.userRepo.GetForSubsMap(userLogins)
	if err != nil {
		return nil
	}
	tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
		EntityIds: photoIds,
	})
	if err != nil {
		return nil
	}
	for i := range photos {
		photoTags := tags.GetTags()[photoIds[i]]
		var resultPhotoTags []models.ShortTags

		for _, tag := range photoTags.GetTags() {
			resultPhotoTags = append(resultPhotoTags, models.ShortTags{
				ID:        tag.GetId(),
				Name:      tag.GetName(),
				Slug:      tag.GetSlug(),
				PrimaryID: tag.GetPrimaryId(),
				SeriesID:  tag.GetSeriesId(),
			})
		}

		createdAtAsTime, err := parseDate(photos[i].GetCreatedAt())
		if err != nil {
			log.Printf("Error parsing date: %v", err)
		}
		canLike := utils.CanLike(photos[i].GetLikes(), authUserLogin, photos[i].GetAuthor())
		resultPhotos = append(resultPhotos, models.Photos{
			ID:            photos[i].GetId(),
			Author:        usersMap[photos[i].GetAuthor()],
			CollectionID:  photos[i].GetCollectionId(),
			Path:          photos[i].GetPath(),
			LikesCount:    int64(len(photos[i].GetLikes())),
			CommentsCount: photos[i].GetCommentsCount(),
			CreatedAt:     createdAtAsTime,
			Collection:    collectionsMap[photos[i].GetCollectionId()],
			Tags:          resultPhotoTags,
			CanLike:       canLike,
			IsOwner:       authUserLogin == photos[i].GetAuthor(),
			Description:   photos[i].GetDescription(),
		})
	}
	return resultPhotos
}

func (s *PhotosService) convertProtoToModel(photo *microservices.PhotoItem, authUserLogin string) *models.Photos {
	var resultPhoto *models.Photos
	var userLogins []string
	var collectionIds []string
	var photoIds []string
	userLogins = append(userLogins, photo.GetAuthor())
	collectionIds = append(collectionIds, photo.GetCollectionId())
	id := strconv.Itoa(int(photo.GetId()))
	photoIds = append(photoIds, id)
	collectionsMap, err := s.collectionRepo.GetCollectionsByIdsMapForShort(collectionIds)
	if err != nil {
		return nil
	}
	usersMap, err := s.userRepo.GetForSubsMap(userLogins)
	if err != nil {
		return nil
	}
	tags, err := s.tagsClient.GetTagsByEntityIdsMap(context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
		EntityIds: photoIds,
	})
	if err != nil {
		return nil
	}
	photoTags := tags.GetTags()[photoIds[0]]
	var resultPhotoTags []models.ShortTags

	for _, tag := range photoTags.GetTags() {
		resultPhotoTags = append(resultPhotoTags, models.ShortTags{
			ID:        tag.GetId(),
			Name:      tag.GetName(),
			Slug:      tag.GetSlug(),
			PrimaryID: tag.GetPrimaryId(),
			SeriesID:  tag.GetSeriesId(),
		})
	}
	canLike := utils.CanLike(photo.GetLikes(), authUserLogin, photo.GetAuthor())
	createdAtAsTime, err := parseDate(photo.GetCreatedAt())
	if err != nil {
		log.Printf("Error parsing date: %v", err)
	}
	resultPhoto = &models.Photos{
		ID:            photo.GetId(),
		Author:        usersMap[photo.GetAuthor()],
		CollectionID:  photo.GetCollectionId(),
		Path:          photo.GetPath(),
		LikesCount:    int64(len(photo.GetLikes())),
		CommentsCount: photo.GetCommentsCount(),
		CreatedAt:     createdAtAsTime,
		Collection:    collectionsMap[photo.GetCollectionId()],
		Tags:          resultPhotoTags,
		CanLike:       canLike,
		IsOwner:       authUserLogin == photo.GetAuthor(),
		Description:   photo.GetDescription(),
	}

	return resultPhoto
}

func (s *PhotosService) stringToUint64(id string) uint64 {
	val, _ := strconv.ParseUint(id, 10, 64)
	return val
}

func parseDate(dateStr string) (string, error) {
	layout := "2006-01-02 15:04:05.999999 -0700 MST"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", err
	}

	return t.UTC().Format(time.RFC3339), nil
}
