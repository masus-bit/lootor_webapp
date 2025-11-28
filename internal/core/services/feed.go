package services

import (
	"context"
	"github.com/google/uuid"
	"lootor/internal/infrastructure/newsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"strconv"
	"time"
)

type FeedService struct {
	feedClient *newsclient.GRPCClient
}

func NewFeedService(feedClient *newsclient.GRPCClient) *FeedService {
	return &FeedService{feedClient: feedClient}
}

func (s *FeedService) CreateFeed(ctx context.Context, request *dto.FeedRequest) (*dto.FeedResponse, error) {
	feed, err := s.feedClient.CreateNews(ctx, request.Type, request.Date, request.Content)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(feed.Data.Id)
	if err != nil {
		return nil, err
	}
	content := utils.NormalizeContent(feed.Data.Content)

	createdAtAsTime, _ := time.Parse(time.RFC3339, feed.Data.CreatedAt)
	updatedAtAsTime, _ := time.Parse(time.RFC3339, feed.Data.UpdatedAt)

	resultFeed := dto.Feed{
		ID:        id,
		Type:      feed.Data.Type,
		Date:      feed.Data.Date,
		Content:   content,
		CreatedAt: createdAtAsTime,
		UpdatedAt: updatedAtAsTime,
	}

	return &dto.FeedResponse{Data: resultFeed}, nil
}

func (s *FeedService) GetFeed(ctx context.Context, id string) (*dto.FeedResponse, error) {
	feed, err := s.feedClient.GetNews(ctx, id)
	if err != nil {
		return nil, err
	}

	parsedId, err := uuid.Parse(feed.Data.Id)
	if err != nil {
		return nil, err
	}

	content := utils.NormalizeContent(feed.Data.Content)

	createdAtAsTime, _ := time.Parse(time.RFC3339, feed.Data.CreatedAt)
	updatedAtAsTime, _ := time.Parse(time.RFC3339, feed.Data.UpdatedAt)

	resultFeed := dto.Feed{
		ID:        parsedId,
		Type:      feed.Data.Type,
		Date:      feed.Data.Date,
		Content:   content,
		CreatedAt: createdAtAsTime,
		UpdatedAt: updatedAtAsTime,
	}

	return &dto.FeedResponse{Data: resultFeed}, nil
}

func (s *FeedService) GetAllFeed(ctx context.Context, limit, offset string) (*dto.FeedDataResponse, error) {
	feed, err := s.feedClient.GetAllNews(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var resultFeed []dto.Feed
	for _, f := range feed.Data {
		parsedId, err := uuid.Parse(f.Id)
		if err != nil {
			return nil, err
		}
		content := utils.NormalizeContent(f.Content)

		createdAtAsTime, _ := time.Parse(time.RFC3339, f.CreatedAt)
		updatedAtAsTime, _ := time.Parse(time.RFC3339, f.UpdatedAt)

		resultFeed = append(resultFeed, dto.Feed{
			ID:        parsedId,
			Type:      f.Type,
			Date:      f.Date,
			Content:   content,
			CreatedAt: createdAtAsTime,
			UpdatedAt: updatedAtAsTime,
		})
	}

	totalInt, _ := strconv.Atoi(feed.Total)

	return &dto.FeedDataResponse{Data: resultFeed, Total: int64(totalInt)}, nil
}

func (s *FeedService) DeleteFeed(ctx context.Context, id string) (*dto.CommonResponse, error) {
	_, err := s.feedClient.DeleteNews(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *FeedService) UpdateFeed(ctx context.Context, id string, request *dto.FeedRequest) (*dto.FeedResponse, error) {
	feed, err := s.feedClient.UpdateNews(ctx, id, request.Type, request.Date, request.Content)
	if err != nil {
		return nil, err
	}

	parsedId, err := uuid.Parse(feed.Data.Id)
	if err != nil {
		return nil, err
	}
	content := utils.NormalizeContent(feed.Data.Content)

	createdAtAsTime, _ := time.Parse(time.RFC3339, feed.Data.CreatedAt)
	updatedAtAsTime, _ := time.Parse(time.RFC3339, feed.Data.UpdatedAt)

	resultFeed := dto.Feed{
		ID:        parsedId,
		Type:      feed.Data.Type,
		Date:      feed.Data.Date,
		Content:   content,
		CreatedAt: createdAtAsTime,
		UpdatedAt: updatedAtAsTime,
	}

	return &dto.FeedResponse{Data: resultFeed}, nil
}
