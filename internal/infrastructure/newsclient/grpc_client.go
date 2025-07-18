package newsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/datatypes"
	"lootor/gen/go/feed"
	"lootor/internal/pkg/utils"
)

type GRPCClient struct {
	client feed.NewsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		client: feed.NewNewsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCClient) CreateNews(ctx context.Context, feedType, date string, content datatypes.JSON) (*feed.NewsResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(content)

	resp, err := c.client.CreateNews(ctx, &feed.CreateNewsRequest{
		Content: normalizedContent,
		Date:    date,
		Type:    feedType,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCClient) GetNews(ctx context.Context, id string) (*feed.GetNewsByIdResponse, error) {
	return c.client.GetNewsById(ctx, &feed.GetNewsByIdRequest{
		Id: id,
	})
}

func (c *GRPCClient) GetAllNews(ctx context.Context, limit, offset string) (*feed.GetAllNewsResponse, error) {
	return c.client.GetAllNews(ctx, &feed.GetAllNewsRequest{Limit: limit, Offset: offset})
}

func (c *GRPCClient) DeleteNews(ctx context.Context, id string) (*feed.DeleteNewsResponse, error) {
	return c.client.DeleteNews(ctx, &feed.DeleteNewsRequest{
		Id: id,
	})
}

func (c *GRPCClient) UpdateNews(ctx context.Context, id, feedType, date string, content datatypes.JSON) (*feed.UpdateNewsResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(content)
	return c.client.UpdateNews(ctx, &feed.UpdateNewsRequest{
		Content: normalizedContent,
		Date:    date,
		Id:      id,
		Type:    feedType,
	})
}

func (c *GRPCClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
