package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	SessionID string `json:"session_id"`
	ProductID string `json:"product_id"`
	Quantity  string `json:"quantity"`
}
