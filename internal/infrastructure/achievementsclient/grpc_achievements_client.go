package achievementsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
)

type GRPCAchievementsClient struct {
	client microservices.AchievementsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCAchievementsClient(addr string) (*GRPCAchievementsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCAchievementsClient{
		client: microservices.NewAchievementsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCAchievementsClient) AddOrUpdateAchievement(ctx context.Context, req *microservices.AddOrUpdateAchievementRequest) (*microservices.AddOrUpdateAchievementResponse, error) {
	resp, err := c.client.AddOrUpdateAchievement(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCAchievementsClient) GetUserAchievements(ctx context.Context, req *microservices.GetUserAchievementsRequest) (*microservices.GetUserAchievementsResponse, error) {
	resp, err := c.client.GetUserAchievements(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *GRPCAchievementsClient) GetAchievedUserAchievements(ctx context.Context, req *microservices.GetUserAchievementsRequest) (*microservices.GetUserAchievementsResponse, error) {
	resp, err := c.client.GetAchievedUserAchievements(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCAchievementsClient) GetOneAchievement(ctx context.Context, req *microservices.GetOneAchievementRequest) (*microservices.GetOneAchievementResponse, error) {
	resp, err := c.client.GetOneAchievement(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCAchievementsClient) GetAllAchievementsItems(ctx context.Context, req *microservices.GetAllAchievementsItemsRequest) (*microservices.GetAllAchievementsItemsResponse, error) {
	resp, err := c.client.GetAllAchievementsItems(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCAchievementsClient) AddZeroAchievements(ctx context.Context, req *microservices.GetUserAchievementsRequest) (*microservices.GetUserAchievementsRequest, error) {
	resp, err := c.client.AddZeroAchievements(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCAchievementsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
