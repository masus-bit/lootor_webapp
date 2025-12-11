package commentsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
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

func (c *GRPCCommentsClient) CreateComment(ctx context.Context, req *dto.CommentsRequest) (*microservices.CommentResponse, error) {
	normalizedContent, _ := utils.GormJSONToProtoStruct(req.Content)

	resp, err := c.client.CreateComment(ctx, &microservices.CreateCommentRequest{
		Content:    normalizedContent,
		Date:       req.Date,
		Author:     req.Author,
		ParentId:   *req.ParentID,
		TargetId:   req.TargetID,
		EntityType: req.EntityType,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCCommentsClient) GetAllComments(ctx context.Context, id, limit, offset, entityType string) (*microservices.GetAllCommentsResponse, error) {
	return c.client.GetAllComments(ctx, &microservices.GetAllCommentsRequest{Limit: limit, Offset: offset, TargetId: id, EntityType: entityType})
}

func (c *GRPCCommentsClient) DeleteComments(ctx context.Context, id string) (*microservices.DeleteCommentResponse, error) {
	return c.client.DeleteComment(ctx, &microservices.DeleteCommentRequest{
		Id: id,
	})
}

func (c *GRPCCommentsClient) LoadAnswers(ctx context.Context, req *dto.AnswersRequest) (*microservices.AnswersResponse, error) {
	return c.client.LoadAnswers(ctx, &microservices.AnswersRequest{
		Id:     req.ID,
		Limit:  strconv.Itoa(req.Limit),
		Offset: strconv.Itoa(req.Offset),
	})
}

func (c *GRPCCommentsClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
