package commentsclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lootor/gen/go/microservices"
)

type GRPCLikesClient struct {
	client microservices.LikesServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCLikesClient(addr string) (*GRPCLikesClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCLikesClient{
		client: microservices.NewLikesServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCLikesClient) Like(
	ctx context.Context,
	commentID, author string,
	isLiked bool,
) (*microservices.LikeResponse, error) {

	resp, err := c.client.Like(
		ctx, &microservices.LikeRequest{
			CommentId: commentID,
			Author:    author,
			IsLiked:   isLiked,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCLikesClient) Dislike(
	ctx context.Context,
	commentID, author string,
	isLiked bool,
) (*microservices.LikeResponse, error) {

	resp, err := c.client.Dislike(
		ctx, &microservices.DislikeRequest{
			CommentId:  commentID,
			Author:     author,
			IsDisliked: isLiked,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCLikesClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
