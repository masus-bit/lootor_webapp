package dto

import "lootor/internal/core/models"

type PaymentsData struct {
	Data []models.Payments `json:"data"`
}

type PaymentsTotalDonations struct {
	UserLogin      string
	TotalDonations string
}
