package dto

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
)

type PaymentDto struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID            string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Amount        string   `json:"amount"`
	TransactionID string   `json:"transactionId"`
	InnerOrderID  string   `json:"innerOrderId"`
	User          SubUsers `json:"user"`
	PaymentType   string   `json:"paymentType"`
}
type PaymentsData struct {
	Data []models.Payments `json:"data"`
}

type PaymentsTotalDonations struct {
	UserLogin      string
	TotalDonations string
}

type PaymentsInfo struct {
	Data struct {
		YearDonationsSum string       `json:"yearDonationsSum"`
		LastPayments     []PaymentDto `json:"lastPayments"`
	} `json:"data"`
}
