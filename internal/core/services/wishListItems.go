package services

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
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

func (s *WLService) AddItem(requestDto *models.WishListCreateRequest, userLogin string) (*dto.CommonResponse, error) {
	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	collectionItem, err := s.ciRepo.GetCIByID(requestDto.CollectionItemId)
	if err != nil {
		return nil, fmt.Errorf("colelction item not найден")
	}

	wishListItem := &models.WishListItems{
		UserLogin:        userLogin,
		User:             *user,
		CollectionItemID: &collectionItem.Id,
		CollectionItem:   collectionItem,
		PurchaseLinks:    requestDto.PurchaseLinks,
		Priority:         requestDto.Priority,
		Notes:            requestDto.Notes,
	}
	wlItem, err := s.wlRepo.AddItem(wishListItem)
	if err != nil {
		return nil, err
	}
	err = s.eventRepo.AddEvent(
		userLogin,
		utils.EventActionCreate,
		utils.EventTargetWL,
		collectionItem.Name,
		nil,
		nil,
		nil,
		&wlItem.Id)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
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
