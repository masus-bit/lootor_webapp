package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *SubscriptionRepository) GetActiveByUser(ctx context.Context, login string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).
		Where("user_login = ? AND end_date > ?", login, time.Now()).
		Order("end_date desc").
		First(&subscription).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &subscription, err
}

func (r *SubscriptionRepository) DeactivateExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&models.Subscription{}).
		Where("end_date <= ?", time.Now()).
		Update("updated_at", time.Now()).Error
}
