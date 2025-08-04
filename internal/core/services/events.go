package services

import (
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
)

type EventsService struct {
	evRepo   *repositories.EventsRepository
	userRepo *repositories.UsersRepository
}

func NewEventsService(evRepo *repositories.EventsRepository, userRepo *repositories.UsersRepository) *EventsService {
	return &EventsService{
		evRepo:   evRepo,
		userRepo: userRepo,
	}
}

func (s *EventsService) GetEvents(authUserLogin, limit, offset string) (*models.EventsDataResponse, error) {
	dbUser, err := s.userRepo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}
	var events = make([]models.Events, 0)
	subscriptions := dbUser.Subscriptions
	var totalCount int64
	if subscriptions != nil {
		events, totalCount, err = s.evRepo.GetEvents(subscriptions, limit, offset)
		if err != nil {
			return nil, err
		}
		return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
	}
	return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
}

func (s *EventsService) GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset string) (*models.EventsDataResponse, error) {
	events, totalCount, err := s.evRepo.GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
}
