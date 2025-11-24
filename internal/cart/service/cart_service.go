package service

import (
	"eduVix_backend/internal/cart/dto"
	"eduVix_backend/internal/cart/repository"
	ProductRepo "eduVix_backend/internal/product/repository"
	"github.com/google/uuid"
)

type CartService struct {
	repo     *repository.CartRepository
	procRepo *ProductRepo.ProcessorRepository
}

func NewCartService(repo *repository.CartRepository, procRepo *ProductRepo.ProcessorRepository) *CartService {
	return &CartService{repo: repo, procRepo: procRepo}
}

func (s *CartService) AddNewItem(cartItem dto.CartItemDto) (uuid.UUID, error) {
	var currentPrice float64
	switch cartItem.ProductType {
	case "P":
		proc, err := s.procRepo.GetProcessorById(cartItem.ProductId)
		if err != nil {
			return uuid.Nil, err
		}

		if proc != nil {
			currentPrice = proc.WholesalePrice
			if cartItem.Quantity < proc.WholesaleMinQty {
				currentPrice = proc.RetailPrice
			}
		}
	case "FD":

	}

	return s.repo.AddNewCartItem(cartItem, currentPrice)
}

func (s *CartService) RemoveItem(cartItemId, userId uuid.UUID) error {
	return s.repo.RemoveCartItem(cartItemId, userId)
}

func (s *CartService) ChangeItem(cartItemId, userId uuid.UUID, quantity int) error {
	return s.repo.ChangeQuantity(cartItemId, userId, quantity)
}

func (s *CartService) GetAllCartItems(cartItemId, userId uuid.UUID) ([]dto.GetCartItemsResponse, error) {
	return s.repo.GetAllCartItems(cartItemId, userId)
}
