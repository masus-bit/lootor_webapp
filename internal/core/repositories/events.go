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
		Date:            time.Now().Format(time.RFC3339),
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

func (r *EventsRepository) GetEvents(subscriptions []string, limit, offset string) ([]models.Events, int64, error) {
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var count int64

	var events []models.Events

	query := r.db.Model(&models.Events{}).
		Where("LOWER(initiator_login) IN ?", subscriptions).
		Preload("Initiator").
		Order("date DESC")

	query.Count(&count)

	if err := query.Limit(intLimit).Offset(intOffset).Find(&events).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}

	if err := r.loadEventRelations(&events); err != nil {
		return nil, 0, fmt.Errorf("failed to load event relations: %w", err)
	}

	return events, count, nil
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
		Where("LOWER(initiator_login) = LOWER(?)", userLogin).Preload("Initiator").
		Order("date DESC")

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
		if err := r.db.Where("id IN ?", itemIDs).Preload("Collections").Find(&items).Error; err != nil {
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

func convertToCollectionResponse(collection models.Collections) *types.CommonShortTypeCollection {
	return &types.CommonShortTypeCollection{
		CommonShortType: &types.CommonShortType{
			ID:   collection.Id.String(),
			Name: collection.Name,
		},
		Owner: collection.UserLogin,
	}
}

func convertToItemResponse(item models.CollectionItems) *types.CommonShortTypeItem {
	return &types.CommonShortTypeItem{
		CommonShortType: &types.CommonShortType{
			ID:   item.Id.String(),
			Name: item.Name,
		},
		Collection: item.Collections[0].Transliteration,
		Owner:      item.UserLogin,
	}
}

func convertToWishListItemResponse(item models.WishListItems) *types.CommonShortType {
	return &types.CommonShortType{
		ID:   item.Id.String(),
		Name: item.ItemName,
	}
}
