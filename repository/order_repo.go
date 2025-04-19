package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
	"errors"
)

// CreateOrderFromCart creates an Order by copying current Cart items.
// It returns the newly created Order, with its Items preloaded.
func CreateOrderFromCart(sessionID string) (*models.Order, error) {
	// 1) Load the cart and its items
	var cart models.Cart
	if err := config.DB.Preload("Items").Where("session_id = ?", sessionID).First(&cart).Error; err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 2) Create the order record
	order := models.Order{
		SessionID: sessionID,
		CartID:    cart.ID,
		IsPaid:    false,
	}
	if err := config.DB.Create(&order).Error; err != nil {
		return nil, err
	}

	// 3) Copy cart items into order_items
	for _, ci := range cart.Items {
		oi := models.OrderItem{
			OrderID:   order.ID,
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
		}
		if err := config.DB.Create(&oi).Error; err != nil {
			return nil, err
		}
	}

	// 4) Clear the cart
	if err := ClearCart(sessionID); err != nil {
		return nil, err
	}

	// 5) Reload order with its items
	if err := config.DB.Preload("Items.Product").First(&order, order.ID).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderBySession returns the most recent order for a session.
func GetOrderBySession(sessionID string) (*models.Order, error) {
	var order models.Order
	err := config.DB.Preload("Items.Product").
		Where("session_id = ?", sessionID).
		Order("created_at DESC").
		First(&order).Error
	return &order, err
}

// GetAllOrders returns every order in the system (with items).
func GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := config.DB.Preload("Items.Product").Find(&orders).Error
	return orders, err
}

// UpdateOrder marks an order as paid/unpaid.
func UpdateOrder(order *models.Order) error {
	return config.DB.Model(&models.Order{}).
		Where("id = ?", order.ID).
		Update("is_paid", order.IsPaid).
		Error
}

// DeleteOrderBySession deletes all orders for the given session.
func DeleteOrderBySession(sessionID string) error {
	return config.DB.Where("session_id = ?", sessionID).
		Delete(&models.Order{}).Error
}

