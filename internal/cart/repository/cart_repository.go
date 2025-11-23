package repository

import (
	"eduVix_backend/internal/cart/dto"
	"eduVix_backend/internal/common"
	"eduVix_backend/internal/common/types"
	"eduVix_backend/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{db: common.DB}
}

func (r *CartRepository) CreateCart(userId uuid.UUID) error {
	return r.db.Create(&models.Cart{
		ID:     uuid.New(),
		UserID: userId,
	}).Error
}

func (r *CartRepository) AddNewCartItem(cartItem dto.CartItemDto, price float64) (uuid.UUID, error) {
	var existingItem models.CartItem

	err := r.db.
		Where("cart_id = ? AND product_id = ? AND product_type = ?",
			cartItem.CartID,
			cartItem.ProductId,
			cartItem.ProductType,
		).
		First(&existingItem).Error

	// Если уже есть, просто обновляем количество
	if err == nil {
		existingItem.Quantity += cartItem.Quantity
		existingItem.UnitPrice = price
		err = r.db.Save(&existingItem).Error
		return existingItem.ID, err
	}

	// Если не найден — создаём новую запись
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	newItem := models.CartItem{
		ID:          uuid.New(),
		CartID:      cartItem.CartID,
		ProductID:   cartItem.ProductId,
		ProductType: cartItem.ProductType,
		Quantity:    cartItem.Quantity,
		UnitPrice:   price,
	}

	err = r.db.Create(&newItem).Error
	return newItem.ID, err
}

func (r *CartRepository) RemoveCartItem(cartItemId, userId uuid.UUID) error {
	res := r.db.
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", cartItemId, userId).
		Delete(&models.CartItem{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("cart item not found or does not belong to user")
	}

	return nil
}

func (r *CartRepository) ChangeCartItem(cartItemId, userId uuid.UUID, quantity int) error {
	res := r.db.
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", cartItemId, userId).
		Update("quantity", quantity)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("cart item not found or does not belong to user")
	}

	return nil
}

func (r *CartRepository) GetAllCartItems(userId, cartId uuid.UUID) ([]dto.GetCartItemsResponse, error) {
	// 1. Проверяем, что корзина принадлежит пользователю
	var cart models.Cart
	if err := r.db.First(&cart, "id = ? AND user_id = ?", cartId, userId).Error; err != nil {
		return nil, err
	}

	// 2. Получаем все позиции корзины
	var cartItems []models.CartItem
	if err := r.db.Where("cart_id = ?", cartId).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// 3. Группируем ID товаров по типам
	var procIDs []uuid.UUID
	var flashIDs []uuid.UUID

	for _, ci := range cartItems {
		switch ci.ProductType {
		case types.Processor:
			procIDs = append(procIDs, ci.ProductID)
		case types.FlashDriver:
			flashIDs = append(flashIDs, ci.ProductID)
		}
	}

	// 4. Грузим изображения товаров
	imageMap := make(map[uuid.UUID]string)

	// Процессоры
	if len(procIDs) > 0 {
		var images []models.Image
		if err := r.db.Where("processor_id IN ?", procIDs).Find(&images).Error; err == nil {
			for _, img := range images {
				if img.ProcessorID != nil {
					imageMap[*img.ProcessorID] = img.URL
				}
			}
		}
	}

	// Флешки
	if len(flashIDs) > 0 {
		var images []models.Image
		if err := r.db.Where("flash_drive_id IN ?", flashIDs).Find(&images).Error; err == nil {
			for _, img := range images {
				if img.FlashDriveID != nil {
					imageMap[*img.FlashDriveID] = img.URL
				}
			}
		}
	}

	// 5. Формируем итоговый ответ
	resultItems := make([]dto.GetCartItemsResponse, 0, len(cartItems))

	for _, ci := range cartItems {
		resultItems = append(resultItems, dto.GetCartItemsResponse{
			ProductId:   ci.ProductID,
			ProductType: ci.ProductType,
			Quantity:    ci.Quantity,
			Price:       ci.UnitPrice,
			ImageUrl:    imageMap[ci.ProductID], // O(1)
		})
	}

	return resultItems, nil
}
