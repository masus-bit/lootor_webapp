package repositories

import (
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
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

func (r *PaymentsRepository) FindAllPaymentsByLogin(login string) ([]models.Payments, error) {
	var payments []models.Payments
	err := r.db.Where("user_login = ?", login).Find(&payments).Error
	return payments, err
}

func (r *PaymentsRepository) FindYearPayments() (string, []models.Payments, error) {
	var donations string
	var payments []models.Payments
	currentYear := time.Now().Year()

	err := r.db.Model(&models.Payments{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", currentYear).
		Select("COALESCE(SUM(CAST(NULLIF(amount, '') AS NUMERIC)), 0) as total_donations").
		Scan(&donations).Error

	if err != nil {
		return "", payments, err
	}

	err = r.db.Model(&models.Payments{}).
		Order("created_at DESC").
		Limit(10).
		Find(&payments).Error
	if err != nil {
		return "", payments, err
	}

	return donations, payments, nil
}
