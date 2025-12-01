package repository

import (
	"Market_backend/internal/cart/dto"
	"Market_backend/internal/common"
	"Market_backend/internal/common/types"
	"Market_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: common.DB}
}

func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}

func (r *OrderRepository) CreateOrder(userId uuid.UUID, status types.OrderStatus, total float64) (uuid.UUID, error) {
	var orderId uuid.UUID = uuid.New()
	if err := r.db.Create(&models.Order{ID: orderId, UserID: userId, Status: status, Total: total}).Error; err != nil {
		return orderId, err
	}
	return orderId, nil
}

func (r *OrderRepository) CreateOrderTx(tx *gorm.DB, userId uuid.UUID, status types.OrderStatus, total float64) (uuid.UUID, error) {
	var orderId uuid.UUID = uuid.New()
	if err := tx.Create(&models.Order{ID: orderId, UserID: userId, Status: status, Total: total}).Error; err != nil {
		return orderId, err
	}
	return orderId, nil
}

func (r *OrderRepository) GetOrderById(orderId, userId uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("Items").Where("id = ? AND user_id = ?", orderId, userId).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetAllOrders(userId uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Find(&orders, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) CreateOrderItems(orderId uuid.UUID, createOrderItem []dto.GetCartItemsResponse) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range createOrderItem {
			if err := tx.Create(&models.OrderItem{
				ID:          uuid.New(),
				OrderID:     orderId,
				Quantity:    item.Quantity,
				ProductType: item.ProductType,
				ProductID:   item.ProductId,
				UnitPrice:   item.Price,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OrderRepository) CreateOrderItemsTx(tx *gorm.DB, orderId uuid.UUID, createOrderItem []dto.GetCartItemsResponse) error {
	for _, item := range createOrderItem {
		if err := tx.Create(&models.OrderItem{
			ID:          uuid.New(),
			OrderID:     orderId,
			Quantity:    item.Quantity,
			ProductType: item.ProductType,
			ProductID:   item.ProductId,
			UnitPrice:   item.Price}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepository) ChangeStatus(orderId uuid.UUID, status types.OrderStatus) error {
	if err := r.db.Model(&models.Order{}).Where("id = ?", orderId).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
