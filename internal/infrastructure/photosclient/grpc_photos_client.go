package photosclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"lootor/gen/go/microservices"
)

type GRPCPhotosClient struct {
	client microservices.PhotosServiceClient
	conn   *grpc.ClientConn
}

func NewGRPPhotosClient(addr string) (*GRPCPhotosClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCPhotosClient{
		client: microservices.NewPhotosServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCPhotosClient) CreatePhotos(ctx context.Context, req *microservices.CreatePhotosRequest) (*microservices.PhotosResponse, error) {
	resp, err := c.client.CreatePhotos(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) FindByUser(ctx context.Context, req *microservices.FindByUserRequest) (*microservices.PhotosResponse, error) {
	resp, err := c.client.FindByUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) DeletePhoto(ctx context.Context, req *microservices.DeletePhotoRequest) (*microservices.DeletePhotoResponse, error) {
	resp, err := c.client.DeletePhoto(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) FindOneByID(ctx context.Context, req *microservices.FindOneByIdRequest) (*microservices.PhotoResponse, error) {
	resp, err := c.client.FindOneByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) FindAllByCollectionID(ctx context.Context, req *microservices.FindByCollectionIdRequest) (*microservices.PhotosResponse, error) {
	resp, err := c.client.FindAllByCollectionID(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) LikePhoto(ctx context.Context, req *microservices.LikePhotoRequest) (*microservices.DeletePhotoResponse, error) {
	resp, err := c.client.LikePhoto(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) GetPhotosByCollectionsIds(ctx context.Context, req *microservices.GetPhotosByCollectionsIdsMapRequest) (*microservices.GetPhotosByIdsMapResponse, error) {
	resp, err := c.client.GetPhotosByCollectionsIds(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) GetPhotosByIDs(ctx context.Context, req *microservices.GetByIdsRequest) (*microservices.PhotosResponse, error) {
	resp, err := c.client.GetPhotosByIDs(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) GetCountByCollection(ctx context.Context, req *microservices.CountRequestPhoto) (*microservices.CountResponsePhoto, error) {
	resp, err := c.client.GetCountByCollection(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) IncrementCommentsCount(ctx context.Context, id uint64) (*microservices.CommentsCountResponsePhoto, error) {
	return c.client.IncrementCommentsCount(ctx, &microservices.CommentsCountRequestPhoto{PhotoId: id})
}

func (c *GRPCPhotosClient) DecrementCommentsCount(ctx context.Context, id uint64) (*microservices.CommentsCountResponsePhoto, error) {
	return c.client.DecrementCommentsCount(ctx, &microservices.CommentsCountRequestPhoto{PhotoId: id})
}

func (c *GRPCPhotosClient) UpdatePhoto(ctx context.Context, req *microservices.UpdatePhotoRequest) (*microservices.PhotoResponse, error) {
	resp, err := c.client.UpdatePhoto(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GRPCPhotosClient) GetPhotosByIDsMap(ctx context.Context, req *microservices.GetByIdsRequest) (*microservices.GetPhotosByIdsMapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	photosResult, err := c.client.GetPhotosByIDsMap(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return photosResult, nil
}

func (c *GRPCPhotosClient) GetCountsByCollectionsIDsMap(ctx context.Context, req *microservices.GetCountsByCollectionsIdsMapRequest) (*microservices.GetCountsByCollectionsIdsMapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	countsResult, err := c.client.GetCountsByCollectionsIDsMap(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return countsResult, nil
}

func (c *GRPCPhotosClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
