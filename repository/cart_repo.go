package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
)

func AddCartItem(item *models.CartItem) error {
	return config.DB.Create(item).Error
}

func GetCartItems(session string) ([]models.CartItem, error) {
	var items []models.CartItem
	err := config.DB.Where("session_id = ?", session).Find(&items).Error
	return items, err
}
