package service

import (
	"Market_backend/internal/cart/dto"
	"Market_backend/internal/cart/repository"
	"Market_backend/internal/common/types"
	ProductRepo "Market_backend/internal/product/repository"

	"Market_backend/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService struct {
	repo      *repository.CartRepository
	flashRepo *ProductRepo.FlashDriveRepository
	procRepo  *ProductRepo.ProcessorRepository
}

func NewCartService(repo *repository.CartRepository, procRepo *ProductRepo.ProcessorRepository, flashRepo *ProductRepo.FlashDriveRepository) *CartService {
	return &CartService{repo: repo, procRepo: procRepo, flashRepo: flashRepo}
}

func (s *CartService) AddNewItem(cartItem dto.CartItemDto) (uuid.UUID, error) {
	var currentPrice float64
	var stock int

	switch cartItem.ProductType {
	case "P": // Processor
		proc, err := s.procRepo.GetProcessorById(cartItem.ProductId)
		if err != nil {
			return uuid.Nil, err
		}
		if proc == nil {
			return uuid.Nil, fmt.Errorf("processor not found")
		}

		stock = proc.Stock
		if cartItem.Quantity > stock {
			return uuid.Nil, fmt.Errorf("cannot add %d items, only %d in stock", cartItem.Quantity, stock)
		}

		if cartItem.Quantity >= proc.WholesaleMinQty {
			currentPrice = proc.WholesalePrice
		} else {
			currentPrice = proc.RetailPrice
		}

	case "FD": // FlashDrive
		flash, err := s.flashRepo.GetFlashDriveById(cartItem.ProductId)
		if err != nil {
			return uuid.Nil, err
		}
		if flash == nil {
			return uuid.Nil, fmt.Errorf("flash drive not found")
		}

		stock = flash.Stock
		if cartItem.Quantity > stock {
			return uuid.Nil, fmt.Errorf("cannot add %d items, only %d in stock", cartItem.Quantity, stock)
		}

		if cartItem.Quantity >= flash.WholesaleMinQty {
			currentPrice = flash.WholesalePrice
		} else {
			currentPrice = flash.RetailPrice
		}

	default:
		return uuid.Nil, fmt.Errorf("unknown product type: %s", cartItem.ProductType)
	}

	return s.repo.AddNewCartItem(cartItem, currentPrice)
}

func (s *CartService) RemoveItem(cartItemId, userId uuid.UUID) error {
	return s.repo.RemoveCartItem(cartItemId, userId)
}

func (s *CartService) ChangeItem(cartItemId, userId uuid.UUID, quantity int) error {
	cartItem, err := s.repo.GetItem(userId, cartItemId)
	if err != nil {
		return err
	}

	if quantity <= 0 {
		return s.repo.RemoveCartItem(cartItemId, userId)
	}

	var newPrice float64

	switch p := cartItem.Product.(type) {
	case *models.Processor:
		newPrice = p.RetailPrice
		if quantity >= p.WholesaleMinQty {
			newPrice = p.WholesalePrice
		}
	case *models.FlashDrive:
		newPrice = p.RetailPrice
		if quantity >= p.WholesaleMinQty {
			newPrice = p.WholesalePrice
		}
	default:
		return fmt.Errorf("unknown product type")
	}

	return s.repo.ChangeQuantity(cartItemId, userId, quantity, newPrice)
}

func (s *CartService) GetAllCartItems(cartItemId, userId uuid.UUID) ([]dto.GetCartItemsResponse, error) {
	return s.repo.GetAllCartItems(cartItemId, userId)
}

func (s *CartService) ClearCart(userId, cartId uuid.UUID) error {
	return s.repo.ClearCart(userId, cartId)
}

func (s *CartService) ValidateCart(userId, cartId uuid.UUID) (float64, error) {
	cartItems, err := s.repo.GetAllCartItems(userId, cartId)
	if err != nil {
		return 0, err
	}

	if len(cartItems) == 0 {
		return 0, errors.New("корзина пуста")
	}

	total := 0.0
	for _, ci := range cartItems {
		var price float64
		switch ci.ProductType {
		case types.Processor:
			proc, err := s.procRepo.GetProcessorById(ci.ProductId)
			if err != nil {
				return 0, err
			}
			if ci.Quantity > proc.Stock {
				return 0, fmt.Errorf("товара %s не хватает на складе", proc.Name)
			}
			if ci.Quantity >= proc.WholesaleMinQty {
				price = proc.WholesalePrice
			} else {
				price = proc.RetailPrice
			}
		case types.FlashDriver:
			flash, err := s.flashRepo.GetFlashDriveById(ci.ProductId)
			if err != nil {
				return 0, err
			}
			if ci.Quantity > flash.Stock {
				return 0, fmt.Errorf("товара %s не хватает на складе", flash.Name)
			}
			if ci.Quantity >= flash.WholesaleMinQty {
				price = flash.WholesalePrice
			} else {
				price = flash.RetailPrice
			}
		default:
			return 0, fmt.Errorf("unknown product type: %s", ci.ProductType)
		}
		total += price * float64(ci.Quantity)
	}

	return total, nil
}

func (s *CartService) ValidateCartTx(tx *gorm.DB, userId, cartId uuid.UUID) (float64, error) {
	cartItems, err := s.repo.GetAllCartItemsTx(tx, userId, cartId)
	if err != nil {
		return 0, err
	}

	if len(cartItems) == 0 {
		return 0, errors.New("корзина пуста")
	}

	total := 0.0
	for _, ci := range cartItems {
		var price float64

		switch ci.ProductType {
		case types.Processor:
			proc, err := s.procRepo.GetProcessorByIdTx(tx, ci.ProductId)
			if err != nil {
				return 0, err
			}
			if ci.Quantity > proc.Stock {
				return 0, fmt.Errorf("товара %s не хватает на складе", proc.Name)
			}
			if ci.Quantity >= proc.WholesaleMinQty {
				price = proc.WholesalePrice
			} else {
				price = proc.RetailPrice
			}
		case types.FlashDriver:
			flash, err := s.flashRepo.GetFlashDriveByIdTx(tx, ci.ProductId)
			if err != nil {
				return 0, err
			}
			if ci.Quantity > flash.Stock {
				return 0, fmt.Errorf("товара %s не хватает на складе", flash.Name)
			}
			if ci.Quantity >= flash.WholesaleMinQty {
				price = flash.WholesalePrice
			} else {
				price = flash.RetailPrice
			}
		default:
			return 0, fmt.Errorf("unknown product type: %s", ci.ProductType)
		}

		total += price * float64(ci.Quantity)
	}

	return total, nil
}

func (s *CartService) DeductStockAfterPayment(orderId uuid.UUID) error {
	items, err := s.repo.GetCartItemsByOrder(orderId)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ProductType {
		case types.Processor:
			proc, err := s.procRepo.GetProcessorById(item.ProductID)
			if err != nil {
				return err
			}
			newStock := proc.Stock - item.Quantity
			if newStock < 0 {
				return fmt.Errorf("processor %s: not enough stock", proc.Name)
			}
			if err := s.procRepo.UpdateStock(proc.ID, newStock); err != nil {
				return err
			}

		case types.FlashDriver:
			fd, err := s.flashRepo.GetFlashDriveById(item.ProductID)
			if err != nil {
				return err
			}
			newStock := fd.Stock - item.Quantity
			if newStock < 0 {
				return fmt.Errorf("flash drive %s: not enough stock", fd.Name)
			}
			if err := s.flashRepo.UpdateStock(fd.ID, newStock); err != nil {
				return err
			}
		}
	}

	return nil
}
