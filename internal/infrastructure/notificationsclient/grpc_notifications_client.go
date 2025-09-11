package notificationsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/pkg/dto"
)

type GRPCNotificationsClient struct {
	client microservices.NotificationsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCNotificationsClient(addr string) (*GRPCNotificationsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCNotificationsClient{
		client: microservices.NewNotificationsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCNotificationsClient) AddNotification(ctx context.Context, dto *dto.NotificationsRequest) (*microservices.NotificationResponse, error) {

	resp, err := c.client.CreateNotification(ctx, &microservices.CreateNotificationRequest{
		Date:        dto.Date,
		TargetId:    dto.TargetId,
		UserLogin:   dto.Login,
		Type:        dto.Type,
		SenderLogin: dto.SenderLogin,
		Action:      dto.Action,
		TargetType:  dto.TargetItem.TargetType,
		Target: &microservices.TargetItem{
			Id:              &dto.TargetItem.Id,
			Name:            &dto.TargetItem.Name,
			Transliteration: &dto.TargetItem.Transliteration,
			TargetType:      &dto.TargetItem.TargetType,
		},
		TargetUser: &microservices.User{
			Login:       &dto.TargetUser.Login,
			IsPremium:   &dto.TargetUser.IsPremium,
			AvatarUrl:   &dto.TargetUser.AvatarUrl,
			ProfileName: &dto.TargetUser.ProfileName,
		},
		SenderUser: &microservices.User{
			Login:       &dto.SenderUser.Login,
			IsPremium:   &dto.SenderUser.IsPremium,
			AvatarUrl:   &dto.SenderUser.AvatarUrl,
			ProfileName: &dto.SenderUser.ProfileName,
		},
		Owner: &microservices.User{
			Login:       &dto.Owner.Login,
			IsPremium:   &dto.Owner.IsPremium,
			AvatarUrl:   &dto.Owner.AvatarUrl,
			ProfileName: &dto.Owner.ProfileName,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCNotificationsClient) GetAllNotifications(ctx context.Context, login, limit, offset string) (*microservices.GetAllNotificationsResponse, error) {
	return c.client.GetAllNotifications(ctx, &microservices.GetAllNotificationsRequest{Limit: limit, Offset: offset, UserLogin: login})
}

func (c *GRPCNotificationsClient) ReadNotifications(ctx context.Context, ids []string) (*microservices.ReadResponse, error) {
	return c.client.ReadNotifications(ctx, &microservices.ReadRequest{
		Ids:    ids,
		IsRead: true,
	})
}

func (c *GRPCNotificationsClient) DeleteNotification(ctx context.Context, targetId, senderLogin string) (*microservices.DeleteNotificationResponse, error) {
	notification, err := c.client.RemoveNotification(ctx, &microservices.DeleteNotificationRequest{
		TargetId:  targetId,
		UserLogin: senderLogin,
	})
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (c *GRPCNotificationsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
