package dto

import "lootor/internal/core/models"

type PlatformsDataResponse struct {
	Data []models.Platforms `json:"data"`
}
