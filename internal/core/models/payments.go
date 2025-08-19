package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Payments struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Amount        string    `json:"amount"`
	TransactionId string    `json:"transactionId"`
	InnerOrderId  string    `json:"innerOrderId"`
	UserLogin     string    `json:"userLogin"`
	PaymentType   string    `json:"paymentType"`
}

type PaymentsData struct {
	Data []Payments `json:"data"`
}

type PaymentsTotalDonations struct {
	UserLogin      string
	TotalDonations string
}
