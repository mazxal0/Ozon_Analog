package repository

import (
	"Market_backend/internal/cart/dto"
	"Market_backend/internal/common"
	"Market_backend/internal/common/types"
	"Market_backend/models"
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

func (r *CartRepository) CreateCart(userId uuid.UUID) (uuid.UUID, error) {
	cart := models.Cart{
		ID:     uuid.New(),
		UserID: userId,
	}

	if err := r.db.Create(&cart).Error; err != nil {
		return uuid.Nil, err
	}

	return cart.ID, nil
}

func (r *CartRepository) GetItem(userId, itemId uuid.UUID) (*dto.CartItemWithProduct, error) {
	var item models.CartItem

	// 1. Загружаем сам CartItem
	err := r.db.
		Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)", itemId, userId).
		First(&item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("item was not found")
	}

	if err != nil {
		return nil, err
	}

	// 2. Подгружаем продукт по типу
	var product any

	switch item.ProductType {
	case types.Processor:
		var processor models.Processor
		if err := r.db.Preload("Images").First(&processor, "id = ?", item.ProductID).Error; err == nil {
			product = &processor
		}

	case types.FlashDriver:
		var flash models.FlashDrive
		if err := r.db.Preload("Images").First(&flash, "id = ?", item.ProductID).Error; err == nil {
			product = &flash
		}
	}

	// 3. Возвращаем DTO
	return &dto.CartItemWithProduct{
		CartItem: item,
		Product:  product,
	}, nil
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
		Exec(`
            DELETE FROM cart_items
            USING carts
            WHERE cart_items.id = ?
              AND cart_items.cart_id = carts.id
              AND carts.user_id = ?
        `, cartItemId, userId)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("cart item not found or does not belong to user")
	}

	return nil
}

func (r *CartRepository) ChangeQuantity(cartItemId, userId uuid.UUID, quantity int, unitPrice float64) error {
	res := r.db.Model(&models.CartItem{}).
		Where("id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)", cartItemId, userId).
		Updates(map[string]interface{}{
			"quantity":   quantity,
			"unit_price": unitPrice,
		})

	return res.Error
}

func (r *CartRepository) GetAllCartItems(userId, cartId uuid.UUID) ([]dto.GetCartItemsResponse, error) {
	// 1. Проверяем владельца корзины
	var cart models.Cart
	if err := r.db.First(&cart, "id = ? AND user_id = ?", cartId, userId).Error; err != nil {
		return nil, err
	}

	// 2. Карточные элементы
	var cartItems []models.CartItem
	if err := r.db.Where("cart_id = ?", cartId).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// Собираем ID по категориям
	var procIDs, flashIDs []uuid.UUID
	for _, ci := range cartItems {
		switch ci.ProductType {
		case types.Processor:
			procIDs = append(procIDs, ci.ProductID)
		case types.FlashDriver:
			flashIDs = append(flashIDs, ci.ProductID)
		}
	}

	// 3. Загружаем товары
	procMap := map[uuid.UUID]models.Processor{}
	flashMap := map[uuid.UUID]models.FlashDrive{}
	imageMap := map[uuid.UUID]string{}
	priceMap := map[uuid.UUID]float64{}

	// PROCESSORS
	if len(procIDs) > 0 {
		var procs []models.Processor
		if err := r.db.Preload("Images").Where("id IN ?", procIDs).Find(&procs).Error; err == nil {
			for _, p := range procs {
				procMap[p.ID] = p
				// цена по умолчанию — RetailPrice
				priceMap[p.ID] = p.RetailPrice
				if len(p.Images) > 0 {
					imageMap[p.ID] = p.Images[0].URL
				}
			}
		}
	}

	// FLASH DRIVES
	if len(flashIDs) > 0 {
		var flash []models.FlashDrive
		if err := r.db.Preload("Images").Where("id IN ?", flashIDs).Find(&flash).Error; err == nil {
			for _, f := range flash {
				flashMap[f.ID] = f
				// цена по умолчанию — RetailPrice
				priceMap[f.ID] = f.RetailPrice
				if len(f.Images) > 0 {
					imageMap[f.ID] = f.Images[0].URL
				}
			}
		}
	}

	// 4. Формируем ответ с учетом количества
	result := make([]dto.GetCartItemsResponse, 0, len(cartItems))

	for _, ci := range cartItems {
		price := priceMap[ci.ProductID]
		name := ""

		switch ci.ProductType {
		case types.Processor:
			if proc, ok := procMap[ci.ProductID]; ok {
				name = proc.Name
				if ci.Quantity >= proc.WholesaleMinQty {
					price = proc.WholesalePrice
				}
			}
		case types.FlashDriver:
			if flash, ok := flashMap[ci.ProductID]; ok {
				name = flash.Name
				if ci.Quantity >= flash.WholesaleMinQty {
					price = flash.WholesalePrice
				}
			}
		}

		result = append(result, dto.GetCartItemsResponse{
			ID:          ci.ID,
			ProductId:   ci.ProductID,
			ProductType: ci.ProductType,
			Quantity:    ci.Quantity,
			Price:       price,
			Name:        name, // <-- добавлено
			ImageUrl:    imageMap[ci.ProductID],
		})
	}

	return result, nil
}

