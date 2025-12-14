package repository

import (
	"Market_backend/internal/common"
	"Market_backend/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: common.DB}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetByPaymentID(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("payment_id = ?", paymentID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}
