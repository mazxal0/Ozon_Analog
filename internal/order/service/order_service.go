package service

import (
	CartRepository "eduVix_backend/internal/cart/repository"
	CartService "eduVix_backend/internal/cart/service"
	"eduVix_backend/internal/order/repository"
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

		orderId, err = s.repo.CreateOrderTx(tx, userId, "in_progress", totalBalance)
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
