package service

import (
	CartRepository "Market_backend/internal/cart/repository"
	CartService "Market_backend/internal/cart/service"
	"Market_backend/internal/common/types"
	"Market_backend/internal/order/repository"
	"Market_backend/internal/product/service"
	"Market_backend/models"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repo        *repository.OrderRepository
	cartRepo    *CartRepository.CartRepository
	cartService *CartService.CartService

	ProcService  *service.ProcessorService
	FlashService *service.FlashDriveService
}

func NewOrderService(
	repo *repository.OrderRepository,
	cartRepo *CartRepository.CartRepository,
	cartService *CartService.CartService,
	procS *service.ProcessorService,
	flashS *service.FlashDriveService,
) *OrderService {
	return &OrderService{repo: repo, cartRepo: cartRepo, cartService: cartService, ProcService: procS, FlashService: flashS}
}

func (s *OrderService) CreateOrder(userId, cartId uuid.UUID) (uuid.UUID, error) {
	var orderId uuid.UUID

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		totalBalance, err := s.cartService.ValidateCartTx(tx, userId, cartId)

		if err != nil {
			return err
		}

		items, err := s.cartRepo.GetAllCartItemsTx(tx, userId, cartId)
		if err != nil {
			return err
		}

		orderId, err = s.repo.CreateOrderTx(tx, userId, types.InProgress, totalBalance)
		if err != nil {
			return err
		}

		if err = s.repo.CreateOrderItemsTx(tx, orderId, items); err != nil {
			return err
		}

		if err = s.cartRepo.ClearCartTx(tx, userId, cartId); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return orderId, nil
}

func (s *OrderService) ChangeOrderStatus(orderId uuid.UUID, status types.OrderStatus) error {
	return s.repo.ChangeStatus(orderId, status)
}

func (s *OrderService) CancelOrder(orderId uuid.UUID) error {
	return s.repo.ChangeStatus(orderId, types.Cancelled)
}

func (s *OrderService) GetOrderById(userId, orderId uuid.UUID) (*models.Order, error) {
	return s.repo.GetOrderById(orderId, userId)
}

//func (s *OrderService) GetAllOrders(userId uuid.UUID) ([]models.Order, error) {
//	return s.repo.GetAllOrders(userId)
//}

func (s *OrderService) GetAllUserOrders(userId uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	err := s.repo.DB().
		Where("user_id = ?", userId).
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	err := s.repo.DB().
		Preload("Items").
		Preload("User"). // подтянуть имя покупателя
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (s *OrderService) UpdateOrderStatus(orderId uuid.UUID, newStatus types.OrderStatus) error {
	var order models.Order

	// Находим заказ
	if err := s.repo.DB().
		Preload("Items").
		First(&order, "id = ?", orderId).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Обновляем статус
	order.Status = newStatus

	// Если нужно вычитать остатки со склада при оплате
	if newStatus == types.Completed {
		for _, item := range order.Items {
			switch item.ProductType {
			case types.Processor:
				if err := s.ProcService.DB().UpdateStock(item.ProductID, -item.Quantity); err != nil {
					return fmt.Errorf("cannot update processor stock: %w", err)
				}
			case types.FlashDriver:
				if err := s.FlashService.DB().UpdateStock(item.ProductID, -item.Quantity); err != nil {
					return fmt.Errorf("cannot update flash drive stock: %w", err)
				}
			}
		}
	}

	// Сохраняем изменения
	if err := s.repo.DB().Save(&order).Error; err != nil {
		return fmt.Errorf("cannot save order: %w", err)
	}

	return nil
}
