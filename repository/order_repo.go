package repository

import (
	"errors"

	"bogbon-api/config"
	"bogbon-api/models"
	"gorm.io/gorm"
)

// Create a new order
func CreateOrder(o *models.Order) error {
	return config.DB.Create(o).Error
}

// Get a single order by session ID
func GetOrderBySession(sessionID string) (*models.Order, error) {
	var order models.Order
	err := config.DB.Preload("Items").Where("session_id = ?", sessionID).First(&order).Error
	return &order, err
}

// Get all orders
func GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := config.DB.Preload("Items").Find(&orders).Error
	return orders, err
}

// Update an order (e.g., mark as paid)
func UpdateOrder(order *models.Order) error {
	return config.DB.Save(order).Error
}

// Delete order by session ID
func DeleteOrderBySession(sessionID string) error {
	return config.DB.Where("session_id = ?", sessionID).Delete(&models.Order{}).Error
}

// Ensure that an order exists for a session
func EnsureOrder(sessionID string) (*models.Order, error) {
	var order models.Order
	err := config.DB.Where("session_id = ?", sessionID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		order = models.Order{SessionID: sessionID, IsPaid: false}
		if err := CreateOrder(&order); err != nil {
			return nil, err
		}
	}
	return &order, err
}
