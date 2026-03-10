package tagsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
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

func (c *GRPCTagsClient) CreateTag(ctx context.Context, req *dto.TagCreateRequest) (
	*microservices.TagsCreateResponse,
	error,
) {

	resp, err := c.client.CreateTags(
		ctx, &microservices.CreateTagsRequest{
			Names:  req.Names,
			Author: req.Author,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) SearchTags(ctx context.Context, req *dto.TagsSearchRequest) (
	*microservices.TagsDataResponse,
	error,
) {

	resp, err := c.client.SearchTags(
		ctx, &microservices.TagsSearchRequest{
			Name: req.Name,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) MergeTags(
	ctx context.Context,
	req *microservices.MergeTagsRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.MergeTags(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) MergeSeries(
	ctx context.Context,
	req *microservices.MergeTagsRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.MergeSeries(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) FindAllEntitiesByTag(
	ctx context.Context,
	req *microservices.GetEntitiesByTagRequest,
) (*microservices.TagDataResponse, error) {
	resp, err := c.client.FindAllEntitiesByTag(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) AddTagToEntity(
	ctx context.Context,
	req *microservices.AddTagsToEntityRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.AddTagToEntity(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) RemoveTagsFromEntity(
	ctx context.Context,
	req *microservices.RemoveTagsRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.RemoveTagsFromEntity(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) UpdateTag(ctx context.Context, req *microservices.UpdateTagRequest) (
	*microservices.TagItem,
	error,
) {
	resp, err := c.client.UpdateTag(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetAllTags(
	ctx context.Context,
	req *microservices.GetAllTagsRequest,
) (*microservices.TagsDataResponse, error) {
	resp, err := c.client.GetAllTags(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsByEntityId(
	ctx context.Context,
	dto *microservices.GetTagsByEntityIdRequest,
) (*microservices.GetShortsResponse, error) {
	resp, err := c.client.GetTagsByEntityId(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) AddTagsToEntity(
	ctx context.Context,
	dto *microservices.AddFewTagsToEntityRequest,
) (*microservices.GetShortsResponse, error) {
	resp, err := c.client.AddTagsToEntity(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsByEntityIdsMap(
	ctx context.Context,
	dto *microservices.GetTagsByEntityIdsMapRequest,
) (*microservices.TagsMapResponse, error) {
	resp, err := c.client.GetTagsByEntityIdsMap(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) UpdateTagsOfEntity(
	ctx context.Context,
	dto *microservices.UpdateTagsOfEntityRequest,
) (*microservices.TagsCreateResponse, error) {
	resp, err := c.client.UpdateTagsOfEntity(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagBySlug(ctx context.Context, slug string) (*microservices.TagItem, error) {
	resp, err := c.client.GetTagBySlug(
		ctx, &microservices.GetTagBySlugRequest{
			Slug: slug,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTag(ctx context.Context, dto *microservices.GetTagRequest) (
	*microservices.TagItemShort,
	error,
) {
	resp, err := c.client.GetTag(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsByIDs(
	ctx context.Context,
	dto *microservices.GetTagsByIDsRequest,
) (*microservices.GetShortsResponse, error) {
	resp, err := c.client.GetTagsByIDs(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsByIdsMap(
	ctx context.Context,
	dto *microservices.GetTagsByIDsRequest,
) (*microservices.TagsMapResponseForEvents, error) {
	resp, err := c.client.GetTagsByIdsMap(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) FindAllTags(
	ctx context.Context,
	dto *microservices.FindAllTagsRequest,
) (*microservices.TagsDataResponse, error) {
	resp, err := c.client.FindAllTags(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) SearchGameTitles(
	ctx context.Context,
	dto *microservices.SearchGameTitlesRequest,
) (*microservices.GameSearchResponse, error) {
	resp, err := c.client.SearchGameTitles(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) RecordUserChoice(
	ctx context.Context,
	dto *microservices.RecordUserChoiceRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.RecordUserChoice(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetSearchSuggestions(
	ctx context.Context,
	dto *microservices.GetSearchSuggestionsRequest,
) (*microservices.SearchSuggestionsResponse, error) {
	resp, err := c.client.GetSearchSuggestions(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) UpdateVisibleLinks(
	ctx context.Context,
	dto *microservices.UpdateVisibleLinksRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.UpdateVisibleLinks(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) RemoveEntityTags(
	ctx context.Context,
	dto *microservices.RemoveEntityTagsRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.RemoveEntityTags(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) DeleteTags(
	ctx context.Context,
	dto *microservices.GetTagsByIDsRequest,
) (*microservices.TagsSuccessResponse, error) {
	resp, err := c.client.DeleteTags(ctx, dto)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCTagsClient) GetTagsTotalByUserLogin(
	ctx context.Context,
	userLogin string,
) (*microservices.TagsTotalByUserResponse, error) {
	return c.client.GetTagsTotalByUser(ctx, &microservices.GetTagsTotalByUserRequest{UserLogin: userLogin})
}

func (c *GRPCTagsClient) PublicUpdateTag(
	ctx context.Context,
	dto *microservices.PublicUpdateTagRequest,
) (*microservices.PublicUpdateTagResponse, error) {
	return c.client.PublicUpdateTag(ctx, dto)
}

func (c *GRPCTagsClient) MoveTagLinks(
	ctx context.Context,
	dto *microservices.MoveTagLinksRequest,
) (*microservices.MoveTagLinksResponse, error) {
	return c.client.MoveTagLinks(ctx, dto)
}

func (c *GRPCTagsClient) GetOnModerationTags(
	ctx context.Context,
	dto *microservices.OnModerationTagsRequest,
) (*microservices.OnModerationTagsResponse, error) {
	return c.client.GetOnModerationTags(ctx, dto)
}

func (c *GRPCTagsClient) TriggerMoveTags(
	ctx context.Context,
	dto *microservices.TriggerMoveTagsRequest,
) (*microservices.TagsSuccessResponse, error) {
	return c.client.TriggerMoveTags(ctx, dto)
}

func (c *GRPCTagsClient) TriggerReindex(
	ctx context.Context,
	dto *microservices.TriggerReindexRequest,
) (*microservices.ReindexStatusResponse, error) {
	return c.client.TriggerReindex(ctx, dto)
}

func (c *GRPCTagsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
