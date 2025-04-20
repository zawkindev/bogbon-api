package models

import (
	"gorm.io/gorm"
	"time"
)

// Category: fixed taxonomy
type Category struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	Name     string    `gorm:"unique;not null"`
	Products []Product `gorm:"many2many:category_products;constraint:OnDelete:CASCADE;"`
}

// Product: can be "plant" or "service" determined by Type
type Product struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Price       int        `gorm:"not null"`
	Stock       int        `gorm:"not null"`
	Type        string     `gorm:"type:VARCHAR(20);not null"`
	Image       string     `gorm:"size:255"`
	Categories  []Category `gorm:"many2many:category_products;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Cart model: holds the cart items before checkout
type Cart struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	SessionID string     `gorm:"index;not null;unique"`
	Items     []CartItem `gorm:"foreignKey:CartID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CartItem model: now belongs to a Cart (not directly to SessionID or OrderID)
type CartItem struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CartID    uint `gorm:"index;not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null;default:1"`
	Product   Product
}

// Order model: created from a Cart
type Order struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	SessionID string `gorm:"index;not null"`
	CartID    uint   `gorm:"not null"`
	IsPaid    bool   `gorm:"default:false"`
	CreatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem model: copies data from CartItems into Order
type OrderItem struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	OrderID   uint `gorm:"index;not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	Product   Product
}
