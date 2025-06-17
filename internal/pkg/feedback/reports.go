package feedback

import (
	"github.com/joho/godotenv"
	"lootor/internal/core/repositories"
	"lootor/internal/pkg/dto"
)

type ReportsService struct {
	feedService    *FeedbackService
	ciRepo         *repositories.CiRepository
	userRepo       *repositories.UsersRepository
	collectionRepo *repositories.CollectionsRepository
}

func NewReportsService(
	feedService *FeedbackService,
	ciRepo *repositories.CiRepository,
	userRepo *repositories.UsersRepository,
	collectionRepo *repositories.CollectionsRepository) *ReportsService {
	_ = godotenv.Load()
	return &ReportsService{feedService: feedService, ciRepo: ciRepo, userRepo: userRepo, collectionRepo: collectionRepo}
}

func (s *ReportsService) ReportAnything(id string) (*dto.CommonResponse, error) {
	title := "Жалоба"
	var description string
	collection, err := s.collectionRepo.GetByIdWithoutCollectionItems(id)
	if err != nil {
		collection = nil
	}
	collectionItem, err := s.ciRepo.GetCIByID(id)
	if err != nil {
		collectionItem = nil
	}
	user, err := s.userRepo.GetUserByLogin(id)
	if err != nil {
		user = nil
	}

	if collection != nil {
		if collection.ReportsCount == 2 {
			description = description + "\n\n\n__________\n\nКоллекция: " + collection.Id.String() + "\n\nTransliteration: " + collection.Transliteration + "\n\nВладелец: " + collection.User.Login
			_, err := s.feedService.SendTextTicket(title, description, true)
			if err != nil {
				return nil, err
			}
			collection.ReportsCount = 0
			_, err = s.collectionRepo.UpdateCollectionFull(collection)
			if err != nil {
				return nil, err
			}
		}
		collection.ReportsCount = collection.ReportsCount + 1
		_, err := s.collectionRepo.UpdateCollectionFull(collection)
		if err != nil {
			return nil, err
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}

	if collectionItem != nil {
		if collectionItem.ReportsCount == 2 {
			description = description + "\n\n\n__________\n\nЭкземпляр коллекции: " + collectionItem.Id.String() + "\n\nВладелец: " + collectionItem.Owner.Login
			_, err := s.feedService.SendTextTicket(title, description, true)
			if err != nil {
				return nil, err
			}
			collectionItem.ReportsCount = 0
			_, err = s.ciRepo.UpdateCI(collectionItem, collectionItem)
			if err != nil {
				return nil, err
			}
		}
		collectionItem.ReportsCount = collectionItem.ReportsCount + 1
		_, err := s.ciRepo.UpdateCI(collectionItem, collectionItem)
		if err != nil {
			return nil, err
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}

	if user != nil {
		if user.ReportsCount == 2 {
			description = description + "\n\n\n__________\n\nПользователь: " + user.Login
			_, err := s.feedService.SendTextTicket(title, description, true)
			if err != nil {
				return nil, err
			}
			user.ReportsCount = 0
			_, err = s.userRepo.UpdateUser(user, *user)
			if err != nil {
				return nil, err
			}
		}
		user.ReportsCount = user.ReportsCount + 1
		_, err := s.userRepo.UpdateUser(user, *user)
		if err != nil {
			return nil, err
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: false}}, nil
}
