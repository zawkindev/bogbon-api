package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
	"errors"

	"gorm.io/gorm"
)

// EnsureCart returns the existing cart for this session or creates one.
func EnsureCart(sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := config.DB.Preload("Items").Where("session_id = ?", sessionID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart = models.Cart{SessionID: sessionID}
		if err := config.DB.Create(&cart).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &cart, nil
}

// AddCartItem adds a new item to the cart.
func AddCartItem(item *models.CartItem) error {
	return config.DB.Create(item).Error
}

// GetCart returns the cart (with items) for a session.
func GetCart(sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := config.DB.Preload("Items.Product").
		Where("session_id = ?", sessionID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// UpdateCartItem updates the quantity of a CartItem by its ID.
func UpdateCartItem(id uint, quantity int) error {
	return config.DB.Model(&models.CartItem{}).
		Where("id = ?", id).
		Update("quantity", quantity).Error
}

// DeleteCartItem removes a CartItem by its ID.
func DeleteCartItem(id uint) error {
	return config.DB.Delete(&models.CartItem{}, id).Error
}

// ClearCart deletes all items in the user's cart.
func ClearCart(sessionID string) error {
	cart, err := GetCart(sessionID)
	if err != nil {
		return err
	}
	if len(cart.Items) == 0 {
		return nil
	}
	return config.DB.
		Where("cart_id = ?", cart.ID).
		Delete(&models.CartItem{}).Error
}

