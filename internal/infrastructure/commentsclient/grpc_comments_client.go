package commentsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"strconv"
)

type GRPCCommentsClient struct {
	client microservices.CommentsServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCCommentsClient(addr string) (*GRPCCommentsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCCommentsClient{
		client: microservices.NewCommentsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCCommentsClient) CreateComment(ctx context.Context, dto *dto.CommentsRequest) (*microservices.CommentResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(dto.Content)

	resp, err := c.client.CreateComment(ctx, &microservices.CreateCommentRequest{
		Content:    normalizedContent,
		Date:       dto.Date,
		Author:     dto.Author,
		ParentID:   *dto.ParentID,
		TargetID:   dto.TargetID,
		EntityType: dto.EntityType,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCCommentsClient) GetAllComments(ctx context.Context, id, limit, offset, entityType string) (*microservices.GetAllCommentsResponse, error) {
	return c.client.GetAllComments(ctx, &microservices.GetAllCommentsRequest{Limit: limit, Offset: offset, TargetID: id, EntityType: entityType})
}

func (c *GRPCCommentsClient) DeleteComments(ctx context.Context, id string) (*microservices.DeleteCommentResponse, error) {
	return c.client.DeleteComment(ctx, &microservices.DeleteCommentRequest{
		ID: id,
	})
}

func (c *GRPCCommentsClient) LoadAnswers(ctx context.Context, dto *dto.AnswersRequest) (*microservices.AnswersResponse, error) {
	return c.client.LoadAnswers(ctx, &microservices.AnswersRequest{
		ID:     dto.ID,
		Limit:  strconv.Itoa(dto.Limit),
		Offset: strconv.Itoa(dto.Offset),
	})
}

func (c *GRPCCommentsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
