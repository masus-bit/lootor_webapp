package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
)

type PaymentsRepository struct {
	db *gorm.DB
}

func NewPaymentsRepository(db *gorm.DB) *PaymentsRepository {
	return &PaymentsRepository{db: db}
}

func (r *PaymentsRepository) CreateRecord(payment *models.Payments) error {
	err := r.db.Create(payment)
	if err.Error != nil {
		return err.Error
	}

	return nil
}
