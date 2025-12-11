package notificationsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
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

func (c *GRPCNotificationsClient) AddNotification(ctx context.Context, req *dto.NotificationsRequest) (*microservices.NotificationResponse, error) {

	resp, err := c.client.CreateNotification(ctx, &microservices.CreateNotificationRequest{
		Date:        req.Date,
		TargetId:    req.TargetID,
		UserLogin:   req.Login,
		Type:        req.Type,
		SenderLogin: req.SenderLogin,
		Action:      req.Action,
		TargetType:  req.TargetItem.TargetType,
		Target: &microservices.TargetItem{
			Id:              &req.TargetItem.ID,
			Name:            &req.TargetItem.Name,
			Transliteration: &req.TargetItem.Transliteration,
			TargetType:      &req.TargetItem.TargetType,
		},
		TargetUser: &microservices.User{
			Login:       &req.TargetUser.Login,
			IsPremium:   &req.TargetUser.IsPremium,
			AvatarUrl:   &req.TargetUser.AvatarURL,
			ProfileName: &req.TargetUser.ProfileName,
		},
		SenderUser: &microservices.User{
			Login:       &req.SenderUser.Login,
			IsPremium:   &req.SenderUser.IsPremium,
			AvatarUrl:   &req.SenderUser.AvatarURL,
			ProfileName: &req.SenderUser.ProfileName,
		},
		Owner: &microservices.User{
			Login:       &req.Owner.Login,
			IsPremium:   &req.Owner.IsPremium,
			AvatarUrl:   &req.Owner.AvatarURL,
			ProfileName: &req.Owner.ProfileName,
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

func (c *GRPCNotificationsClient) DeleteNotification(ctx context.Context, targetID, senderLogin string) (*microservices.DeleteNotificationResponse, error) {
	notification, err := c.client.RemoveNotification(ctx, &microservices.DeleteNotificationRequest{
		TargetId:  targetID,
		UserLogin: senderLogin,
	})
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (c *GRPCNotificationsClient) DeleteAllNotificationsByTargetID(ctx context.Context, targetID string) (*microservices.DeleteNotificationResponse, error) {
	return c.client.RemoveAllNotificationsByTargetId(ctx, &microservices.DeleteNotificationRequest{
		TargetId: targetID,
	})
}

func (c *GRPCNotificationsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
