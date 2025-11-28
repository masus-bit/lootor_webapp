package services

import (
	"context"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/notificationsclient"
	"lootor/internal/infrastructure/photosclient"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/pkg/dto"
	"strconv"
)

type NotificationsService struct {
	notificationsClient *notificationsclient.GRPCNotificationsClient
	userRepo            *repositories.UsersRepository
	collectionRepo      *repositories.CollectionsRepository
	ciRepo              *repositories.CiRepository
	postsClient         *postsclient.GRPCPostsClient
	photosClient        *photosclient.GRPCPhotosClient
}

func NewNotificationsService(notificationsClient *notificationsclient.GRPCNotificationsClient, userRepo *repositories.UsersRepository, collectionRepo *repositories.CollectionsRepository, ciRepo *repositories.CiRepository, postsClient *postsclient.GRPCPostsClient, photosClient *photosclient.GRPCPhotosClient) *NotificationsService {
	return &NotificationsService{
		notificationsClient: notificationsClient, userRepo: userRepo, collectionRepo: collectionRepo, ciRepo: ciRepo, postsClient: postsClient, photosClient: photosClient}
}

func (s *NotificationsService) SendNotification(ctx context.Context, request *dto.NotificationsRequest, targetReq *dto.TargetItem) error {
	targetUser, _ := s.userRepo.GetUserByLogin(request.Login)
	senderUser, _ := s.userRepo.GetUserByLogin(request.SenderLogin)
	owner, _ := s.userRepo.GetUserByLogin(request.OwnerLogin)

	targetUserReq := dto.User{
		Login:       targetUser.Login,
		IsPremium:   targetUser.IsPremium,
		AvatarURL:   targetUser.AvatarURL,
		ProfileName: targetUser.ProfileName,
	}

	senderUserReq := dto.User{
		Login:       senderUser.Login,
		IsPremium:   senderUser.IsPremium,
		AvatarURL:   senderUser.AvatarURL,
		ProfileName: senderUser.ProfileName,
	}

	ownerUserReq := dto.User{
		Login:       owner.Login,
		IsPremium:   owner.IsPremium,
		AvatarURL:   owner.AvatarURL,
		ProfileName: owner.ProfileName,
	}

	request.TargetUser = targetUserReq
	request.SenderUser = senderUserReq
	request.TargetItem = *targetReq
	request.Owner = ownerUserReq

	finalRequest := dto.NotificationsRequest{
		Login:       request.Login,
		TargetID:    request.TargetID,
		Type:        request.Type,
		SenderLogin: request.SenderLogin,
		Date:        request.Date,
		TargetUser:  request.TargetUser,
		SenderUser:  request.SenderUser,
		TargetItem:  request.TargetItem,
		Action:      request.Action,
		Owner:       request.Owner,
	}

	_, err := s.notificationsClient.AddNotification(ctx, &finalRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationsService) GetAllNotifications(ctx context.Context, login, limit, offset, authUser string) (*dto.NotificationsDataResponse, error) {
	notifications, err := s.notificationsClient.GetAllNotifications(ctx, login, limit, offset)
	if err != nil {
		return nil, err
	}
	targetUser, _ := s.userRepo.GetUserByLogin(login)
	var photosIDs []string
	var collectionsIDs []string
	var collectionItemsIDs []string
	var postsIDs []string
	var userLogins []string

	var photosMap map[string]models.Photos
	var collectionsMap map[string]models.Collections
	var collectionItemsMap map[string]models.CollectionItems
	var postsMap map[string]models.Posts
	var resultNotifications []dto.Notifications
	photosMap = make(map[string]models.Photos)
	collectionsMap = make(map[string]models.Collections)
	collectionItemsMap = make(map[string]models.CollectionItems)
	postsMap = make(map[string]models.Posts)
	for _, n := range notifications.GetData() {
		switch n.GetTargetType() {
		case "photo":
			photosIDs = append(photosIDs, n.GetTargetId())
		case "collection":
			collectionsIDs = append(collectionsIDs, n.GetTargetId())
		case "collectionItem":
			collectionItemsIDs = append(collectionItemsIDs, n.GetTargetId())
		case "post":
			postsIDs = append(postsIDs, n.GetTargetId())
		}
	}
	if len(photosIDs) > 0 {
		photos, err := s.photosClient.GetPhotosByIDsMap(ctx, &microservices.GetByIdsRequest{
			Ids: photosIDs,
		})
		if err != nil {
			return nil, err
		}
		for _, photo := range photos.GetData() {
			userLogins = append(userLogins, photo.GetAuthor())
			stringID := strconv.Itoa(int(photo.GetId()))
			photosMap[stringID] = models.Photos{
				ID:           photo.GetId(),
				CollectionID: photo.GetCollectionId(),
				Path:         photo.GetPath(),
				Author: models.SubUsers{
					Login: photo.GetAuthor(),
				},
			}
		}
	}

	if len(collectionsIDs) > 0 {
		collections, err := s.collectionRepo.GetCollectionsByIdsMapForShort(collectionsIDs)
		if err != nil {
			return nil, err
		}
		for _, collection := range collections {
			userLogins = append(userLogins, collection.UserLogin)
			collectionsMap[collection.ID.String()] = models.Collections{
				ID:              collection.ID,
				Name:            collection.Name,
				Transliteration: collection.Transliteration,
				UserLogin:       collection.UserLogin,
			}
		}
	}

	if len(collectionItemsIDs) > 0 {
		collectionItems, err := s.ciRepo.GetCollectionItemsByIdsMap(collectionItemsIDs, authUser)
		if err != nil {
			return nil, err
		}
		for _, collectionItem := range collectionItems {
			userLogins = append(userLogins, collectionItem.Owner.Login)
			collectionItemsMap[collectionItem.ID.String()] = models.CollectionItems{
				ID:              collectionItem.ID,
				Name:            collectionItem.Name,
				Transliteration: collectionItem.CollectionTransliteration,
				UserLogin:       collectionItem.Owner.Login,
				Description:     collectionItem.CollectionName,
			}
		}
	}

	if len(postsIDs) > 0 {
		posts, err := s.postsClient.GetPostsByIds(ctx, postsIDs, authUser, false)
		if err != nil {
			return nil, err
		}
		for _, post := range posts.GetData() {
			userLogins = append(userLogins, post.GetAuthor())
			stringID := strconv.Itoa(int(post.GetId()))
			postsMap[stringID] = models.Posts{
				ID:       post.GetId(),
				Title:    post.GetTitle(),
				Translit: post.GetTranslit(),
				Author: models.SubUsers{
					Login: post.GetAuthor(),
				},
			}
		}
	}

	users, err := s.userRepo.GetUsersByLogins(userLogins)
	if err != nil {
		return nil, err
	}

	for _, n := range notifications.GetData() {
		initUser, _ := s.userRepo.GetUserByLogin(n.SenderLogin)
		tempItem := dto.Notifications{
			ID:          n.Id,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
			UserLogin:   n.UserLogin,
			Type:        n.Type,
			SenderLogin: n.SenderLogin,
			TargetID:    n.TargetId,
			IsRead:      n.IsRead,
			Date:        n.Date,
			Action:      n.Action,
			TargetUser: dto.User{
				Login:       targetUser.Login,
				IsPremium:   targetUser.IsPremium,
				AvatarURL:   targetUser.AvatarURL,
				ProfileName: targetUser.ProfileName,
			},
			SenderUser: dto.User{
				Login:       initUser.Login,
				IsPremium:   initUser.IsPremium,
				AvatarURL:   initUser.AvatarURL,
				ProfileName: initUser.ProfileName,
			},
		}
		switch n.TargetType {
		case "post":
			tId, _ := strconv.ParseUint(n.TargetId, 10, 64)
			post := postsMap[strconv.FormatUint(tId, 10)]
			owner := users[post.Author.Login]
			tempItem.Target = dto.TargetItem{
				ID:              strconv.FormatUint(post.ID, 10),
				Name:            post.Title,
				Transliteration: post.Translit,
				TargetType:      "post",
			}
			tempItem.Owner = dto.User{
				Login:       owner.Login,
				IsPremium:   owner.IsPremium,
				AvatarURL:   owner.AvatarURL,
				ProfileName: owner.ProfileName,
			}

		case "collection":
			owner := users[collectionsMap[n.TargetId].UserLogin]
			tempItem.Target = dto.TargetItem{
				ID:              collectionsMap[n.TargetId].ID.String(),
				Name:            collectionsMap[n.TargetId].Name,
				Transliteration: collectionsMap[n.TargetId].Transliteration,
				TargetType:      "collection",
			}
			tempItem.Owner = dto.User{
				Login:       owner.Login,
				IsPremium:   owner.IsPremium,
				AvatarURL:   owner.AvatarURL,
				ProfileName: owner.ProfileName,
			}
		case "collectionItem":
			owner := users[collectionItemsMap[n.TargetId].UserLogin]
			tempItem.Target = dto.TargetItem{
				ID:               collectionItemsMap[n.TargetId].ID.String(),
				Name:             collectionItemsMap[n.TargetId].Name,
				Transliteration:  collectionItemsMap[n.TargetId].Transliteration,
				TargetType:       "collectionItem",
				TargetParentName: collectionItemsMap[n.TargetId].Description,
			}
			tempItem.Owner = dto.User{
				Login:       owner.Login,
				IsPremium:   owner.IsPremium,
				AvatarURL:   owner.AvatarURL,
				ProfileName: owner.ProfileName,
			}

		case "photo":
			owner := users[photosMap[n.TargetId].Author.Login]
			dbCollection, _ := s.collectionRepo.GetCollectionByIdWithoutLimits(photosMap[n.TargetId].CollectionID)
			tempItem.Target = dto.TargetItem{
				ID:               strconv.FormatUint(photosMap[n.TargetId].ID, 10),
				Name:             photosMap[n.TargetId].Path,
				Transliteration:  dbCollection.Transliteration,
				TargetType:       "photo",
				TargetParentName: dbCollection.Name,
			}
			tempItem.Owner = dto.User{
				Login:       owner.Login,
				IsPremium:   owner.IsPremium,
				AvatarURL:   owner.AvatarURL,
				ProfileName: owner.ProfileName,
			}

		}

		resultNotifications = append(resultNotifications, tempItem)
	}

	return &dto.NotificationsDataResponse{Data: resultNotifications}, nil
}

func (s *NotificationsService) ReadNotification(ctx context.Context, ids []string) (*dto.CommonResponse, error) {
	_, err := s.notificationsClient.ReadNotifications(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *NotificationsService) DeleteNotification(ctx context.Context, targetId, senderLogin string) error {
	_, err := s.notificationsClient.DeleteNotification(ctx, targetId, senderLogin)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationsService) DeleteAllNotificationsByTargetID(ctx context.Context, targetId string) error {
	_, err := s.notificationsClient.DeleteAllNotificationsByTargetID(ctx, targetId)
	if err != nil {
		return err
	}
	return nil
}
