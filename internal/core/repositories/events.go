package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
)

type EventsRepository struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) *EventsRepository {
	return &EventsRepository{db: db}
}

func (r *EventsRepository) AddEvent(
	initiatorLogin string,
	action string,
	targetType string,
	targetName string,
	targetUserLogin *string,
	targetCollectionID *uuid.UUID,
	targetItemID *uuid.UUID,
) error {
	event := &models.Events{
		Action:          action,
		EventTargetType: targetType,
		TargetName:      targetName,
		InitiatorLogin:  initiatorLogin,
		Date:            time.Now().String(),
	}

	switch targetType {
	case "user":
		event.TargetUserLogin = targetUserLogin
	case "collection":
		event.TargetCollectionID = targetCollectionID
	case "collectionItem":
		event.TargetItemID = targetItemID
	default:
		return fmt.Errorf("unknown target type: %s", targetType)
	}

	return r.db.Create(event).Error
}

func (r *EventsRepository) GetEvents(subscriptions []string) ([]models.Events, error) {
	var events []models.Events

	err := r.db.
		Preload("Initiator").
		Where("initiator_login IN ?", subscriptions).
		Order("date DESC").
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	for i := range events {
		switch events[i].EventTargetType {
		case "user":
			if events[i].TargetUserLogin != nil {
				var user models.Users
				if err := r.db.Where("login = ?", *events[i].TargetUserLogin).First(&user).Error; err == nil {
					events[i].TargetUser = &user
				}
				events[i].TargetCollection = nil
				events[i].TargetItem = nil
			}

		case "collection":
			if events[i].TargetCollectionID != nil {
				var collection models.Collections
				if err := r.db.Preload("User").First(&collection, *events[i].TargetCollectionID).Error; err == nil {
					events[i].TargetCollection = &collection
				}
				events[i].TargetUser = nil
				events[i].TargetItem = nil
			}

		case "item":
			if events[i].TargetItemID != nil {
				var item models.CollectionItems
				if err := r.db.Preload("Owner").First(&item, *events[i].TargetItemID).Error; err == nil {
					events[i].TargetItem = &item
				}
				events[i].TargetUser = nil
				events[i].TargetCollection = nil
			}
		}
	}

	return events, nil
}

func (r *EventsRepository) GetFilteredEvents(
	userLogin string,
	collectionId string,
	collectionItemId string,
) ([]models.Events, error) {
	var events []models.Events

	query := r.db.Model(&models.Events{}).
		Where("initiator_login = ?", userLogin).
		Or("target_collection_id = ?", collectionId).
		Or("target_item_id = ?", collectionItemId).
		Order("date DESC")

	query = query.
		Joins("LEFT JOIN users ON users.login = events.target_user_login").
		Joins("LEFT JOIN collections ON collections.id = events.target_collection_id").
		Joins("LEFT JOIN collection_items ON collection_items.id = events.target_collection_item_id").
		Preload("TargetUserLogin").
		Preload("TargetCollectionId").
		Preload("TargetCollectionItemId")

	if err := query.Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}
