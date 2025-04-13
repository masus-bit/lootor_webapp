package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	models2 "lootor/internal/core/models"
	"time"
)

type EventsRepository struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) *EventsRepository {
	return &EventsRepository{db: db}
}

func (r *EventsRepository) AddEvent(
	user string,
	action string,
	eventTarget string,
	targetName string,
	targetUser *string,
	targetCollection *uuid.UUID,
	targetCollectionItem *uuid.UUID,
) error {
	event := &models2.Events{
		Action:               action,
		EventTarget:          eventTarget,
		TargetName:           targetName,
		Date:                 time.Now().Format(time.RFC3339),
		User:                 models2.Users{Login: user},
		TargetUser:           models2.Users{Login: *targetUser},
		TargetCollection:     models2.Collections{Id: *targetCollection},
		TargetCollectionItem: models2.CollectionItems{Id: *targetCollectionItem},
	}

	if err := r.db.Create(event).Error; err != nil {
		return err
	}

	return nil
}

func (r *EventsRepository) GetEvents(subscriptions []string) ([]models2.Events, error) {
	var events []models2.Events

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
