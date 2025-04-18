package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
)

// Add a new cart item
func AddCartItem(item *models.CartItem) error {
	return config.DB.Create(item).Error
}

// Get all cart items by session ID
func GetCartItems(sessionID string) ([]models.CartItem, error) {
	var items []models.CartItem
	err := config.DB.Where("session_id = ?", sessionID).Find(&items).Error
	return items, err
}

// Update cart item quantity by ID
func UpdateCartItem(id uint, quantity int) error {
	return config.DB.Model(&models.CartItem{}).Where("id = ?", id).Update("quantity", quantity).Error
}

// Delete a specific cart item by ID
func DeleteCartItem(id uint) error {
	return config.DB.Delete(&models.CartItem{}, id).Error
}

// Clear all cart items for a session
func ClearCart(sessionID string) error {
	return config.DB.Where("session_id = ?", sessionID).Delete(&models.CartItem{}).Error
}

