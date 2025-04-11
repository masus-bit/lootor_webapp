package events

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lootor/internal/models"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddEvent(
	user string,
	action string,
	eventTarget string,
	targetName string,
	targetUser *string,
	targetCollection *string,
	targetCollectionItem *string,
) (*models.Events, error) {
	event := &models.Events{
		Action:               action,
		EventTarget:          eventTarget,
		TargetName:           targetName,
		Date:                 time.Now().Format(time.RFC3339),
		User:                 models.Users{Login: user},
		TargetUser:           models.Users{Login: *targetUser},
		TargetCollection:     models.Collections{Id: uuid.MustParse(*targetCollection)},
		TargetCollectionItem: models.CollectionItems{Id: uuid.MustParse(*targetCollectionItem)},
	}

	if err := r.db.Create(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (r *Repository) GetEvents(subscriptions []string) ([]models.Events, error) {
	var events []models.Events

	err := r.db.
		Preload("User").
		Preload("TargetUser").
		Preload("TargetCollection").
		Preload("TargetCollectionItem").
		Where("user_login IN ?", subscriptions).
		Order("date DESC").
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}
