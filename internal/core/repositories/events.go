package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"lootor/internal/pkg/types"
	"strconv"
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
	targetWLID *uuid.UUID,
) error {
	event := &models.Events{
		Action:          action,
		EventTargetType: targetType,
		TargetName:      targetName,
		InitiatorLogin:  initiatorLogin,
		Date:            time.Now().String(),
	}
	fmt.Println(targetWLID)
	switch targetType {
	case "user":
		event.TargetUserLogin = targetUserLogin
	case "collection":
		event.TargetCollectionID = targetCollectionID
	case "collectionItem":
		event.TargetItemID = targetItemID
	case "wishListItem":
		event.TargetWishListItemID = targetWLID

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
				var user models.SubUsers
				if err := r.db.Where("login = ?", *events[i].TargetUserLogin).First(&user).Error; err == nil {
					events[i].TargetUser = &user
				}
				events[i].TargetCollection = nil
				events[i].TargetItem = nil
				events[i].TargetWishListItem = nil
			}

		case "collection":
			if events[i].TargetCollectionID != nil {
				var collection types.CommonShortType
				if err := r.db.Preload("User").First(&collection, *events[i].TargetCollectionID).Error; err == nil {
					events[i].TargetCollection = &collection
				}
				events[i].TargetUser = nil
				events[i].TargetItem = nil
				events[i].TargetWishListItem = nil
			}

		case "collectionItem":
			if events[i].TargetItemID != nil {
				var item types.CommonShortType
				if err := r.db.Preload("Owner").First(&item, *events[i].TargetItemID).Error; err == nil {
					events[i].TargetItem = &item
				}
				events[i].TargetUser = nil
				events[i].TargetCollection = nil
				events[i].TargetWishListItem = nil
			}

		case "wishListItem":
			if events[i].TargetWishListItemID != nil {
				var wlItem types.CommonShortType
				if err := r.db.Preload("User").First(&wlItem, *events[i].TargetWishListItemID).Error; err == nil {
					events[i].TargetWishListItem = &wlItem
				}
				events[i].TargetUser = nil
				events[i].TargetCollection = nil
				events[i].TargetItem = nil
			}

		}
	}

	return events, nil
}

func (r *EventsRepository) GetFilteredEvents(
	userLogin,
	collectionId,
	collectionItemId,
	wlItemId, limit, offset string,
) ([]models.Events, int64, error) {
	var events []models.Events

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var count int64

	query := r.db.Model(&models.Events{}).
		Where("initiator_login = ?", userLogin).Preload("Initiator")

	switch {
	case collectionId != "":
		query = query.Where("target_collection_id = ?", collectionId)
	case collectionItemId != "":
		query = query.Where("target_item_id = ?", collectionItemId)
	case wlItemId != "":
		query = query.Where("target_wish_list_item_id = ?", wlItemId)
	}

	query.Count(&count)

	if err := query.Limit(intLimit).Offset(intOffset).Find(&events).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}

	if err := r.loadEventRelations(&events); err != nil {
		return nil, 0, fmt.Errorf("failed to load event relations: %w", err)
	}

	return events, count, nil
}

func (r *EventsRepository) loadEventRelations(events *[]models.Events) error {
	// Сначала загружаем полные модели из БД
	var (
		userLogins      []string
		collectionIDs   []uuid.UUID
		itemIDs         []uuid.UUID
		wishListItemIDs []uuid.UUID
	)

	for _, event := range *events {
		if event.TargetUserLogin != nil {
			userLogins = append(userLogins, *event.TargetUserLogin)
		}
		if event.TargetCollectionID != nil {
			collectionIDs = append(collectionIDs, *event.TargetCollectionID)
		}
		if event.TargetItemID != nil {
			itemIDs = append(itemIDs, *event.TargetItemID)
		}
		if event.TargetWishListItemID != nil {
			wishListItemIDs = append(wishListItemIDs, *event.TargetWishListItemID)
		}
	}

	// 1. Загружаем полные модели из БД
	usersMap := make(map[string]models.Users)
	if len(userLogins) > 0 {
		var users []models.Users
		if err := r.db.Where("login IN ?", userLogins).Find(&users).Error; err != nil {
			return err
		}
		for _, user := range users {
			usersMap[user.Login] = user
		}
	}

	collectionsMap := make(map[uuid.UUID]models.Collections)
	if len(collectionIDs) > 0 {
		var collections []models.Collections
		if err := r.db.Where("id IN ?", collectionIDs).Find(&collections).Error; err != nil {
			return err
		}
		for _, collection := range collections {
			collectionsMap[collection.Id] = collection
		}
	}

	itemsMap := make(map[uuid.UUID]models.CollectionItems)
	if len(itemIDs) > 0 {
		var items []models.CollectionItems
		if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			itemsMap[item.Id] = item
		}
	}

	wishListItemsMap := make(map[uuid.UUID]models.WishListItems)
	if len(wishListItemIDs) > 0 {
		var wishListItems []models.WishListItems
		if err := r.db.Where("id IN ?", wishListItemIDs).Find(&wishListItems).Error; err != nil {
			return err
		}
		for _, item := range wishListItems {
			wishListItemsMap[item.Id] = item
		}
	}
	// 2. Конвертируем в response-модели
	for i, event := range *events {
		if event.TargetUserLogin != nil {
			if user, ok := usersMap[*event.TargetUserLogin]; ok {
				(*events)[i].TargetUser = convertToUserResponse(user)
			}
		}
		if event.TargetCollectionID != nil {
			if collection, ok := collectionsMap[*event.TargetCollectionID]; ok {
				(*events)[i].TargetCollection = convertToCollectionResponse(collection)
			}
		}
		if event.TargetItemID != nil {
			if item, ok := itemsMap[*event.TargetItemID]; ok {
				(*events)[i].TargetItem = convertToItemResponse(item)
			}
		}
		if event.TargetWishListItemID != nil {
			if item, ok := wishListItemsMap[*event.TargetWishListItemID]; ok {
				(*events)[i].TargetWishListItem = convertToWishListItemResponse(item)
			}
		}
	}

	return nil
}

func convertToUserResponse(user models.Users) *models.SubUsers {
	return &models.SubUsers{
		Login:       user.Login,
		AvatarUrl:   user.AvatarUrl,
		ProfileName: user.ProfileName,
		IsPremium:   user.IsPremium,
	}
}

func convertToCollectionResponse(collection models.Collections) *types.CommonShortType {
	return &types.CommonShortType{
		ID:   collection.Id.String(),
		Name: collection.Name,
	}
}

func convertToItemResponse(item models.CollectionItems) *types.CommonShortType {
	return &types.CommonShortType{
		ID:   item.Id.String(),
		Name: item.Name,
	}
}

func convertToWishListItemResponse(item models.WishListItems) *types.CommonShortType {
	return &types.CommonShortType{
		ID:   item.Id.String(),
		Name: item.ItemName,
	}
}
