package services

import (
	"cmp"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/utils"
	"reflect"
	"slices"
)

type WLService struct {
	wlRepo        *repositories.WLRepository
	userRepo      *repositories.UsersRepository
	ciRepo        *repositories.CiRepository
	eventsService *EventsService
}

func NewWLService(
	wlRepo *repositories.WLRepository,
	userRepo *repositories.UsersRepository,
	ciRepo *repositories.CiRepository,
	eventsService *EventsService,
) *WLService {
	return &WLService{
		wlRepo:        wlRepo,
		userRepo:      userRepo,
		ciRepo:        ciRepo,
		eventsService: eventsService,
	}
}

func (s *WLService) AddItem(requestDto *dto.WishListCreateRequest, userLogin string) (*dto.WishListSingleDataResponse, error) {
	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	var collectionItem *models.CollectionItems

	allUsersItems, _, err := s.wlRepo.GetAllByUserLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var newSlice []models.WishListItems

	newSlice = append(newSlice, allUsersItems...)

	for _, wlItem := range newSlice {
		wlItem.Priority++
		_, _, err = s.wlRepo.UpdatePriority(&wlItem)
		if err != nil {
			return nil, err
		}
	}

	if requestDto.CollectionItemID != "" {
		collectionItem, err = s.ciRepo.GetCIByID(requestDto.CollectionItemID)
		if err != nil {
			return nil, fmt.Errorf("collection item not found")
		}
	}

	var wishListItem models.WishListItems
	if collectionItem != nil {
		wishListItem = models.WishListItems{
			UserLogin:        userLogin,
			User:             *user,
			CollectionItemID: &collectionItem.ID,
			CollectionItem:   collectionItem,
			PurchaseLinks:    requestDto.PurchaseLinks,
			Priority:         1,
			Notes:            requestDto.Notes,
			ItemName:         collectionItem.Name,
			Images:           collectionItem.Images,
		}
	} else {
		wishListItem = models.WishListItems{
			UserLogin:     userLogin,
			User:          *user,
			ItemName:      requestDto.ItemName,
			Images:        requestDto.Images,
			PurchaseLinks: requestDto.PurchaseLinks,
			Priority:      1,
			Notes:         requestDto.Notes,
		}
	}

	var name string
	if collectionItem != nil {
		name = utils.FirstNonZero(collectionItem.Name, requestDto.ItemName)
	} else {
		name = requestDto.ItemName
	}

	wlItem, err := s.wlRepo.AddItem(&wishListItem)
	if err != nil {
		return nil, err
	}
	err = s.eventsService.AddEvent(
		userLogin,
		utils.EventActionCreate,
		utils.EventTargetWL,
		name,
		&dto.EventsParams{TargetWLID: wlItem.ID})

	if err != nil {
		return nil, err
	}

	var resultItem dto.WishListItemResponse
	err = mapstructure.Decode(wlItem, &resultItem)
	if err != nil {
		return nil, err
	}
	if collectionItem != nil {
		resultItem.CollectionItem.Collection = wlItem.CollectionItem.Collections[0].ID
	}

	return &dto.WishListSingleDataResponse{Data: resultItem}, nil
}

func (s *WLService) GetAllUserItems(userLogin string) (*dto.WishListDataResponse, error) {
	allWLItems, total, err := s.wlRepo.GetAllByUserLogin(userLogin)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var wlItems = make([]dto.WishListItemResponse, 0)

	for _, item := range allWLItems {
		var resultItem dto.WishListItemResponse
		err = mapstructure.Decode(item, &resultItem)
		if err != nil {
			return nil, err
		}
		wlItems = append(wlItems, resultItem)
	}

	return &dto.WishListDataResponse{
		Data:        wlItems,
		Total:       total,
		ProfileName: user.ProfileName,
	}, nil
}

func (s *WLService) UpdatePriority(id string, req dto.WishListItemUpdatePriority, authUser string) (*dto.WishListDataResponse, error) {
	allUsersItems, _, err := s.wlRepo.GetAllByUserLogin(authUser)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(allUsersItems, func(a, b models.WishListItems) int {
		return cmp.Compare(a.Priority, b.Priority)
	})

	var targetItem models.WishListItems
	var targetIndex int
	var found bool

	for idx, item := range allUsersItems {
		if item.ID.String() == id {
			targetItem = item
			targetIndex = idx
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("item not found")
	}

	oldPriority := targetItem.Priority
	newPriority := req.Priority

	if oldPriority == newPriority {
		return s.GetAllUserItems(authUser)
	}

	tmpSlice := append(allUsersItems[:targetIndex], allUsersItems[targetIndex+1:]...)

	var newSlice []models.WishListItems

	for _, tmpItem := range tmpSlice {
		if newPriority > oldPriority {
			if tmpItem.Priority > oldPriority && tmpItem.Priority <= newPriority {
				tmpItem.Priority--
			}
		} else {
			if tmpItem.Priority >= newPriority && tmpItem.Priority < oldPriority {
				tmpItem.Priority++
			}
		}
		newSlice = append(newSlice, tmpItem)
	}

	targetItem.Priority = newPriority
	newSlice = append(newSlice, targetItem)

	slices.SortFunc(newSlice, func(a, b models.WishListItems) int {
		return cmp.Compare(a.Priority, b.Priority)
	})

	for _, item := range newSlice {
		_, _, err = s.wlRepo.UpdatePriority(&item)
		if err != nil {
			return nil, err
		}
	}

	return s.GetAllUserItems(authUser)
}

func (s *WLService) DeleteItem(id string) (*dto.CommonResponse, error) {
	err := s.wlRepo.DeleteItem(id)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *WLService) Update(id string, req dto.WishListUpdateRequest, login string) (*dto.WishListSingleDataResponse, error) {
	if login == "" {
		return nil, errors.New("not authorized")
	}

	existsItem, err := s.wlRepo.GetByID(id)

	if err != nil {
		return nil, err
	}
	if existsItem.CollectionItemID == nil {
		existsItem.CollectionItem = nil
		existsItem.CollectionItemID = nil
	}

	dst := reflect.ValueOf(existsItem)
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}

	if dst.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", dst.Kind())
	}

	src := reflect.ValueOf(req)
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}

	if src.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", src.Kind())
	}

	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		fieldType := src.Type().Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		if (field.Kind() == reflect.Ptr && !field.IsNil()) ||
			(field.Kind() == reflect.Slice && field.Len() > 0) {
			fieldName := fieldType.Name
			dstField := dst.FieldByName(fieldName)

			if dstField.IsValid() && dstField.CanSet() {
				if field.Kind() == reflect.Ptr && field.Elem().Type().AssignableTo(dstField.Type()) {
					dstField.Set(field.Elem())
				}
				if field.Kind() == reflect.Slice && field.Type().AssignableTo(dstField.Type()) {
					dstField.Set(field)
				}
			}
		}
	}
	updated, err := s.wlRepo.UpdateItem(existsItem)
	if err != nil {
		return nil, err
	}

	var resultItem dto.WishListItemResponse
	err = mapstructure.Decode(updated, &resultItem)
	if err != nil {
		return nil, err
	}

	return &dto.WishListSingleDataResponse{Data: resultItem}, nil
}