func (r *CartRepository) GetAllCartItemsTx(tx *gorm.DB, userId, cartId uuid.UUID) ([]dto.GetCartItemsResponse, error) {
	// 1. Проверяем владельца корзины
	var cart models.Cart
	if err := tx.First(&cart, "id = ? AND user_id = ?", cartId, userId).Error; err != nil {
		return nil, err
	}

	// 2. Карточные элементы
	var cartItems []models.CartItem
	if err := tx.Where("cart_id = ?", cartId).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// Собираем ID по категориям
	var procIDs, flashIDs []uuid.UUID
	for _, ci := range cartItems {
		switch ci.ProductType {
		case types.Processor:
			procIDs = append(procIDs, ci.ProductID)
		case types.FlashDriver:
			flashIDs = append(flashIDs, ci.ProductID)
		}
	}

	// 3. Загружаем товары
	procMap := map[uuid.UUID]models.Processor{}
	flashMap := map[uuid.UUID]models.FlashDrive{}
	imageMap := map[uuid.UUID]string{}
	priceMap := map[uuid.UUID]float64{}

	// PROCESSORS
	if len(procIDs) > 0 {
		var procs []models.Processor
		if err := tx.Preload("Images").Where("id IN ?", procIDs).Find(&procs).Error; err == nil {
			for _, p := range procs {
				procMap[p.ID] = p
				priceMap[p.ID] = p.RetailPrice
				if len(p.Images) > 0 {
					imageMap[p.ID] = p.Images[0].URL
				}
			}
		}
	}

	// FLASH DRIVES
	if len(flashIDs) > 0 {
		var flash []models.FlashDrive
		if err := tx.Preload("Images").Where("id IN ?", flashIDs).Find(&flash).Error; err == nil {
			for _, f := range flash {
				flashMap[f.ID] = f
				priceMap[f.ID] = f.RetailPrice
				if len(f.Images) > 0 {
					imageMap[f.ID] = f.Images[0].URL
				}
			}
		}
	}

	// 4. Формируем ответ
	result := make([]dto.GetCartItemsResponse, 0, len(cartItems))

	for _, ci := range cartItems {
		price := priceMap[ci.ProductID]

		switch ci.ProductType {
		case types.Processor:
			if proc, ok := procMap[ci.ProductID]; ok && ci.Quantity >= proc.WholesaleMinQty {
				price = proc.WholesalePrice
			}
		case types.FlashDriver:
			if flash, ok := flashMap[ci.ProductID]; ok && ci.Quantity >= flash.WholesaleMinQty {
				price = flash.WholesalePrice
			}
		}

		result = append(result, dto.GetCartItemsResponse{
			ID:          ci.ID,
			ProductId:   ci.ProductID,
			ProductType: ci.ProductType,
			Quantity:    ci.Quantity,
			Price:       price,
			ImageUrl:    imageMap[ci.ProductID],
		})
	}

	return result, nil
}

func (r *CartRepository) ClearCart(userId, cartId uuid.UUID) error {
	return r.db.
		Where("cart_id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)", cartId, userId).
		Delete(&models.CartItem{}).
		Error
}

func (r *CartRepository) ClearCartTx(tx *gorm.DB, userId, cartId uuid.UUID) error {
	return tx.
		Where("cart_id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)", cartId, userId).
		Delete(&models.CartItem{}).
		Error
}

func (r *CartRepository) GetCartItemsByOrder(orderId uuid.UUID) ([]models.CartItem, error) {
	var items []models.CartItem
	if err := r.db.Where("order_id = ?", orderId).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
