package newsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/datatypes"
	"lootor/gen/go/microservices"
	"lootor/internal/pkg/utils"
)

type GRPCClient struct {
	client microservices.NewsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		client: microservices.NewNewsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCClient) CreateNews(ctx context.Context, feedType, date string, content datatypes.JSON) (*microservices.NewsResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(content)

	resp, err := c.client.CreateNews(ctx, &microservices.CreateNewsRequest{
		Content: normalizedContent,
		Date:    date,
		Type:    feedType,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCClient) GetNews(ctx context.Context, id string) (*microservices.GetNewsByIdResponse, error) {
	return c.client.GetNewsById(ctx, &microservices.GetNewsByIdRequest{
		ID: id,
	})
}

func (c *GRPCClient) GetAllNews(ctx context.Context, limit, offset string) (*microservices.GetAllNewsResponse, error) {
	return c.client.GetAllNews(ctx, &microservices.GetAllNewsRequest{Limit: limit, Offset: offset})
}

func (c *GRPCClient) DeleteNews(ctx context.Context, id string) (*microservices.DeleteNewsResponse, error) {
	return c.client.DeleteNews(ctx, &microservices.DeleteNewsRequest{
		ID: id,
	})
}

func (c *GRPCClient) UpdateNews(ctx context.Context, id, feedType, date string, content datatypes.JSON) (*microservices.UpdateNewsResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(content)
	return c.client.UpdateNews(ctx, &microservices.UpdateNewsRequest{
		Content: normalizedContent,
		Date:    date,
		ID:      id,
		Type:    feedType,
	})
}

func (c *GRPCClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
