package services

import (
	"cmp"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"reflect"
	"slices"
)

type WLService struct {
	wlRepo    *repositories.WLRepository
	userRepo  *repositories.UsersRepository
	ciRepo    *repositories.CiRepository
	eventRepo *repositories.EventsRepository
}

func NewWLService(
	wlRepo *repositories.WLRepository,
	userRepo *repositories.UsersRepository,
	ciRepo *repositories.CiRepository,
	eventRepo *repositories.EventsRepository,
) *WLService {
	return &WLService{
		wlRepo:    wlRepo,
		userRepo:  userRepo,
		ciRepo:    ciRepo,
		eventRepo: eventRepo,
	}
}

func (s *WLService) AddItem(requestDto *models.WishListCreateRequest, userLogin string) (*models.WishListSingleDataResponse, error) {
	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	var collectionItem *models.CollectionItems

	allUsersItems, err := s.wlRepo.GetAllByUserLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var newSlice []models.WishListItems

	for _, tmpItem := range allUsersItems {
		newSlice = append(newSlice, tmpItem)
	}

	for _, wlItem := range newSlice {
		wlItem.Priority++
		_, err = s.wlRepo.UpdatePriority(&wlItem)
		if err != nil {
			return nil, err
		}
	}

	if requestDto.CollectionItemId != "" {
		collectionItem, err = s.ciRepo.GetCIByID(requestDto.CollectionItemId)
		if err != nil {
			return nil, fmt.Errorf("collection item not found")
		}
	}

	var wishListItem models.WishListItems
	if collectionItem != nil {
		wishListItem = models.WishListItems{
			UserLogin:        userLogin,
			User:             *user,
			CollectionItemID: &collectionItem.Id,
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
	err = s.eventRepo.AddEvent(
		userLogin,
		utils.EventActionCreate,
		utils.EventTargetWL,
		name,
		nil,
		nil,
		nil,
		&wlItem.Id)

	if err != nil {
		return nil, err
	}

	var resultItem models.WishListItemResponse
	err = mapstructure.Decode(wlItem, &resultItem)
	if err != nil {
		return nil, err
	}

	resultItem.CollectionItem.Collection = wlItem.CollectionItem.Collections[0].Id

	return &models.WishListSingleDataResponse{Data: resultItem}, nil
}

func (s *WLService) GetAllUserItems(userLogin string) (*models.WishListDataResponse, error) {
	allWLItems, err := s.wlRepo.GetAllByUserLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var wlItems = make([]models.WishListItemResponse, 0)

	for _, item := range allWLItems {
		var resultItem models.WishListItemResponse
		err = mapstructure.Decode(item, &resultItem)
		if err != nil {
			return nil, err
		}
		wlItems = append(wlItems, resultItem)
	}

	return &models.WishListDataResponse{
		Data: wlItems,
	}, nil
}

func (s *WLService) UpdatePriority(id string, dto models.WishListItemUpdatePriority, authUser string) (*models.WishListDataResponse, error) {
	allUsersItems, err := s.wlRepo.GetAllByUserLogin(authUser)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(allUsersItems, func(a, b models.WishListItems) int {
		return cmp.Compare(a.Priority, b.Priority)
	})

	var targetItem models.WishListItems
	var targetIndex int

	for idx, item := range allUsersItems {
		if item.Id.String() == id {
			targetItem = item
			targetIndex = idx
		}
	}

	tmpSlice := utils.RemoveOrdered(allUsersItems, targetIndex)

	var newSlice []models.WishListItems

	for _, tmpItem := range tmpSlice {
		if tmpItem.Priority > dto.Priority && targetItem.Priority < tmpItem.Priority {
			newSlice = append(newSlice, tmpItem)
			continue
		} else if tmpItem.Priority >= dto.Priority {
			tmpItem.Priority++
		}
		newSlice = append(newSlice, tmpItem)

	}

	targetItem.Priority = dto.Priority

	newSlice = append(newSlice, targetItem)

	slices.SortFunc(newSlice, func(a, b models.WishListItems) int {
		return cmp.Compare(a.Priority, b.Priority)
	})

	for _, item := range newSlice {
		_, err = s.wlRepo.UpdatePriority(&item)
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

func (s *WLService) Update(id string, req models.WishListUpdateRequest, login string) (*models.WishListSingleDataResponse, error) {
	if login == "" {
		return nil, errors.New("not authorized")
	}

	existsItem, err := s.wlRepo.GetById(id)

	fmt.Println(existsItem)
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

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			fieldName := fieldType.Name
			dstField := dst.FieldByName(fieldName)

			if dstField.IsValid() && dstField.CanSet() {
				if field.Elem().Type().AssignableTo(dstField.Type()) {
					dstField.Set(field.Elem())
				}
			}
		}
	}

	updated, err := s.wlRepo.UpdateItem(existsItem)
	if err != nil {
		return nil, err
	}

	var resultItem models.WishListItemResponse
	err = mapstructure.Decode(updated, &resultItem)
	if err != nil {
		return nil, err
	}

	return &models.WishListSingleDataResponse{Data: resultItem}, nil
}
