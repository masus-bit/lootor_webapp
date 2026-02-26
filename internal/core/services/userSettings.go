package services

import (
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
)

type UserSettingsService struct {
	userSettingsRepo *repositories.UserSettingsRepository
}

func NewUserSettingsService(
	userSettingsRepo *repositories.UserSettingsRepository,
) *UserSettingsService {
	return &UserSettingsService{
		userSettingsRepo: userSettingsRepo,
	}
}

func (s *UserSettingsService) CreateSettings(userLogin string, settings *dto.UserSettingsRequest) (
	*dto.CommonResponse,
	error,
) {
	var userSettings *models.UserSettings

	userSettings = &models.UserSettings{
		UserLogin: userLogin,
		Settings:  settings.Settings,
	}
	_, err := s.userSettingsRepo.CreateRecord(userSettings)
	if err != nil {
		return nil, err
	}

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *UserSettingsService) UpdateSettings(userLogin string, settings *dto.UserSettingsRequest) (
	*dto.CommonResponse,
	error,
) {
	existsSettings, _ := s.userSettingsRepo.FindRecordByUserLogin(userLogin)

	if existsSettings == nil {
		var userSettings *models.UserSettings

		userSettings = &models.UserSettings{
			UserLogin: userLogin,
			Settings:  settings.Settings,
		}
		_, err := s.userSettingsRepo.CreateRecord(userSettings)
		if err != nil {
			return nil, err
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	} else {
		existsSettings.Settings = settings.Settings

		_, err := s.userSettingsRepo.UpdateFull(existsSettings)
		if err != nil {
			return nil, err
		}
		return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
	}

}
