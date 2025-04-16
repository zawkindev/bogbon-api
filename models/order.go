package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	SessionID string     `json:"session_id"`
	IsPaid    bool       `json:"is_paid"`
	Items     []CartItem `json:"items" gorm:"foreignKey:SessionID;references:SessionID`
}
