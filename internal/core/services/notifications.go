package services

import (
	"context"
	"lootor/gen/go/microservices"
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
		AvatarUrl:   targetUser.AvatarUrl,
		ProfileName: targetUser.ProfileName,
	}

	senderUserReq := dto.User{
		Login:       senderUser.Login,
		IsPremium:   senderUser.IsPremium,
		AvatarUrl:   senderUser.AvatarUrl,
		ProfileName: senderUser.ProfileName,
	}

	ownerUserReq := dto.User{
		Login:       owner.Login,
		IsPremium:   owner.IsPremium,
		AvatarUrl:   owner.AvatarUrl,
		ProfileName: owner.ProfileName,
	}

	request.TargetUser = targetUserReq
	request.SenderUser = senderUserReq
	request.TargetItem = *targetReq
	request.Owner = ownerUserReq

	finalRequest := dto.NotificationsRequest{
		Login:       request.Login,
		TargetId:    request.TargetId,
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

func (s *NotificationsService) GetAllNotifications(ctx context.Context, login, limit, offset string) (*dto.NotificationsDataResponse, error) {
	notifications, err := s.notificationsClient.GetAllNotifications(ctx, login, limit, offset)
	if err != nil {
		return nil, err
	}
	targetUser, _ := s.userRepo.GetUserByLogin(login)
	var resultNotifications []dto.Notifications
	for _, n := range notifications.Data {
		initUser, _ := s.userRepo.GetUserByLogin(n.SenderLogin)
		tempItem := dto.Notifications{
			Id:          n.Id,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
			UserLogin:   n.UserLogin,
			Type:        n.Type,
			SenderLogin: n.SenderLogin,
			TargetId:    n.TargetId,
			IsRead:      n.IsRead,
			Date:        n.Date,
			Action:      n.Action,
			TargetUser: dto.User{
				Login:       targetUser.Login,
				IsPremium:   targetUser.IsPremium,
				AvatarUrl:   targetUser.AvatarUrl,
				ProfileName: targetUser.ProfileName,
			},
			SenderUser: dto.User{
				Login:       initUser.Login,
				IsPremium:   initUser.IsPremium,
				AvatarUrl:   initUser.AvatarUrl,
				ProfileName: initUser.ProfileName,
			},
		}
		switch n.TargetType {
		case "post":
			tId, _ := strconv.ParseUint(n.TargetId, 10, 64)
			post, err := s.postsClient.GetPostById(ctx, tId, false)
			if err == nil {
				owner, _ := s.userRepo.GetUserByLogin(post.Data.GetAuthor())
				tempItem.Target = dto.TargetItem{
					Id:              strconv.FormatUint(post.Data.Id, 10),
					Name:            post.Data.Title,
					Transliteration: post.Data.Translit,
					TargetType:      "post",
				}
				tempItem.Owner = dto.User{
					Login:       owner.Login,
					IsPremium:   owner.IsPremium,
					AvatarUrl:   owner.AvatarUrl,
					ProfileName: owner.ProfileName,
				}
			}
		case "collection":
			collection, err := s.collectionRepo.GetByIdWithoutCollectionItems(n.TargetId)
			if err == nil {
				owner, _ := s.userRepo.GetUserByLogin(collection.UserLogin)
				tempItem.Target = dto.TargetItem{
					Id:              collection.Id.String(),
					Name:            collection.Name,
					Transliteration: collection.Transliteration,
					TargetType:      "collection",
				}
				tempItem.Owner = dto.User{
					Login:       owner.Login,
					IsPremium:   owner.IsPremium,
					AvatarUrl:   owner.AvatarUrl,
					ProfileName: owner.ProfileName,
				}
			}
		case "collectionItem":
			ci, err := s.ciRepo.GetCIByID(n.TargetId)
			if err == nil {
				owner, _ := s.userRepo.GetUserByLogin(ci.UserLogin)
				tempItem.Target = dto.TargetItem{
					Id:              ci.Id.String(),
					Name:            ci.Name,
					Transliteration: ci.Collections[0].Transliteration,
					TargetType:      "collectionItem",
				}
				tempItem.Owner = dto.User{
					Login:       owner.Login,
					IsPremium:   owner.IsPremium,
					AvatarUrl:   owner.AvatarUrl,
					ProfileName: owner.ProfileName,
				}
			}
		case "photo":
			tId, _ := strconv.ParseUint(n.TargetId, 10, 64)
			post, err := s.photosClient.FindOneById(ctx, &microservices.FindOneByIdRequest{
				Id:            tId,
				AuthUserLogin: login,
			})
			if err == nil {
				owner, _ := s.userRepo.GetUserByLogin(post.Data.GetAuthor())
				tempItem.Target = dto.TargetItem{
					Id:              strconv.FormatUint(post.Data.Id, 10),
					Name:            post.Data.Path,
					Transliteration: post.Data.Path,
					TargetType:      "photo",
				}
				tempItem.Owner = dto.User{
					Login:       owner.Login,
					IsPremium:   owner.IsPremium,
					AvatarUrl:   owner.AvatarUrl,
					ProfileName: owner.ProfileName,
				}
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

func (s *NotificationsService) DeleteAllNotificationsByTargetId(ctx context.Context, targetId string) error {
	_, err := s.notificationsClient.DeleteAllNotificationsByTargetId(ctx, targetId)
	if err != nil {
		return err
	}
	return nil
}
