package tagsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
)

type GRPCTagsClient struct {
	client microservices.TagsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCTagsClient(addr string) (*GRPCTagsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCTagsClient{
		client: microservices.NewTagsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCTagsClient) CreateTag(ctx context.Context, dto *models.TagCreateRequest) (*microservices.TagsCreateResponse, error) {

	resp, err := c.client.CreateTags(ctx, &microservices.CreateTagsRequest{
		Names:  dto.Names,
		Author: dto.Author,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) SearchTags(ctx context.Context, dto *models.TagsSearchRequest) (*microservices.TagsDataResponse, error) {

	resp, err := c.client.SearchTags(ctx, &microservices.TagsSearchRequest{
		Name: dto.Name,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) MergeTags(ctx context.Context, dto *microservices.MergeTagsRequest) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.MergeTags(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) FindAllEntitiesByTag(ctx context.Context, dto *microservices.GetEntitiesByTagRequest) (*microservices.TagDataResponse, error) {
	resp, err := c.client.FindAllEntitiesByTag(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) AddTagToEntity(ctx context.Context, dto *microservices.AddTagsToEntityRequest) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.AddTagToEntity(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) RemoveTagsFromEntity(ctx context.Context, dto *microservices.RemoveTagsRequest) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.RemoveTagsFromEntity(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) UpdateTag(ctx context.Context, dto *microservices.UpdateTagRequest) (*microservices.TagItem, error) {
	resp, err := c.client.UpdateTag(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetAllTags(ctx context.Context) (*microservices.GetShortsResponse, error) {
	resp, err := c.client.GetAllTags(ctx, &microservices.GetAllTagsRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsByEntityId(ctx context.Context, dto *microservices.GetTagsByEntityIdRequest) (*microservices.GetShortsResponse, error) {
	resp, err := c.client.GetTagsByEntityId(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) AddTagsToEntity(ctx context.Context, dto *microservices.AddFewTagsToEntityRequest) (*microservices.GetShortsResponse, error) {
	resp, err := c.client.AddTagsToEntity(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsByEntityIdsMap(ctx context.Context, dto *microservices.GetTagsByEntityIdsMapRequest) (*microservices.TagsMapResponse, error) {
	resp, err := c.client.GetTagsByEntityIdsMap(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) UpdateTagsOfEntity(ctx context.Context, dto *microservices.UpdateTagsOfEntityRequest) (*microservices.TagsCreateResponse, error) {
	resp, err := c.client.UpdateTagsOfEntity(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagBySlug(ctx context.Context, slug string) (*microservices.TagItem, error) {
	resp, err := c.client.GetTagBySlug(ctx, &microservices.GetTagBySlugRequest{
		Slug: slug,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
