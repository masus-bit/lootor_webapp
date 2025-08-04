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

func (s *EventsService) GetEvents(authUserLogin string) (*models.EventsDataResponse, error) {
	dbUser, err := s.userRepo.GetUserByLogin(authUserLogin)
	if err != nil {
		return nil, err
	}
	var events = make([]models.Events, 0)
	subscriptions := dbUser.Subscriptions
	if subscriptions != nil {
		events, err = s.evRepo.GetEvents(subscriptions)
		if err != nil {
			return nil, err
		}
		return &models.EventsDataResponse{Data: events}, nil
	}
	return &models.EventsDataResponse{Data: events}, nil
}

func (s *EventsService) GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset string) (*models.EventsDataResponse, error) {
	events, totalCount, err := s.evRepo.GetFilteredEvents(userLogin, collectionId, collectionItem, wlId, limit, offset)
	if err != nil {
		return nil, err
	}
	return &models.EventsDataResponse{Data: events, Total: totalCount}, nil
}
