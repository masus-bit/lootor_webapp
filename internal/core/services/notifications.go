package services

import (
	"context"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/notificationsclient"
	"lootor/internal/pkg/dto"
)

type NotificationsService struct {
	notificationsClient *notificationsclient.GRPCNotificationsClient
	userRepo            *repositories.UsersRepository
	collectionRepo      *repositories.CollectionsRepository
	ciRepo              *repositories.CiRepository
}

func NewNotificationsService(notificationsClient *notificationsclient.GRPCNotificationsClient, userRepo *repositories.UsersRepository, collectionRepo *repositories.CollectionsRepository, ciRepo *repositories.CiRepository) *NotificationsService {
	return &NotificationsService{
		notificationsClient: notificationsClient, userRepo: userRepo, collectionRepo: collectionRepo, ciRepo: ciRepo}
}

func (s *NotificationsService) SendNotification(ctx context.Context, request *dto.NotificationsRequest) error {
	targetUser, _ := s.userRepo.GetUserByLogin(request.Login)
	senderUser, _ := s.userRepo.GetUserByLogin(request.SenderLogin)
	var targetReq dto.TargetItem
	collection, err := s.collectionRepo.GetByIdWithoutCollectionItems(request.TargetId)
	if err == nil {
		targetReq = dto.TargetItem{
			Id:              collection.Id.String(),
			Name:            collection.Name,
			Transliteration: collection.Transliteration,
			TargetType:      "collection",
		}
	}
	ci, err := s.ciRepo.GetCIByID(request.TargetId)
	if err == nil {
		targetReq = dto.TargetItem{
			Id:              ci.Id.String(),
			Name:            ci.Name,
			Transliteration: ci.Collections[0].Transliteration,
			TargetType:      "collectionItem",
		}
	}
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

	request.TargetUser = targetUserReq
	request.SenderUser = senderUserReq
	request.TargetItem = targetReq

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
	}

	_, err = s.notificationsClient.AddNotification(ctx, &finalRequest)
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
		collection, err := s.collectionRepo.GetByIdWithoutCollectionItems(n.TargetId)
		if err == nil {
			tempItem.Target = dto.TargetItem{
				Id:              collection.Id.String(),
				Name:            collection.Name,
				Transliteration: collection.Transliteration,
				TargetType:      "collection",
			}
		}
		ci, err := s.ciRepo.GetCIByID(n.TargetId)
		if err == nil {
			tempItem.Target = dto.TargetItem{
				Id:              ci.Id.String(),
				Name:            ci.Name,
				Transliteration: ci.Transliteration,
				TargetType:      "collectionItem",
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
