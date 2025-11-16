package postsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/pkg/utils"
)

type GRPCPostsClient struct {
	client microservices.PostsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCPostsClient(addr string) (*GRPCPostsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCPostsClient{
		client: microservices.NewPostsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCPostsClient) CreatePost(ctx context.Context, dto *models.PostRequest) (*microservices.PostResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(dto.Content)

	resp, err := c.client.CreatePost(ctx, &microservices.CreatePostRequest{
		Content: normalizedContent,
		Date:    dto.Date,
		Author:  dto.Author,
		IsDraft: dto.IsDraft,
		Title:   dto.Title,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPostsClient) GetAllPosts(ctx context.Context, order, limit, offset string, authUserIsPremium bool, authUser string) (*microservices.GetAllPostsResponse, error) {
	return c.client.GetAllPosts(ctx, &microservices.GetAllPostsRequest{Limit: limit, Offset: offset, Order: order, IsPremium: authUserIsPremium, AuthUserLogin: authUser})
}

func (c *GRPCPostsClient) GetPostsByUser(ctx context.Context, limit, offset, login string, authUserIsPremium bool, authUser string, isDraft bool) (*microservices.GetAllPostsResponse, error) {
	return c.client.GetPostsByUser(ctx, &microservices.GetPostsByUserRequest{
		Login:         login,
		Offset:        offset,
		Limit:         limit,
		IsPremium:     authUserIsPremium,
		AuthUserLogin: authUser,
		IsDraft:       isDraft,
	})
}

func (c *GRPCPostsClient) DeletePost(ctx context.Context, id uint64) (*microservices.DeletePostResponse, error) {
	return c.client.DeletePost(ctx, &microservices.DeletePostRequest{
		Id: id,
	})
}

func (c *GRPCPostsClient) ReactPost(ctx context.Context, dto *models.ReactRequest) (*microservices.ReactResponse, error) {
	return c.client.IncrementReaction(ctx, &microservices.ReactRequest{
		UserLogin: dto.UserLogin,
		PostId:    dto.PostId,
		Reaction:  dto.Reaction,
	})
}

func (c *GRPCPostsClient) ReactPostDecrement(ctx context.Context, dto *models.ReactRequest) (*microservices.ReactResponse, error) {
	return c.client.DecrementReaction(ctx, &microservices.ReactDecrementRequest{
		UserLogin: dto.UserLogin,
		PostId:    dto.PostId,
		Reaction:  dto.Reaction,
	})
}

func (c *GRPCPostsClient) UpdatePost(ctx context.Context, dto *models.PostUpdateRequest) (*microservices.PostResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(dto.Content)

	resp, err := c.client.UpdatePost(ctx, &microservices.UpdatePostRequest{
		Id:       dto.Id,
		Content:  normalizedContent,
		IsDraft:  dto.IsDraft,
		Title:    dto.Title,
		Translit: dto.Translit,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPostsClient) GetPostById(ctx context.Context, id uint64, authUserIsPremium bool) (*microservices.PostResponse, error) {
	return c.client.GetPost(ctx, &microservices.PostRequest{
		Id:        id,
		IsPremium: authUserIsPremium,
	})
}

func (c *GRPCPostsClient) GetPostByTranslit(ctx context.Context, translit string, authUserIsPremium bool) (*microservices.PostResponse, error) {
	return c.client.GetPostByTranslit(ctx, &microservices.PostRequestByTranslit{
		Translit:  translit,
		IsPremium: authUserIsPremium,
	})
}

func (c *GRPCPostsClient) IncrementViews(ctx context.Context, req *models.IncrementRequest) (*microservices.ViewsResponse, error) {
	return c.client.IncrementViews(ctx, &microservices.ViewsRequest{
		PostIds: req.PostIds,
	})
}

func (c *GRPCPostsClient) IncrementCommentsCount(ctx context.Context, id uint64) (*microservices.CommentsCountResponse, error) {
	return c.client.IncrementCommentsCount(ctx, &microservices.CommentsCountRequest{PostId: id})
}

func (c *GRPCPostsClient) DecrementCommentsCount(ctx context.Context, id uint64) (*microservices.CommentsCountResponse, error) {
	return c.client.DecrementCommentsCount(ctx, &microservices.CommentsCountRequest{PostId: id})
}

func (c *GRPCPostsClient) GetCount(ctx context.Context, userLogin string) (*microservices.CountResponse, error) {
	return c.client.GetCountByUser(ctx, &microservices.CountRequest{UserLogin: userLogin})
}

func (c *GRPCPostsClient) GetPostsByIds(ctx context.Context, ids []string, authUserLogin string, authUserIsPremium bool) (*microservices.GetPostsByIdsMapResponse, error) {
	return c.client.GetPostsByIds(ctx, &microservices.GetPostsByIdsMapRequest{PostIds: ids, AuthUserLogin: authUserLogin, IsPremium: authUserIsPremium})
}

func (c *GRPCPostsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
