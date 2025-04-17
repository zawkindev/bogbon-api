package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
)

func CreateOrder(o *models.Order) error {
	return config.DB.Create(o).Error
}

func GetOrderBySession(session string) (*models.Order, error) {
	var o models.Order
	err := config.DB.Preload("Items").Where("session_id = ?", session).First(&o).Error
	return &o, err
}
