package eventsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
)

type GRPCEventsClient struct {
	client microservices.EventsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCClient(addr string) (*GRPCEventsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCEventsClient{
		client: microservices.NewEventsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCEventsClient) AddEvent(ctx context.Context, req *models.AddEventRequest) (*microservices.AddEventResponse, error) {
	return c.client.AddEvent(ctx, &microservices.AddEventRequest{
		Action:         req.Action,
		TargetType:     req.EventTargetType,
		TargetName:     req.TargetName,
		InitiatorLogin: req.InitiatorLogin,
		Params: &microservices.Params{
			TargetUserLogin:      req.Params.TargetUserLogin,
			TargetCollectionId:   req.Params.TargetCollectionID,
			TargetItemId:         req.Params.TargetItemID,
			TargetWishListItemId: req.Params.TargetWLID,
			TargetPostID:         req.Params.TargetPostID,
			TargetTagId:          req.Params.TargetTagID,
			TargetPhotoId:        req.Params.TargetPhotoID,
		},
		TagRelatedEntityType: req.Params.TagRelatedEntityType,
	})
}

func (c *GRPCEventsClient) GetEvents(ctx context.Context, req *models.GetEventsRequest) (*microservices.GetEventsResponse, error) {
	return c.client.GetEvents(ctx, &microservices.GetEventsRequest{
		Limit:            req.Limit,
		Offset:           req.Offset,
		Subscriptions:    req.Subscriptions,
		Actions:          req.Actions,
		EventTargetTypes: req.EventTargetTypes,
	})
}

func (c *GRPCEventsClient) GetFilteredEvents(ctx context.Context, req *models.GetFilteredEventsRequest) (*microservices.GetEventsResponse, error) {
	return c.client.GetFilteredEvents(ctx, &microservices.GetFilteredEventsRequest{
		Limit:             req.Limit,
		Offset:            req.Offset,
		UserLogin:         req.UserLogin,
		CollectionID:      req.CollectionID,
		CollectionItemID:  req.CollectionItemID,
		WishListItemID:    req.WishListItemID,
		CollectionItemIDs: req.CollectionItemIDs,
	})
}

func (c *GRPCEventsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
