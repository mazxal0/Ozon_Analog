package service

import (
	CartRepository "eduVix_backend/internal/cart/repository"
	CartService "eduVix_backend/internal/cart/service"
	"eduVix_backend/internal/common/types"
	"eduVix_backend/internal/order/repository"
	"eduVix_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	repo        *repository.OrderRepository
	cartRepo    *CartRepository.CartRepository
	cartService *CartService.CartService
}

func NewOrderService(repo *repository.OrderRepository, cartRepo *CartRepository.CartRepository, cartService *CartService.CartService) *OrderService {
	return &OrderService{repo: repo, cartRepo: cartRepo, cartService: cartService}
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

func (s *OrderService) GetAllOrders(userId uuid.UUID) ([]models.Order, error) {
	return s.repo.GetAllOrders(userId)
}
