package services

import (
	"context"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"time"
)

type SubscriptionService struct {
	repo *repositories.SubscriptionRepository
}

func NewSubscriptionService(repo *repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, login string, subType string, duration time.Duration) (*models.Subscription, error) {
	now := time.Now()
	subscription := &models.Subscription{
		UserLogin: login,
		Type:      subType,
		StartDate: now,
		EndDate:   now.Add(duration),
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) CheckActiveSubscription(ctx context.Context, login string) (bool, *models.Subscription, error) {
	sub, err := s.repo.GetActiveByUser(ctx, login)
	if err != nil {
		return false, nil, err
	}

	return sub != nil, sub, nil
}
